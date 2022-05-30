# NGB-Notification

NGB的消息通知系统，与后端主进程之间通过RabbitMQ连接。

2种通知方式：

- websocket通知：

  `ws://localhost:8081/notification`

  - 如果用户在线，系统直接推送消息给用户

  - 否则将消息存储到Redis，等用户上线后主动拉取未读消息

- 邮件通知：发送一封包含关键信息的邮件到用户邮箱中