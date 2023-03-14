package spiders

import (
	"context"
	"regexp"
	"strings"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type scylla struct {
	source    string //
	urls      []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func ScyllaSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile(`<td>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}</td>\n<td>\d{2,5}</td>`)
	svc := &scylla{
		source: "scylla",
		urls: []string{
			"http://31f.cn/socks-proxy/",
			"http://31f.cn/http-proxy/",
			"http://31f.cn/https-proxy/",
		},

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}
func (srv scylla) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}

	re := regexp.MustCompile(`\n`)
	for _, v := range srv.urls {
		slug := v
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
		}
		if document == nil { // timeout error next slug
			continue
		}
		pageContent, _ := document.Html()
		proxyPaths, _ := srv.Find(ctx, pageContent)

		for _, v := range proxyPaths {
			t := strings.ReplaceAll(re.ReplaceAllString(v, ""), "</td><td>", ":")
			t = strings.ReplaceAll(t, "<td>", "")
			t = strings.ReplaceAll(t, "</td>", "")
			proxy := models.Proxy{
				IP:     t,
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

func (srv scylla) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv scylla) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
