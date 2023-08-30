learn by [bilibili](“https://www.bilibili.com/video/BV1gf4y1r79E/”) 刘丹冰老师

    client
        client.go// 客户端
    service：//服务端代码
        main.go// 服务端启动入口
        server.go//服务端主要模块
        user.go// 用户相关模块

**个人改动**

- Before:服务端退出后client 不会断开链接
- After: 将client.go main方法中 client.run()方法放入goroutine中，使用chan 变量quit 进行阻塞主线程实现程序的挂起



-----

- Before: 使用`who`指令来查看当前在线人员列表
- After：跳过自己



## use

` go run service/main.go`



`go run client/client.go` 