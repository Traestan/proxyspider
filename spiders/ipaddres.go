package spiders

import (
	"context"
	"regexp"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type ipaddres struct {
	source    string //
	url       string
	paths     []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
	//gch       chan string
}

func IpaddrespxSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &ipaddres{
		source:    "ipaddres",
		url:       "https://www.ipaddress.com/proxy-list/",
		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}

func (srv ipaddres) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	slug := srv.url

	document, err := parser.GetContent(ctx, slug)
	if err != nil {
		srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
	}
	if document == nil { // timeout error next slug
		return proxys, nil
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

	proxy := models.Proxy{
		IP:     "",
		Source: "stop",
	}
	srv.ch <- proxy

	srv.logger.Info("status stop", zap.String("service", srv.source))
	return proxys, nil
}

func (srv ipaddres) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv ipaddres) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
