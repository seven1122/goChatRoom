##并发聊天室使用说明
* 启动服务端  
启动服务，监听连接，将连接用户保存，根据客户端发送的包含接受者ip的聊天信息处理，将聊天信息转发给指定ip。
* 启动客户端  
可以启动多个客户端，模拟不同的聊天用户  
* 如何聊天？  
ip#聊天内容：ip指定聊天对象，如127.0.0.1:55284#你好 在不在！
* 如何查看在线列表？  
ip#list: ip一般是自己的ip(否则在线列表发送到指定的其他用户窗口了！)
* 如何退出聊天？  
输入exit命令。