package spiders

import (
	"context"

	"gitlab.com/likekimo/goproxyspider/models"
)

type Spider interface {
	GetProxy(context.Context) ([]models.Proxy, error)
	Find(context.Context, string) ([]string, error)

	Run(context.Context)
}
