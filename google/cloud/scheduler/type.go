package scheduler

import (
	"context"
	"google.golang.org/api/cloudscheduler/v1"
)

type IService interface {
	GetJobList(ctx context.Context) ([]*cloudscheduler.Job, error)
	GetJobListInLocation(ctx context.Context, location string) ([]*cloudscheduler.Job, error)
}

type serviceImpl struct {
	projectId        string
	schedulerService *cloudscheduler.Service
}
