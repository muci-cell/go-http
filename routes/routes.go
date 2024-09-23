package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muci-cell/go-http/services"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
    // 读取交易价格的API
    r.GET("/api/trade-price", func(c *gin.Context) {
		symbol := c.Query("symbol")
		period := c.Query("period")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		adjust := c.Query("adjust")
		timeout := c.Query("timeout")

		timeoutFloat, _ := strconv.ParseFloat(timeout, 64)

		price, err := services.GetTradePrice(symbol, period, startDate, endDate, adjust, timeoutFloat)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"trade_price": price})
    })

    // 读取历史净值的API
    r.GET("/api/historical-value", func(c *gin.Context) {
		fundCode := c.Query("fundCode")
		startDate := c.Query("startDate")
		endDate  := c.Query("endDate")
        netValue, err := services.GetHistoricalValue(fundCode,startDate,endDate)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"historical_net_value": netValue})
    })
}