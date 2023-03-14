package parser

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
)

// GetContent получаем контент со страницы
func GetContent(ctx context.Context, slug string) (*goquery.Document, error) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	// Create HTTP client with timeout
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", slug, nil)
	if err != nil {
		logger.Log("err", "request", "slug", slug, err)
	}
	request.Header.Set("User-Agent", "Not Firefox")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		logger.Log("err", "request", "slug", slug, err)
		//continue
		return nil, nil
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	// document.Find("table").Each(func(index int, element *goquery.Selection) {
	// 	imgSrc := element.Find("td").Text()
	// 	fmt.Println(imgSrc)
	// })

	return document, nil
}
func GetContentWithHeaders(ctx context.Context, slug string, headers map[string]string) (*goquery.Document, error) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	// Create HTTP client with timeout
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	request, err := http.NewRequest("GET", slug, nil)
	if err != nil {
		logger.Log("err", "request", "slug", slug, err)
	}
	//request.Header.Set("User-Agent", "Not Firefox")
	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	// Make request
	response, err := client.Do(request)
	if err != nil {
		logger.Log("err", "request", "slug", slug, err)
		//continue
		return nil, nil
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func GetContentJson(ctx context.Context, slug string) (map[string]interface{}, error) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	// Create HTTP client with timeout
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", slug, nil)
	if err != nil {
		logger.Log("err", "request", "slug", slug, err)
	}
	request.Header.Set("User-Agent", "Not Firefox")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		logger.Log("err", "request", "slug", slug, err)
		//continue
		return nil, nil
	}

	defer response.Body.Close()

	resp := make(map[string]interface{})
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
