package scheduler

import (
	"context"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/google/cloud/location"
	"github.com/rosaekapratama/go-starter/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudscheduler/v1"
	"google.golang.org/api/option"
)

const spanJobs = "gcp.scheduler.Jobs"
const parent = "projects/%s/locations/%s"

var (
	Service IService
)

func Init(ctx context.Context, credentials *google.Credentials) {
	schedulerService, err := cloudscheduler.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatal(ctx, err, "Failed to create google cloud scheduler service")
		return
	}

	Service = &ServiceImpl{
		projectId:        credentials.ProjectID,
		schedulerService: schedulerService,
	}
	log.Info(ctx, "Google cloud scheduler client service is initiated")
}

// GetJobList will return job list in jakarta/asia-southeast2 location
func (s *ServiceImpl) GetJobList(ctx context.Context) ([]*cloudscheduler.Job, error) {
	return s.GetJobListInLocation(ctx, location.Jakarta)
}

func (s *ServiceImpl) GetJobListInLocation(ctx context.Context, location string) ([]*cloudscheduler.Job, error) {
	req := s.schedulerService.Projects.Locations.Jobs.List(fmt.Sprintf(parent, s.projectId, location))
	jobs := make([]*cloudscheduler.Job, integer.Zero)
	if err := req.Pages(ctx, func(page *cloudscheduler.ListJobsResponse) error {
		jobs = append(jobs, page.Jobs...)
		for page.NextPageToken != str.Empty {
			req.PageToken(page.NextPageToken)
			if innerErr := req.Pages(ctx, func(innerPage *cloudscheduler.ListJobsResponse) error {
				page = innerPage
				jobs = append(jobs, page.Jobs...)
				return nil
			}); innerErr != nil {
				log.Errorf(ctx, innerErr, "Failed to get google cloud scheduler job list")
				return innerErr
			}
		}
		return nil
	}); err != nil {
		log.Errorf(ctx, err, "failed to get google cloud scheduler job list")
		return nil, err
	}
	return jobs, nil
}
