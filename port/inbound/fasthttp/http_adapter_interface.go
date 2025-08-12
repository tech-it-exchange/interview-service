package fasthttp

import (
	"github.com/valyala/fasthttp"
)

type HttpAdapterInterface interface {
	GetHandler() fasthttp.RequestHandler
	InitRoutes()
}
