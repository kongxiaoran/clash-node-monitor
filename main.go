package main

import (
	"flag"
	"log"
	"time"

	"clash-node-monitor/checker"
	"clash-node-monitor/config"
	"clash-node-monitor/mailer"
)

func main() {
	// 初始化日志
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 读取配置文件
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化邮件客户端
	mailClient := mailer.NewMailer(cfg.Email)

	// 启动定时任务
	ticker := time.NewTicker(time.Duration(cfg.Clash.Interval) * time.Second)
	defer ticker.Stop()

	log.Printf("节点监控服务已启动，检测间隔: %d秒\n", cfg.Clash.Interval)

	// 加载Clash配置
	clashCfg, err := checker.LoadClashConfig(cfg.Clash.ConfigPath)
	if err != nil {
		log.Printf("加载Clash配置失败: %v\n", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			log.Println("开始检测节点状态...")

			// 执行节点检测
			results := checker.CheckAllProxies(clashCfg, cfg.Clash.Timeout)

			// 发送告警邮件
			if err := mailClient.SendAlertEmail(results); err != nil {
				log.Printf("发送告警邮件失败: %v\n", err)
			}
		}
	}
}
