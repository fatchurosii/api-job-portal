package service

import (
	"context"

	"JobPortal/models"
	"JobPortal/repository"
)

type JobTypeService interface {
	GetAllJobTypes(ctx context.Context) ([]*models.JobType, error)
	GetJobTypeById(ctx context.Context, id string) (*models.JobType, error)
}
type JobTypeServiceImpl struct {
	JobTypeRepo repository.JobTypeRepository
}

func NewJobTypeService(repo repository.JobTypeRepository) JobTypeService {
	return &JobTypeServiceImpl{JobTypeRepo: repo}
}

func (s *JobTypeServiceImpl) GetAllJobTypes(ctx context.Context) ([]*models.JobType, error) {
	return s.JobTypeRepo.GetAllJobTypes(ctx)
}

func (s *JobTypeServiceImpl) GetJobTypeById(ctx context.Context, id string) (*models.JobType, error) {
	return s.JobTypeRepo.GetJobTypeById(ctx, id)
}
