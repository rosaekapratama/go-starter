package zeebe

import (
	"context"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/rosaekapratama/go-starter/log"
)

func StartWorker(jobType string, handler func(client worker.JobClient, job entities.Job), opts ...WorkerOption) {
	go startWorker(jobType, handler, opts...)
}

func startWorker(jobType string, handler func(client worker.JobClient, job entities.Job), opts ...WorkerOption) {
	ctx := context.Background()
	builder := Client.NewJobWorker().
		JobType(jobType).
		Handler(handler)
	for _, opt := range opts {
		opt.apply(builder)
	}

	builder.Open()
	log.Infof(ctx, "Zeebe worker started, jobType=%s", jobType)
}
