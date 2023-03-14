package spiders

import (
	"context"
	"regexp"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type miniprox struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func MiniproxSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &miniprox{
		source: "miniprox",
		url:    "http://proxy-ip-list.com/fresh-proxy-list.html",

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}
func (srv miniprox) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	//proxyPaths := make(map[string]string)
	matchesRow := make([]string, 1)

	slug := srv.url
	document, err := parser.GetContent(ctx, slug)
	if err != nil {
		srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
	}

	pageContent, _ := document.Html()

	patternTable, _ := regexp.Compile(`<td>(.*?)</td>`)
	matchesRow = patternTable.FindAllString(pageContent, -1)

	for _, row := range matchesRow {
		proxys, err := srv.Find(ctx, row)
		if err != nil {
			return nil, err
		}
		if len(proxys) != 0 {
			proxy := models.Proxy{
				IP:     proxys[0],
				Source: srv.source,
			}
			srv.ch <- proxy
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
func (srv miniprox) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv miniprox) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
