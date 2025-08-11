package fasthttp

import (
	"log"

	fastHttpLib "github.com/valyala/fasthttp"
)

func NewHandler(logger *log.Logger, h fastHttpLib.RequestHandler) fastHttpLib.RequestHandler {
	return loggerMiddleware(logger, errorHandler(logger, h))
}

func loggerMiddleware(logger *log.Logger, h fastHttpLib.RequestHandler) fastHttpLib.RequestHandler {
	return func(ctx *fastHttpLib.RequestCtx) {
		logger.Printf("Request %s %s %s", ctx.Method(), ctx.Host(), ctx.RequestURI())
		h(ctx)
	}
}

func errorHandler(logger *log.Logger, h fastHttpLib.RequestHandler) fastHttpLib.RequestHandler {
	return func(ctx *fastHttpLib.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("%s", r)

				ctx.SetStatusCode(fastHttpLib.StatusInternalServerError)
				ctx.SetBodyString("Internal Server Error")
			}
		}()
		h(ctx)
	}
}
