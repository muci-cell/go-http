package services

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// ETFNetAssetValue 表示历史净值的结构体
type ETFNetAssetValue struct {
    FSRQ            string  `json:"FSRQ"`      // 净值日期
    DWJZ            string  `json:"DWJZ"`      // 单位净值
    LJJZ            string  `json:"LJJZ"`      // 累计净值
    JZZZL           string  `json:"JZZZL"`     // 日增长率
    SGZT            string  `json:"SGZT"`      // 申购状态
    SHZT            string  `json:"SHZT"`      // 赎回状态
}

// ETFResponse 包含整个API返回的响应结构
type ETFResponse struct {
    Data struct {
        LSJZList []ETFNetAssetValue `json:"LSJZList"`
        FundType string             `json:"FundType"`
    } `json:"Data"`
    ErrCode    int    `json:"ErrCode"`
    ErrMsg     string `json:"ErrMsg"`
    TotalCount int    `json:"TotalCount"`
    PageSize   int    `json:"PageSize"`
    PageIndex  int    `json:"PageIndex"`
}

// GetHistoricalValue 获取历史净值
func GetHistoricalValue(fund string, startDate string, endDate string) (*[]ETFNetAssetValue, error) {
    url := "http://api.fund.eastmoney.com/f10/lsjz"
    headers := map[string]string{
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
        "Referer": fmt.Sprintf("http://fundf10.eastmoney.com/jjjz_%s.html", fund),
    }

    params := fmt.Sprintf("fundCode=%s&pageIndex=1&pageSize=10000&startDate=%s&endDate=%s&_=%.0f",
        fund, formatDate(startDate), formatDate(endDate), float64(time.Now().UnixNano())/1e6)

    // 创建请求
    req, err := http.NewRequest("GET", url+"?"+params, nil)
    if err != nil {
        return nil, err
    }

    // 设置请求头
    for key, value := range headers {
        req.Header.Set(key, value)
    }

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // 解析响应
    var result ETFResponse

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    // 返回净值列表
    return &result.Data.LSJZList, nil
}

// formatDate 格式化日期为YYYY-MM-DD
func formatDate(date string) string {
    return fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:])
}
