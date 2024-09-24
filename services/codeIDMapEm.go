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
	//         // fs: "m:113 t:15",       // 证券过滤器
        // fs: "m:1 t:2",       // 证券过滤器
	// 4.测试市场1 沪港通etf
	if err := fetchData(url, map[string]string{"fs": "b:MK0839"}, codeIDDict,"1"); err != nil {
		return nil, err
	}
	// 测试市场2 深港通etf
	if err := fetchData(url, map[string]string{"fs": "b:MK0840"}, codeIDDict,"0"); err != nil {
		return nil, err
	}

	// 测试市场3 沪港通reits
	if err := fetchData(url, map[string]string{"fs": "m:1 t:9 e:97"}, codeIDDict,"1"); err != nil {
		return nil, err
	}
	// 测试市场4 深港通reits
	if err := fetchData(url, map[string]string{"fs": "m:0 t:10 e:97"}, codeIDDict,"0"); err != nil {
		return nil, err
	}
	// 测试市场5 lof 深圳
	if err := fetchData(url, map[string]string{"fs": "b:MK0404,b:MK0405,b:MK0406,b:MK0407"}, codeIDDict,"0"); err != nil {
		return nil, err
	}
	// 测试市场6 基金 5开头是1上海    1开头是0深圳
	if err := fetchData(url, map[string]string{"fs": "b:MK0021,b:MK0022,b:MK0023,b:MK0024,b:MK0827"}, codeIDDict,"1/0"); err != nil {
		return nil, err
	}
	return codeIDDict, nil
}
func checkCodeID(idKey string, codeIDDict string) string {
	if idKey == "1/0" {
		if len(codeIDDict) > 0 {
			switch codeIDDict[0] {
			case '5':
				return "1"
			case '1':
				return "0"
			}
		}
	}
	return idKey // 如果条件不满足，直接返回市场值
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
		
		codeIDDict[item.F12]=checkCodeID(idKey,item.F12)
	}

	return nil
}

// GetMarketID 根据证券代码返回市场编号
func GetMarketID(code string) (string, error) {
	if marketID, exists := MarketMap[code]; exists {
		return marketID, nil
	}
	return "",errors.New("证券所属市场获取错误")
}