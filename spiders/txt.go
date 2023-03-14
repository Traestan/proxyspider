package spiders

import (
	"context"
	"regexp"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type txt struct {
	source    string //
	paths     []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func TxtpxSpider(logger *logger.Logger, gcht chan models.Proxy) Spider {
	r, _ := regexp.Compile("(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}):([\\d]{1,5})")
	svc := &txt{
		source: "txt",

		paths: []string{
			"http://static.fatezero.org/tmp/proxy.txt",
			"http://pubproxy.com/api/proxy?limit=20&format=txt&type=http",
			"http://comp0.ru/downloads/proxylist.txt",
			"http://www.proxylists.net/http_highanon.txt",
			"http://www.proxylists.net/http.txt",
			"http://ab57.ru/downloads/proxylist.txt",
			"https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list",
			"https://raw.githubusercontent.com/a2u/free-proxy-list/master/free-proxy-list.txt",
			"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list.txt",
			"https://raw.githubusercontent.com/opsxcq/proxy-list/master/list.txt",
			"https://proxyscrape.com/proxies/HTTP_Working_Proxies.txt",
			"https://proxyscrape.com/proxies/Socks4_Working_Proxies.txt",
			"https://proxyscrape.com/proxies/Socks5_Working_Proxies.txt",
			"https://proxyscrape.com/proxies/HTTP_Transparent_Proxies.txt",
			"https://proxyscrape.com/proxies/HTTP_Anonymous_Proxies.txt",
			"https://proxyscrape.com/proxies/HTTP_Elite_Proxies.txt",
			"https://proxyscrape.com/proxies/HTTP_5000ms_Timeout_Proxies.txt",
			"https://proxyscrape.com/proxies/Socks5_5000ms_Timeout_Proxies.txt",
			"https://proxyscrape.com/proxies/Socks4_5000ms_Timeout_Proxies.txt",
			"https://proxyscrape.com/proxies/HTTP_SSL_Proxies_5000ms_Timeout_Proxies.txt",
			"http://pubproxy.com/api/proxy?limit=40&format=txt&type=http",
			"https://raw.githubusercontent.com/stamparm/aux/master/fetch-some-list.txt",
			"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
			"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt",
			"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt",
		},
		ipPattern: r,
		logger:    logger,
		ch:        gcht,
	}

	return svc
}

func (srv txt) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}

	for _, v := range srv.paths {
		slug := v
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

func (srv txt) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv txt) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
