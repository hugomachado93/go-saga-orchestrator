package handlers

import (
	kafka_adapter "main/internal/kafka"
	"main/internal/services"
)

type KafkaListeners struct {
	orchestratorService *services.OrchestratorService
	paymentService      *services.PaymentService
}

func NewKafkaListeners(os *services.OrchestratorService, ps *services.PaymentService) *KafkaListeners {
	return &KafkaListeners{orchestratorService: os, paymentService: ps}
}

func (kl *KafkaListeners) HandleKafkaListeners() {
	kafka_adapter.ListenKafkaHanlderFunc("APP_ORCHESTRATOR", "orquestrator", 0, kl.orchestratorService.OrchestrateSaga)
	kafka_adapter.ListenKafkaHanlderFunc("AUTHORIZE_PAYMENT", "commands", 0, kl.paymentService.AuthorizePayment)
}
