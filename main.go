package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/muci-cell/go-http/routes"
	"github.com/muci-cell/go-http/services"
)

func main() {
    r := gin.Default()
	services.MarketMap, services.Err = services.CodeIDMapEm() // 使用 var 声明的变量进行初始化
	if services.Err != nil {
		log.Fatal("获取市场代码映射失败:", services.Err) // 处理错误并终止程序
	}

    // 注册路由
    routes.SetupRoutes(r)

    // 启动服务器
    r.Run(":8123")
}
