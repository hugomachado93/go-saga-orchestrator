package handlers

import (
	"main/internal/services"

	"github.com/robfig/cron/v3"
)

type JobsHandler struct {
	o *services.OrchestratorService
}

func NewJobsHandler(o *services.OrchestratorService) *JobsHandler {
	return &JobsHandler{o: o}
}

func (j *JobsHandler) HandleJobs() {
	c := cron.New()

	// c.AddFunc("@every 10s", j.o.SendCommand)
	// c.AddFunc("@every 1m", j.o.CheckTimeout)

	c.Start()
}
