package handlers

import (
	"main/internal/infrastructure/kfk"
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
	kfk.ListenKafkaHanlderFunc("APP_ORCHESTRATOR", "orquestrator", 0, kl.orchestratorService.OrchestrateSaga)
	kfk.ListenKafkaHanlderFunc("AUTHORIZE_PAYMENT", "commands", 0, kl.paymentService.AuthorizePayment)
}
