package spiders

import (
	"context"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type proxbug struct {
	source    string //
	paths     []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func ProxbugSpider(logger *logger.Logger, gcht chan models.Proxy) Spider {
	r, _ := regexp.Compile("(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}):([\\d]{1,5})")
	svc := &proxbug{
		source: "proxbug",
		paths: []string{
			"http://www.freshnewproxies24.top/",
			"http://www.live-socks.net/",
			"http://www.sslproxies24.top/",
			"http://www.proxyserverlist24.top/",
		},
		ipPattern: r,
		logger:    logger,
		ch:        gcht,
	}

	return svc
}

func (srv proxbug) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	for _, v := range srv.paths {
		slug := v
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
		}
		document.Find(".jump-link").Each(func(index int, element *goquery.Selection) {
			imgSrc, _ := element.Find("a").Attr("href")

			documentInner, err := parser.GetContent(ctx, imgSrc)
			if err != nil {
				srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
			}
			proxyPaths, _ := srv.Find(ctx, documentInner.Text())
			for _, v := range proxyPaths {
				proxy := models.Proxy{
					IP:     v,
					Source: srv.source,
				}
				srv.ch <- proxy
			}
		})

	}
	proxy := models.Proxy{
		IP:     "",
		Source: "stop",
	}
	srv.ch <- proxy
	srv.logger.Info("status stop", zap.String("service", srv.source))
	return proxys, nil
}

func (srv proxbug) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv proxbug) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
