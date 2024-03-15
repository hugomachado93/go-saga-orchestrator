package kafka_handler

import (
	"main/internal/service/orchestrator_service"
	"main/internal/service/payment_service"
)

type KafkaListeners struct {
	orchestratorService *orchestrator_service.OrchestratorService
	paymentService      *payment_service.PaymentService
}

func NewKafkaListeners(os *orchestrator_service.OrchestratorService, ps *payment_service.PaymentService) *KafkaListeners {
	return &KafkaListeners{orchestratorService: os, paymentService: ps}
}

func (kl *KafkaListeners) HandleKafkaListeners() {
	listenKafkaHanlderFunc("teste3", "teste", 10, kl.orchestratorService.OrchestrateSaga)
	listenKafkaHanlderFunc("AUTHORIZE_PAYMENT", "commands", 10, kl.paymentService.AuthorizePayment)
}
