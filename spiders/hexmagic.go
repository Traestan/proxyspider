package spiders

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type hexmagic struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
	paths     []string
}

func HexmagicSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
	svc := &hexmagic{
		source: "hexmagic",
		url:    "http://www.ip3366.net/?stype=1&page=%d",
		paths: []string{
			"http://www.ip3366.net/?stype=1&page=%d",
			"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-%d",
		},
		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}

func (srv hexmagic) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	for _, v := range srv.paths {
		//go func(parentSlug string) {
		for i := 1; i < 5; i++ {
			slug := fmt.Sprintf(v, i)
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
							proxyIPPort := v + ":" + portLocal[k]
							proxy := models.Proxy{
								IP:     proxyIPPort,
								Source: srv.source,
							}
							srv.ch <- proxy
						}
					}
				}
			})
		}
		//}(v)
	}
	proxy := models.Proxy{
		IP:     "",
		Source: "stop",
	}
	srv.ch <- proxy

	srv.logger.Info("status stop", zap.String("service", srv.source))
	return proxys, nil
}

//
func (srv hexmagic) JsonPx(ctx context.Context) error {
	slug := "http://ip.jiangxianli.com/api/proxy_ips"
	document, err := parser.GetContentJson(ctx, slug)
	if err != nil {
		srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
	}
	if arr, ok := document["data"].(map[string]interface{}); ok {
		for _, v := range arr {
			if arrLoc, ok := v.(map[string]interface{}); ok {
				fmt.Println(arrLoc)
			}
		}
	}
	return nil
}
func (srv hexmagic) Iphai(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	parth := []string{
		"http://www.iphai.com/free/ng",
		"http://www.iphai.com/free/np",
		"http://www.iphai.com/free/wg",
		"http://www.iphai.com/free/wp",
	}

	for _, v := range parth {
		slug := v
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
		}
		if document == nil { // timeout error next slug
			continue
		}

		document.Find("tr").Each(func(index int, element *goquery.Selection) {
			if element.Find("td").Length() > 0 {
				ipLocal := make([]string, 0)
				portLocal := make([]string, 0)

				imgSrc := element.Find("td")
				// rPort, _ := regexp.Compile(`([0-9]+.*){2,4}$`)
				re := regexp.MustCompile(`\n`)
				imgSrc.Each(func(index int, element *goquery.Selection) {
					ip := srv.ipPattern.FindAllString(element.Text(), -1)
					if len(ip) > 0 {
						ipLocal = append(ipLocal, ip[0])
					}
					if index == 1 {
						t := re.ReplaceAllString(element.Text(), "")
						portLocal = append(portLocal, strings.ReplaceAll(t, " ", ""))
					}
				})

				if len(ipLocal) > 0 && len(portLocal) > 0 {
					for k, v := range ipLocal {
						proxyIPPort := v + ":" + strings.ReplaceAll(portLocal[k], `\n`, "")
						proxy := models.Proxy{
							IP:     proxyIPPort,
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

func (srv hexmagic) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv hexmagic) Run(ctx context.Context) {
	srv.GetProxy(ctx)
	//srv.JsonPx(ctx)-доработка
	srv.Iphai(ctx)
}
