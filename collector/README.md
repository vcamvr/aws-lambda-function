event_demo.json 描述了原始的 Event 消息体，在 EventBridge 中基于该消息体
提取了部分字段，简化了调用函数的请求参数。

### Building
```txt
# Remember to build your handler executable for Linux!
GOOS=linux GOARCH=amd64 go build -o main update_collector.go
zip main.zip main
```
