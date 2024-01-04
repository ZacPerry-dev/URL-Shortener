package server

type baseUrlInfo struct {
	LongUrl string `json:"LongUrl"`
}

type newUrlInfo struct {
	LongUrl  string `json:"LongUrl"`
	ShortUrl string `json:"ShortUrl"`
	Key      string `json:"Key"`
}
