package marketstack

import (
	"net/http"
)

type StockData struct {
	Open   float32 `json:"open"`
	High   float32 `json:"high"`
	Low    float32 `json:"low"`
	Close  float32 `json:"close"`
	Volume float32 `json:"volume"`
	Date   string  `json:"date"`
	Symbol string  `json:"symbol"`
}

type Response struct {
	Data []StockData `json:"data"`
}

type Client struct {
	http.Client
}

type T struct {
	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Count  int `json:"count"`
		Total  int `json:"total"`
	} `json:"pagination"`
	Data []struct {
		Open        float64  `json:"open"`
		High        float64  `json:"high"`
		Low         float64  `json:"low"`
		Close       float64  `json:"close"`
		Volume      float64  `json:"volume"`
		AdjHigh     *float64 `json:"adj_high"`
		AdjLow      *float64 `json:"adj_low"`
		AdjClose    float64  `json:"adj_close"`
		AdjOpen     *float64 `json:"adj_open"`
		AdjVolume   *float64 `json:"adj_volume"`
		SplitFactor float64  `json:"split_factor"`
		Dividend    float64  `json:"dividend"`
		Symbol      string   `json:"symbol"`
		Exchange    string   `json:"exchange"`
		Date        string   `json:"date"`
	} `json:"data"`
}

func GetClient() *Client {
	httpClient := &Client{}
	return httpClient
}

//func test(){
//	req, err := http.NewRequest("GET", "https://api.marketstack.com/v1/eod?access_key=bc51d42081fb4b803ff35431e018f3f9", nil)
//	if err != nil {
//		panic(err)
//	}
//
//	q := req.URL.Query()
//	q.Add("access_key", "YOUR_ACCESS_KEY")
//	req.URL.RawQuery = q.Encode()
//
//	res, err := httpClient.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer func(Body io.ReadCloser) {
//		_ = Body.Close()
//	}(res.Body)
//
//	var apiResponse Response
//	_ = json.NewDecoder(res.Body).Decode(&apiResponse)
//
//	for _, stockData := range apiResponse.Data {
//		fmt.Println(fmt.Sprintf("Ticker %s has a day high of %v on %s",
//			stockData.Symbol,
//			stockData.High,
//			stockData.Date))
//	}
//}
