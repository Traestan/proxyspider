package spiders

import (
	"context"
	"regexp"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type stdian struct {
	source    string //
	paths     []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func StdianSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &stdian{
		source: "stdian",
		paths: []string{
			"https://free-proxy-list.net/uk-proxy.html",
			"https://free-proxy-list.net/",
			"https://free-proxy-list.net/anonymous-proxy.html",
			"https://www.sslproxies.org/",
		},
		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}
	return svc
}

func (srv stdian) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	for _, v := range srv.paths {
		slug := v

		//go func() {
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
			proxys = append(proxys, proxy)
		}
	}
	proxy := models.Proxy{
		IP:     "",
		Source: "stop",
	}
	srv.ch <- proxy
	return proxys, nil
}

func (srv stdian) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv stdian) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
