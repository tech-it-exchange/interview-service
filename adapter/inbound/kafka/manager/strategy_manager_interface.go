package manager

import (
	"interview-service/adapter/inbound/kafka/manager/abstract"
)

type StrategyManagerInterface interface {
	GetHandler(contractName string) abstract.EventStrategyInterface
}
