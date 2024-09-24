package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)
var MarketMap map[string]string
var Err error

type DiffData struct {
	F12 string `json:"f12"`
}

type ResponseData struct {
	Data struct {
		Diff []DiffData `json:"diff"`
	} `json:"data"`
}

func CodeIDMapEm() (map[string]string, error) {
	url := "http://80.push2.eastmoney.com/api/qt/clist/get"
	codeIDDict := make(map[string]string)

	// 1. 获取上证市场
	if err := fetchData(url, map[string]string{"fs": "m:1 t:2,m:1 t:23"}, codeIDDict, "1"); err != nil {
		return nil, err
	}

	// 2. 获取深圳市场
	if err := fetchData(url, map[string]string{"fs": "m:0 t:6,m:0 t:80"}, codeIDDict, "0"); err != nil {
		return nil, err
	}

	// 3. 获取北京市场
	if err := fetchData(url, map[string]string{"fs": "m:0 t:81 s:2048"}, codeIDDict,"0"); err != nil {
		return nil, err
	}

	return codeIDDict, nil
}

func fetchData(baseURL string, params map[string]string, codeIDDict map[string]string, idKey string) error {
	// 创建 URL 对象
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	// 构建查询参数
	q := u.Query()
	q.Set("pn", "1")
	q.Set("pz", "50000")
	q.Set("po", "1")
	q.Set("np", "1")
	q.Set("ut", "bd1d9ddb04089700cf9c27f6f7426281")
	q.Set("fltt", "2")
	q.Set("invt", "2")
	q.Set("fid", "f3")
	q.Set("fields", "f12")
	q.Set("_", "1623833739532")

	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}


	defer resp.Body.Close()
	var responseData ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return err
	}

	if len(responseData.Data.Diff) == 0 {
		return nil
	}

	for _, item := range responseData.Data.Diff {
		codeIDDict[item.F12] = idKey
	}

	return nil
}

// GetMarketID 根据证券代码返回市场编号
func GetMarketID(code string) (string, error) {

	if marketID, exists := MarketMap[code]; exists {
		return marketID, nil
	}
	return "", errors.New("证券代码未找到")
}