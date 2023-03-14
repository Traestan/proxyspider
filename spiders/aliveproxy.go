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

type alive struct {
	source    string //
	url       string
	paths     []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
	//gch       chan string
}

func AlivepxSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &alive{
		source: "alivepx",
		url:    "http://www.aliveproxy.com/%s/",
		paths: []string{
			"socks5-list",
			"high-anonymity-proxy-list", "anonymous-proxy-list",
			"fastest-proxies", "us-proxy-list", "gb-proxy-list", "fr-proxy-list",
			"de-proxy-list", "jp-proxy-list", "ca-proxy-list", "ru-proxy-list",
			"proxy-list-port-80", "proxy-list-port-81", "proxy-list-port-3128",
			"proxy-list-port-8000", "proxy-list-port-8080",
		},
		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}

func (srv alive) GetProxy(ctx context.Context) ([]models.Proxy, error) {

	proxys := []models.Proxy{}

	for _, v := range srv.paths {
		slug := fmt.Sprintf(srv.url, v)
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
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

func (srv alive) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv alive) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
