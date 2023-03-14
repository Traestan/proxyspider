package spiders

import (
	"context"
	"regexp"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type freepxlist struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func FreepxlistSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &freepxlist{
		source: "freepxlist",
		url:    "http://www.proxz.com/proxy_list_high_anonymous_0.html",

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}
func (srv freepxlist) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}

	slug := srv.url
	//srv.logger.Log("msg", "Alivepx", "request url", slug)
	// go func(slug string, ch chan models.Proxy) {
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
	// srv.gch <- srv.source + "-stop"
	srv.logger.Info("status stop", zap.String("service", srv.source))
	return proxys, nil
}

func (srv freepxlist) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv freepxlist) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
