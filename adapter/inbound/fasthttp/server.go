package fasthttp

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"interview-service/adapter/inbound/fasthttp/controllers"
	outboundFasthttp "interview-service/port/inbound/fasthttp"
)

type Router struct {
	router      *router.Router
	userGroup   *router.Group
	controllers *controllers.Controllers
}

func NewHttpAdapter(
	dependencies *controllers.Dependencies,
) (outboundFasthttp.HttpAdapterInterface, error) {
	newRouter := router.New()
	newControllers, err := controllers.NewControllers(dependencies)
	if err != nil {
		return nil, err
	}

	return &Router{
		router:      newRouter,
		controllers: newControllers,
	}, nil
}

func (r *Router) GetHandler() fasthttp.RequestHandler {
	return r.router.Handler
}
