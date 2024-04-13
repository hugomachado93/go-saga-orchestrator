package main

import (
	"fmt"
	"main/internal/api"
	"main/internal/api/middlewares"
	db "main/internal/database"
	"main/internal/handlers"
	"main/internal/repository"
	"main/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db := db.NewDBConn()
	tr := repository.NewTransaction(db)
	sagaRepo := repository.NewSagaRepositor()
	stmr := repository.NewStatemachineRepository()

	ps := services.NewPaymentService()

	stms := services.NewStateMachine(stmr, tr)
	cr := repository.NewClientRepository()
	cs := services.NewClientService(cr, tr)
	am := middlewares.NewAuthmiddleware(cs)
	os := services.NewOrchestratorService(tr, sagaRepo, stmr)
	j := handlers.NewJobsHandler(os)

	kh := handlers.NewKafkaListeners(os, ps)
	kh.HandleKafkaListeners()
	j.HandleJobs()

	r := mux.NewRouter()

	api.CreateRoutes(r, stms, am)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}
