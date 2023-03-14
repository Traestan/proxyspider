package spiders

import (
	"context"
	"regexp"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type proxz struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func ProxzSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &proxz{
		source: "proxz",
		url:    "http://www.proxz.com/proxy_list_high_anonymous_0.html",

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}
func (srv proxz) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}

	slug := srv.url
	document, err := parser.GetContent(ctx, slug)
	if err != nil {
		srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
	}

	pageContent := document.Text()
	proxyPaths, _ := srv.Find(ctx, pageContent)
	for _, v := range proxyPaths {
		proxy := models.Proxy{
			IP:     v,
			Source: srv.source,
		}
		srv.ch <- proxy
	}
	proxy := models.Proxy{
		IP:     "",
		Source: "stop",
	}
	srv.ch <- proxy
	srv.logger.Info("status stop", zap.String("service", srv.source))
	return proxys, nil
}

func (srv proxz) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv proxz) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
