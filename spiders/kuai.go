package spiders

import (
	"context"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type kuaidaili struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func KuaidailiSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
	svc := &kuaidaili{
		source: "kuaidaili",
		url:    "https://www.kuaidaili.com/free/inha/%d/",

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}

func (srv kuaidaili) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	//proxyPaths := make(map[string]string)
	for i := 1; i < 10; i++ {
		slug := fmt.Sprintf(srv.url, i)
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
		}
		if document == nil { // timeout error next slug
			continue
		}
		document.Find("tr").Each(func(index int, element *goquery.Selection) {
			if element.Find("td").Length() > 0 {
				imgSrc := element.Find("td")
				ipLocal := make([]string, 0)
				portLocal := make([]string, 0)
				rPort, _ := regexp.Compile(`([0-9]+){4}$`)
				imgSrc.Each(func(index int, element *goquery.Selection) {
					ip := srv.ipPattern.FindAllString(element.Text(), -1)
					if len(ip) > 0 {
						ipLocal = append(ipLocal, ip[0])
					}
					port := rPort.FindAllString(element.Text(), -1)
					if len(port) > 0 {
						portLocal = append(portLocal, port[0])
					}
				})

				if len(ipLocal) > 0 && len(portLocal) > 0 {
					for k, v := range ipLocal {
						proxyIpPort := v + ":" + portLocal[k]
						proxy := models.Proxy{
							IP:     proxyIpPort,
							Source: srv.source,
						}
						srv.ch <- proxy
					}
				}
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

func (srv kuaidaili) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv kuaidaili) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
