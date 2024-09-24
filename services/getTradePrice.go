package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type KlineResponse struct {
	Rc     int      `json:"rc"`
	Rt     int      `json:"rt"`
	Svr    int      `json:"svr"`
	Lt     int      `json:"lt"`
	Full   int      `json:"full"`
	Dlmkts string   `json:"dlmkts"`
	Data   KlineData `json:"data"`
}

type KlineData struct {
	Code      string   `json:"code"`
	Market    int      `json:"market"`
	Name      string   `json:"name"`
	Decimal   int      `json:"decimal"`
	Dktotal   int      `json:"dktotal"`
	PreKPrice float64  `json:"preKPrice"`
	Klines    []string `json:"klines"`
}

func GetTradePrice(symbol, period, startDate, endDate, adjust string, timeout float64) (*[]string, error) {
	adjustDict := map[string]string{"qfq": "1", "hfq": "2", "": "0"}
	periodDict := map[string]string{"daily": "101", "weekly": "102", "monthly": "103"}
	baseURL := "http://push2his.eastmoney.com/api/qt/stock/kline/get"
	macket,err:=GetMarketID(symbol)
	if err!=nil{

		return  nil, err
	}
	// 手动构造查询参数
	query := fmt.Sprintf(
		"fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f116&ut=7eea3edcaed734bea9cbfc24409ed989&klt=%s&fqt=%s&secid=%s.%s&beg=%s&end=%s&_=%d",
		periodDict[period],
		adjustDict[adjust],
		macket,
		symbol,
		startDate,
		endDate,
		time.Now().UnixNano()/1e6,
	)

	url := fmt.Sprintf("%s?%s", baseURL, query)

	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data KlineResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data.Data.Klines, nil
}

