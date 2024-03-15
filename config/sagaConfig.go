package config

import "main/internal/domain/saga"

func InitializeSagas() {
	saga.NewSagaDefinition("payment_saga").
		NextStep(saga.AUTHORIZE_PAYMENT).
		And().
		NextStep(saga.COMPLETE_PAYMENT).
		Build()
}
