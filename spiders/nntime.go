package spiders

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/parser"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"
)

type nntime struct {
	source    string //
	url       string
	ipPattern *regexp.Regexp
	logger    *logger.Logger
	proxys    []models.Proxy
	ch        chan models.Proxy
}

func NntimeSpider(logger *logger.Logger, gch chan models.Proxy) Spider {
	r, _ := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:[0-9]+")
	svc := &nntime{
		source: "nntime",
		url:    "http://nntime.com/proxy-list-0%d.htm",

		ipPattern: r,
		logger:    logger,
		ch:        gch,
	}

	return svc
}

func (srv nntime) GetProxy(ctx context.Context) ([]models.Proxy, error) {
	proxys := []models.Proxy{}
	//proxyPaths := make(map[string]string)
	matchesRow := make([]string, 1)
	codes := make(map[string]string)

	for i := 1; i <= 1; i++ {
		slug := fmt.Sprintf(srv.url, i)
		document, err := parser.GetContent(ctx, slug)
		if err != nil {
			srv.logger.Error("", zap.String("service", srv.source), zap.String("slug", slug), zap.Error(err))
		}

		pageContent, _ := document.Html()

		patternTable, _ := regexp.Compile(`<td>(.*?)</td>`)
		matchesRow = patternTable.FindAllString(pageContent, -1)

		codes, _ = srv.getHackCodes(pageContent) //srv.Find(ctx, pageContent)
	}

	for _, row := range matchesRow {
		pTmp, err := srv.decodeRow(codes, row)
		if err != nil {
			return nil, err
		}
		if len(pTmp) != 0 {
			proxy := models.Proxy{
				IP:     pTmp,
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

func (srv nntime) Find(ctx context.Context, pageText string) ([]string, error) {
	return srv.ipPattern.FindAllString(pageText, -1), nil
}

func (srv nntime) Run(ctx context.Context) {
	srv.GetProxy(ctx)
}

func (srv nntime) getHackCodes(htmlPage string) (map[string]string, error) {
	codes := make(map[string]string)
	parseCode, _ := regexp.Compile("((?:[a-z]=[0-9];)+)")
	matchesCode := parseCode.FindAllString(htmlPage, -1)
	code := matchesCode[0]
	//for _, code := range matchesCode {
	v := strings.Split(code, ";")
	for _, c := range v {
		tm := strings.Split(c, "=")
		if len(tm) == 2 {
			codes[tm[0]] = tm[1]
		}
	}
	return codes, nil
}
func (srv nntime) decodeRow(codes map[string]string, row string) (string, error) {
	parseIp, _ := regexp.Compile(`((?:[0-9]{1,3}\.){3}[0-9]{1,3})`)
	m := parseIp.FindAllString(row, -1)

	if len(m) == 0 {
		return "", nil
	}

	ip := m[0]
	parseCode, _ := regexp.Compile(`document\.write\(":"((?:\+[a-z]){0,6})`)
	pEncs := parseCode.FindAllString(row, -1)
	portTmp := make([]string, 1)
	if len(pEncs) > 0 {
		for _, pEnc := range pEncs {
			pDec := strings.Split(pEnc, "+")
			for _, c := range pDec {
				portTmp = append(portTmp, codes[c])
			}
		}

	} else {
		return "", nil
	}
	port := strings.Join(portTmp, "")
	return ip + ":" + port, nil
}
