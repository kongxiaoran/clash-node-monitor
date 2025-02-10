package checker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/Dreamacro/clash/adapter/outbound"
	"github.com/Dreamacro/clash/constant"
	"gopkg.in/yaml.v3"
)

type ClashConfig struct {
	Proxies []Proxy `yaml:"proxies"`
}

type Proxy struct {
	Name          string    `yaml:"name"`
	Server        string    `yaml:"server"`
	Port          int       `yaml:"port"`
	Type          string    `yaml:"type"`
	Disabled      bool      `yaml:"disabled"`
	Cipher        string    `yaml:"cipher"`   // SS 加密方式
	Password      string    `yaml:"password"` // SS 密码
	FailureCount  int       // 连续失败次数
	LastAlertTime time.Time // 上次发送告警的时间
}

type CheckResult struct {
	Name        string
	Latency     time.Duration
	Error       error
	ShouldAlert bool // 是否需要发送告警
}

func LoadClashConfig(path string) (*ClashConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ClashConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func CheckProxy(proxy *Proxy, timeout int) CheckResult {
	result := CheckResult{Name: proxy.Name}
	start := time.Now()

	// 设置本地 DNS 解析器
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(timeout) * time.Second,
			}
			// 轮询使用多个 DNS 服务器
			dnsServers := []string{
				"223.5.5.5:53",    // 阿里 DNS
				"119.29.29.29:53", // 腾讯 DNS
				"180.76.76.76:53", // 百度 DNS
			}

			// 尝试连接每个 DNS 服务器
			var lastErr error
			for _, server := range dnsServers {
				conn, err := d.DialContext(ctx, "udp", server)
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}
			return nil, fmt.Errorf("所有 DNS 服务器均连接失败，最后错误: %v", lastErr)
		},
	}

	// 使用自定义解析器解析域名
	ctx := context.Background()
	ips, err := resolver.LookupHost(ctx, proxy.Server)
	if err != nil {
		log.Printf("DNS解析失败 %s: %v, 继续使用原始域名", proxy.Server, err)
	} else if len(ips) > 0 {
		//log.Printf("域名 %s 解析到 IP: %v", proxy.Server, ips[0])
		proxy.Server = ips[0]
	}

	// 创建代理适配器
	var proxyAdapter constant.ProxyAdapter

	switch proxy.Type {
	case "ss":
		proxyAdapter, err = outbound.NewShadowSocks(outbound.ShadowSocksOption{
			Cipher:     proxy.Cipher,
			Name:       proxy.Name,
			Password:   proxy.Password,
			Server:     proxy.Server,
			Port:       proxy.Port,
			UDP:        true,
			Plugin:     "",
			PluginOpts: nil,
		})
	case "vmess":
		result.Error = fmt.Errorf("暂不支持 VMess 类型代理")
		log.Printf("节点 %s 错误: %v", proxy.Name, result.Error)
		proxy.FailureCount++
		return result
	default:
		result.Error = fmt.Errorf("不支持的代理类型: %s", proxy.Type)
		log.Printf("节点 %s 错误: %v", proxy.Name, result.Error)
		proxy.FailureCount++
		return result
	}

	if err != nil {
		result.Error = fmt.Errorf("创建代理失败: %v", err)
		log.Printf("节点 %s 错误: %v", proxy.Name, result.Error)
		proxy.FailureCount++
		return result
	}

	// 测试 TCP 连接
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 使用 8.8.8.8:53 作为测试目标（Google DNS）
	conn, err := proxyAdapter.DialContext(ctx, &constant.Metadata{
		Type:    constant.SOCKS5,
		NetWork: constant.TCP,
		Host:    "8.8.8.8",
		DstPort: 53,
	})

	if err != nil {
		result.Error = fmt.Errorf("TCP 连接失败: %v", err)
		log.Printf("节点 %s TCP 测试失败: %v", proxy.Name, err)
		proxy.FailureCount++
		return result
	}
	conn.Close()

	// 如果代理支持 UDP，进行 UDP 测试
	if proxyAdapter.SupportUDP() {
		pc, err := proxyAdapter.ListenPacketContext(ctx, &constant.Metadata{
			NetWork: constant.UDP,
			Host:    "8.8.8.8",
			DstPort: 53,
		})

		if err != nil {
			result.Error = fmt.Errorf("UDP 连接失败: %v", err)
			log.Printf("节点 %s UDP 测试失败: %v", proxy.Name, err)
			proxy.FailureCount++
			return result
		}
		pc.Close()
	}

	// 测试成功，重置失败计数
	proxy.FailureCount = 0
	// 计算延迟并返回结果
	result.Latency = time.Since(start)
	log.Printf("节点 %s 测试成功，延迟: %s", proxy.Name, result.Latency)
	return result
}

func CheckAllProxies(config *ClashConfig, timeout int) []CheckResult {
	var results []CheckResult

	for i := range config.Proxies {
		if !config.Proxies[i].Disabled {
			result := CheckProxy(&config.Proxies[i], timeout)

			// 检查是否连续失败三次且距离上次告警超过1小时
			if config.Proxies[i].FailureCount >= 3 {
				now := time.Now()
				if now.Sub(config.Proxies[i].LastAlertTime) > time.Hour {
					log.Printf("警告：节点 %s 已连续失败 %d 次", config.Proxies[i].Name, config.Proxies[i].FailureCount)
					config.Proxies[i].LastAlertTime = now
					result.ShouldAlert = true
				}
			}

			results = append(results, result)
		}
	}

	return results
}
