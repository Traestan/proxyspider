package container

import (
	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"gitlab.com/likekimo/goproxyspider/spiders"
)

func New(logger *logger.Logger, proxys chan models.Proxy) []spiders.Spider {
	var Spiders = []spiders.Spider{
		spiders.AlivepxSpider(logger, proxys),
		spiders.FoxtoolpxSpider(logger, proxys),
		spiders.TxtpxSpider(logger, proxys),
		spiders.PxlistpxSpider(logger, proxys),
		spiders.FreepxlistSpider(logger, proxys),
		spiders.IpaddrespxSpider(logger, proxys),
		spiders.OpnproxySpider(logger, proxys),
		spiders.ProxzSpider(logger, proxys),
		spiders.NntimeSpider(logger, proxys),
		spiders.MiniproxSpider(logger, proxys),
		spiders.ProxbugSpider(logger, proxys),
		spiders.KuaidailiSpider(logger, proxys),
		spiders.HexmagicSpider(logger, proxys),
		spiders.ScyllaSpider(logger, proxys),
	}
	return Spiders
}
