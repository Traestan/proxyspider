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

type opnproxy struct {
	source    string //
	paths     []string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func OpnproxySpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("(?:((?:\\d|[1-9]\\d|1\\d{2}|2[0-5][0-5])\\.(?:\\d|[1-9]\\d|1\\d{2}|2[0-5][0-5])\\.(?:\\d|[1-9]\\d|1\\d{2}|2[0-5][0-5])\\.(?:\\d|[1-9]\\d|1\\d{2}|2[0-5][0-5]))\\D+?(6[0-5]{2}[0-3][0-5]|[1-5]\\d{4}|[1-9]\\d{1,3}|[0-9]))")
	svc := &opnproxy{
		source: "opnproxy",
		paths: []string{
			"https://www.kuaidaili.com/free", "https://31f.cn/region/北京/", "https://31f.cn/region/广东/#", "https://www.kuaidaili.com/free/inha/2/",
			"https://31f.cn/region/安徽/", "http://www.31f.cn", "https://free-proxy-list.net/anonymous-proxy.html",
			"https://www.kuaidaili.com/free/inha/3/",
			"https://www.us-proxy.org/", "https://www.kuaidaili.com/free/inha/4/",
			"http://www.ip181.com/", "https://www.free-proxy-list.net/", "https://free-proxy-list.net/anonymous-proxy.html",
			"https://www.proxynova.com/proxy-server-list/country-us/", "https://www.ip-adress.com/proxy-list",
			"https://www.proxynova.com/proxy-server-list/",
			"http://www.proxy-daily.com/", "https://www.kuaidaili.com/free/inha/5/", "http://202.112.51.31:5010/get_all/",
			"http://www.data5u.com/", "http://www.goubanjia.com/",
		},
		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}
func (srv opnproxy) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	for _, v := range srv.paths {
		slug := v
		//srv.logger.Log("msg", srv.source, "request url", slug)

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
			strIP := strings.ReplaceAll(v, "\n", "")
			strIP = strings.Replace(strIP, " ", ":", 1)
			strIP = strings.ReplaceAll(strIP, " ", "")
			strIP = strings.ReplaceAll(strIP, "');", "")
			proxy := models.Proxy{
				IP:     strIP,
				Source: srv.source,
			}
			srv.ch <- proxy
			//proxys = append(proxys, proxy)
		}
	}
	proxy := models.Proxy{
		IP:     "",
		Source: "stop",
	}
	srv.ch <- proxy
	return proxys, nil
}

func (srv opnproxy) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv opnproxy) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
