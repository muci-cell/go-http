package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muci-cell/go-http/routes"
)

func main() {
    r := gin.Default()

    // 注册路由
    routes.SetupRoutes(r)

    // 启动服务器
    r.Run(":8123")
}
