
# clash-node-monitor（clash节点连通性监测）
## 一、功能
- 定时检测clash节点是否可用
## 二、编译
下面的命令是编译给linux amd64平台使用的，如果需要在其他平台上部署运行，请自行修改go 编译命令。当然也可以直接使用github release中我提供的对应平台的二进制文件，这样避免自行编译了。
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o clash-node-monitor .
```

## 三、部署
### 3.1 自行准备 clash.yaml 
这个文件就是正常的clash订阅配置文件。 仓库里或者release中的clash.yaml 不是真实可用的clash订阅配置文件，而是给的一个示例。

需要替换成你自己的clash订阅配置文件内容。

如果有些节点不需要检测，则在 proxies:[]，对应节点上新增属性 disabled：true。如下方中的节点配置：
```yaml
    - { name: 香港-直连节点｜不限流量, type: vmess, server: aaaa.one, port: 80, uuid: 08f33264-9bb1-4c56-9623-e6f234ec8894, alterId: 0, cipher: auto, udp: true, network: ws, disabled: true}
```

### 3.2 编写config.yaml
下面的email配置是gmail的，你可以自行替换成你的邮件厂商，不同厂商的smtp_host、smtp_port可能都不一样，可以自己搜一下或者deepseek一下
```yaml
# 邮件配置
email:
  #smtp
  smtp_host: smtp.gmail.com
  smtp_port: 587
  username: sunli***@gmail.com
  password: aaaaaaa
  from: sunli***@gmail.com
  to:
    - 1254****@qq.com
  subject: Clash节点状态告警

# Clash配置
clash:
  # Clash配置文件路径，可以是相对路径或绝对路径，就是之前准备的clash文件（非特定需要 不用调整）
  config_path: ./clash.yaml
  # 节点检测超时时间（秒）
  timeout: 5
  # 检测间隔时间（秒）
  interval: 5
```
- username 就是发送邮件的邮箱
- password 一般是该邮箱的授权码。
- from 也是发送邮件的邮箱
- to 是接受告警邮件的邮箱
- subject 是告警邮件标题

### 3.3 程序二进制文件 clash-node-monitor
这个可以自行编译，也可以直接使用 本仓库release上提供的。

### 3.3 在linux上部署服务
1）在服务器上新增文件夹 clash-node-monitor，用于存放下面提及的各文件

2）准备好clash配置文件（clash.yaml），放在目录clash-node-monitor下

3）同时准备好程序配置文件（config.yaml），放在目录clash-node-monitor下

4）将启动脚本（start.sh）,放在clash-node-monitor目录下

5）将编译后的可执行文件（clash-node-monitor）也放在此目录下

使用启动脚本启动服务，当然也可以自行使用命令来启动服务：
```bash
./start.sh
```

## 四、扩展

### 4.1 start.sh（推荐的部署方式）
启动脚本，比较方便每次启动服务。它主要做了以下几件事：
- 查找并杀掉已有的进程
- 检查程序文件是否存在
- 设置可执行权限
- 循环启动程序，如果程序崩溃则自动重启

### 4.2 deploy.sh
提供了部署脚本（deploy.sh）方便直接在开发环境部署程序到服务器上。
mac电脑
```bash
chmod +x deploy.sh
./deploy.sh
```
window电脑，需要在 Git Bash上执行
```bash
chmod +x deploy.sh
./deploy.sh
```
原因：
我们的 deploy.sh 脚本使用了以下 Unix 命令：
- ssh ：用于远程连接服务器
- scp ：用于文件传输
- chmod ：用于修改文件权限
- rm ：用于删除文件
- mkdir ：用于创建目录

这些命令在 Windows 的 CMD 或 PowerShell 中都不是原生支持的，但在 Git Bash 中都可以直接使用，因为 Git Bash 在安装时就已经包含了这些工具。

