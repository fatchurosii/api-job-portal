package service

import (
	"context"

	"JobPortal/models"
	"JobPortal/repository"
)

type JobTypeService interface {
	GetAllJobTypes(ctx context.Context) ([]*models.JobType, error)
	GetJobTypeById(ctx context.Context, id string) (*models.JobType, error)
	StoreJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error)
	UpdateJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error)
	DeleteJobType(ctx context.Context, id string) error
	ChangeStatusJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error)
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

func (s *JobTypeServiceImpl) StoreJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error) {
	return s.JobTypeRepo.StoreJobType(ctx, jobType)
}

func (s *JobTypeServiceImpl) UpdateJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error) {
	return s.JobTypeRepo.UpdateJobType(ctx, jobType)
}

func (s *JobTypeServiceImpl) DeleteJobType(ctx context.Context, id string) error {
	return s.JobTypeRepo.DeleteJobType(ctx, id)
}

func (s *JobTypeServiceImpl) ChangeStatusJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error) {
	return s.JobTypeRepo.ChangeStatusJobType(ctx, jobType)
}
