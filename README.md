# CloudflareScanner

本项目可以测试Cloudflare节点接入速度

程序运行完毕后，结果会保存在当前目录的`result.csv`下

## 注意事项
#### 协程数请不要调过1000，否则容易出现较大误差
#### linux编译选项为 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build`
#### windows编译命令 `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build`
#### 您可在release界面下载或编译运行
