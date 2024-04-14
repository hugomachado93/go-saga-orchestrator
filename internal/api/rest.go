package api

import (
	"encoding/json"
	"main/internal/api/middlewares"
	"main/internal/api/requests"
	"main/internal/domain/saga"
	kafka_adapter "main/internal/kafka"
	"main/internal/services"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go/protocol"
)

type api struct {
	s  *services.SagaSettingsService
	am *middlewares.AuthMiddleware
}

func CreateRoutes(r *mux.Router, s *services.SagaSettingsService, am *middlewares.AuthMiddleware) {
	ap := &api{s: s, am: am}

	testesubr := r.Methods(http.MethodPost).PathPrefix("/v1/statemachine").Subrouter()
	testesubr.HandleFunc("/create", ap.createNewStateMachine)
	testesubr.Use(am.AuthMiddleware)
	testesubr.HandleFunc("/teste", func(w http.ResponseWriter, r *http.Request) {

		h := protocol.Header{Key: "x-api-key", Value: []byte("644e032b-3ee0-441e-a10c-d265a986ca2c")}

		headers := make([]protocol.Header, 0)
		headers = append(headers, h)

		rjson, _ := json.Marshal(saga.Response{SagaName: "teste", Payload: "", Event: "start"})
		kafka_adapter.SendMessage(string(rjson), "APP_ORCHESTRATOR", headers, "")
	})
}

func (a *api) createNewStateMachine(w http.ResponseWriter, r *http.Request) {
	var stm *requests.Statemachine
	xapk := r.Header.Get("x-api-key")
	json.NewDecoder(r.Body).Decode(&stm)
	err := a.s.InsertStateMachineSettings(stm.ToStateMachineSeetings(), xapk)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
