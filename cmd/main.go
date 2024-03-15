package main

import (
	"context"
	"encoding/json"
	"fmt"
	"main/config"
	db "main/internal/database"
	"main/internal/domain/saga"
	handlers "main/internal/handlers/kafka"
	"main/internal/repository"
	"main/internal/service/orchestrator_service"
	"main/internal/service/payment_service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

func main() {
	db := db.NewDBConn()
	sagaRepo := repository.NewSagaRepositor()
	tr := repository.NewTransaction(db)
	os := orchestrator_service.NewOrchestratorService(tr, sagaRepo)
	ps := payment_service.NewPaymentService()

	config.InitializeSagas()

	kh := handlers.NewKafkaListeners(os, ps)

	kh.HandleKafkaListeners()

	go func() {
		a := &saga.Response{SagaUUID: "123456", Payload: "mamas", Status: saga.SUCCESS, SagaName: "payment_saga"}
		ja, _ := json.Marshal(a)

		w := &kafka.Writer{
			Addr:     kafka.TCP("localhost:9092"),
			Topic:    "teste3",
			Balancer: &kafka.LeastBytes{},
		}

		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Value: ja,
			},
		)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("sent")

	}()

	r := mux.NewRouter()
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}
