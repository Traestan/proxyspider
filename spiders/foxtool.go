package spiders

import (
	"context"
	"fmt"
	"regexp"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type foxtool struct {
	source    string //
	url       string
	paths     []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
	//gch       chan string
}

func FoxtoolpxSpider(logger *logger.Logger, gcht chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &foxtool{
		source:    "foxtool",
		url:       "http://api.foxtools.ru/v2/Proxy.txt?page=%d",
		ipPattern: r,
		logger:    logger,
		ch:        gcht,
	}
	return svc
}

func (srv foxtool) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	for i := 1; i < 5; i++ {

		slug := fmt.Sprintf(srv.url, i)
		// srv.logger.Log("msg", srv.source, "request url", slug)

		//go func() {
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("request url", zap.String("slug", slug), zap.Error(err))
		}
		if document == nil { // timeout error next slug
			continue
		}
		pageContent := document.Text()
		proxyPaths, _ := srv.Find(ctx, pageContent)
		for _, v := range proxyPaths {
			proxy := models.Proxy{
				IP:     v,
				Source: srv.source,
			}
			srv.ch <- proxy
			proxys = append(proxys, proxy)
		}

	}
	proxy := models.Proxy{
		IP:     "",
		Source: "stop",
	}
	srv.ch <- proxy
	srv.logger.Info("status stop", zap.String("service", srv.source))
	return proxys, nil
}

func (srv foxtool) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv foxtool) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
