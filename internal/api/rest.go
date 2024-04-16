package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"main/internal/api/middlewares"
	"main/internal/api/requests"
	"main/internal/domain/saga"
	"main/internal/infrastructure/kfk"
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

	smr := r.PathPrefix("/v1/statemachine").Subrouter()
	smr.HandleFunc("/check", ap.checkStateMachine).Methods(http.MethodGet)

	smrs := smr.NewRoute().Subrouter()
	smrs.Use(am.AuthMiddleware)
	smrs.HandleFunc("/create", ap.createNewStateMachine).Methods(http.MethodPost)

	smrs.HandleFunc("/teste", func(w http.ResponseWriter, r *http.Request) {

		h := protocol.Header{Key: "x-api-key", Value: []byte("644e032b-3ee0-441e-a10c-d265a986ca2c")}

		headers := make([]protocol.Header, 0)
		headers = append(headers, h)

		rjson, _ := json.Marshal(saga.Response{SagaName: "PAYMENT", Payload: "payload1", Event: "STARTED"})
		kfk.SendMessage(string(rjson), "APP_ORCHESTRATOR", headers, "")
	}).Methods(http.MethodPost)
}

func (a *api) createNewStateMachine(w http.ResponseWriter, r *http.Request) {
	var stm *requests.Statemachine
	xapk := r.Header.Get("x-api-key")
	json.NewDecoder(r.Body).Decode(&stm)
	err := a.s.InsertStateMachineSettings(stm.ToStateMachineSeetings(), xapk)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *api) checkStateMachine(w http.ResponseWriter, r *http.Request) {
	var stm *requests.Statemachine
	defBase64 := r.URL.Query().Get("definition")

	if defBase64 != "" {
		defJson, _ := base64.StdEncoding.DecodeString(defBase64)
		fmt.Println(string(defJson))
		json.Unmarshal(defJson, &stm)
	} else {
		json.NewDecoder(r.Body).Decode(&stm)
	}
	gsvg, err := stm.ToStateMachineSeetings().DrawGraph()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "image/svg+xml")
	w.WriteHeader(200)
	w.Write(gsvg)
}
