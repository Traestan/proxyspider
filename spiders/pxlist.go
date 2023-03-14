package spiders

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type pxlist struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func PxlistpxSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("Proxy\\('([\\w=]+)'\\)")
	svc := &pxlist{
		source: "pxlistpx",
		url:    "http://proxy-list.org/english/%s",

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}

func (srv pxlist) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	//r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	proxys := []models.Proxy{}
	url := "http://proxy-list.org/english/index.php?p=1"
	document, err := parser.GetContent(ctx, url)
	if err != nil {
		srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", url), zap.Error(err))
	}
	pageContent, err := document.Html()
	if err != nil {
		return nil, err
	}
	pages, _ := srv.FindPage(ctx, pageContent)
	for _, v := range pages {
		slug := fmt.Sprintf(srv.url, v)

		//go func(slug string, ch chan models.Proxy) {
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
		}
		pageContent := document.Text()
		proxyPaths, _ := srv.Find(ctx, pageContent)
		for _, v := range proxyPaths {
			sDec, _ := base64.URLEncoding.DecodeString(v)
			proxy := models.Proxy{
				IP:     string(sDec),
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

func (srv pxlist) FindPage(ctx context.Context, pageText string) ([]string, error) {
	pages := []string{}
	r, _ := regexp.Compile(`href\s*=\s*['"]\./([^'"]?index\.php\?p=\d+[^'"]*)['"]`)
	for _, v := range r.FindAllString(pageText, -1) {
		page := strings.ReplaceAll(strings.ReplaceAll(v, "href=\"./", ""), "\"", "") //re.ReplaceAllString(v, "")
		pages = append(pages, page)
	}
	return pages, nil
}

func (srv pxlist) Find(ctx context.Context, pageText string) ([]string, error) {
	ips := []string{}
	for _, v := range srv.ipPattern.FindAllStringSubmatch(pageText, -1) {
		ips = append(ips, v[1])
	}
	//fmt.Println(ips)
	return ips, nil
}

func (srv pxlist) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}
