package utils

import (
	"strings" 
	"time"
	"log"
)

const (
	layout = "2006-01-02"
)

func GetRedirectUrl(referer string) string {
	var redirectUrl string
	url := strings.Split(referer, "/")

	if len(url) > 4 {
		redirectUrl = "/" + strings.Join(url[3:], "/")
	} else {
		redirectUrl = "/"
	}
	return redirectUrl
}

func DateToUnix(date string) int64 {
	t, err := time.Parse(layout, date)

	if err != nil  {
		log.Println(err)
	} else {
		return t.Unix()
	}
	return 0
}

func UnixToDate(unix int64) string {
	return "sas"
} 