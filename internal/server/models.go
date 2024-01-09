package server

type baseUrlInfo struct {
	Url string `json:"Url"`
}

type newUrlInfo struct {
	LongUrl  string `json:"LongUrl"`
	ShortUrl string `json:"ShortUrl"`
	Key      string `json:"Key"`
}
