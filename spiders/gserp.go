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

type gserp struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func GserppxSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("/url\\?q=((?!.*webcache)(?:(?:https?|ftp)://|www\\.|ftp\\.)[^'\r\n]+\\.txt)")
	svc := &gserp{
		source: "gserp",
		url:    "https://www.google.com/search?q=+\":8080\" +\":3128\" +\":80\" filetype:txt&start=%s0",

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}
func (srv gserp) GetProxy(ctx context.Context) ([]models.Proxy, error) {

	proxys := []models.Proxy{}

	for i := 1; i <= 7; i++ {
		slug := fmt.Sprintf(srv.url, i)
		// go func(slug string, ch chan models.Proxy) {

		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
		}
		if document == nil { // timeout error next slug
			continue
		}
		pageContent := document.Text()
		fmt.Println(pageContent)
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
		IP:     "none",
		Source: "stop",
	}
	srv.ch <- proxy
	// srv.gch <- srv.source + "-stop"
	srv.logger.Info("status stop", zap.String("service", srv.source))
	return proxys, nil
}

func (srv gserp) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv gserp) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
