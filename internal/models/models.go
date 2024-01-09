package models

type BaseUrlInfo struct {
	Url string `json:"Url"`
}

type NewUrlInfo struct {
	LongUrl  string `json:"LongUrl"`
	ShortUrl string `json:"ShortUrl"`
	Key      string `json:"Key"`
}
