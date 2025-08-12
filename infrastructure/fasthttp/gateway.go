package fasthttp

import (
	"fmt"
	fastHttpLib "github.com/valyala/fasthttp"
	"log/slog"
)

func NewHandler(logger *slog.Logger, h fastHttpLib.RequestHandler) fastHttpLib.RequestHandler {
	return loggerMiddleware(logger, errorHandler(logger, h))
}

func loggerMiddleware(logger *slog.Logger, h fastHttpLib.RequestHandler) fastHttpLib.RequestHandler {
	return func(ctx *fastHttpLib.RequestCtx) {
		logger.InfoContext(
			ctx,
			fmt.Sprintf("Request %s %s %s", ctx.Method(), ctx.Host(), ctx.RequestURI()),
		)
		h(ctx)
	}
}

func errorHandler(logger *slog.Logger, h fastHttpLib.RequestHandler) fastHttpLib.RequestHandler {
	return func(ctx *fastHttpLib.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(ctx, fmt.Sprintf("%s", r))

				ctx.SetStatusCode(fastHttpLib.StatusInternalServerError)
				ctx.SetBodyString("Internal Server Error")
			}
		}()
		h(ctx)
	}
}
