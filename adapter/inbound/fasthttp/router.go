package fasthttp

const (
	clientGroupPath = "/api/v2/client"

	lossPath                 = "/user/getLoss"
	createOrderPath          = "/order/create"
	createInstrumentPath     = "/instrument/create"
	userNamePath             = "/user/by/{name}"
	instrumentBalancePath    = "/instrument/balance/{instrumentId}"
	activeConsumerTopicsPath = "/kafka/topics"
)

// InitRoutes Инициализация маршрутов
func (r *Router) InitRoutes() {
	r.router.RedirectTrailingSlash = false

	clientGroup := r.router.Group(clientGroupPath)

	clientOrderController := r.controllers.ClientOrderController
	{
		clientGroup.POST(lossPath, clientOrderController.GetLoss)
		clientGroup.POST(createOrderPath, clientOrderController.CreateOrder)
		clientGroup.POST(createInstrumentPath, clientOrderController.CreateInstrument)
		clientGroup.GET(userNamePath, clientOrderController.GetUserName)
		clientGroup.GET(instrumentBalancePath, clientOrderController.GetInstrumentBalance)
		clientGroup.GET(activeConsumerTopicsPath, clientOrderController.GetActiveConsumerTopics)
	}
}
