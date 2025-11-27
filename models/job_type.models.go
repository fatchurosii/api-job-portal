package models

import (
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type JobType struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type JobTypeResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (t *JobType) TableName() string {
	return "job_types"
}

type CreateJobTypeInput struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"-"`
}

func NewJobTypeResponse(jobType *JobType) *JobTypeResponse {
	return &JobTypeResponse{
		Id:        jobType.Id,
		Name:      jobType.Name,
		Slug:      jobType.Slug,
		IsActive:  jobType.IsActive,
		CreatedAt: jobType.CreatedAt,
		UpdatedAt: jobType.UpdatedAt,
		DeletedAt: jobType.DeletedAt,
	}
}

func NewJobTypeResponses(jobTypes []*JobType) []*JobTypeResponse {
	responses := make([]*JobTypeResponse, len(jobTypes))
	for i, jt := range jobTypes {
		responses[i] = NewJobTypeResponse(jt)
	}
	return responses
}

func (t *JobType) PrepareForCreate() {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	if strings.TrimSpace(t.Slug) == "" {
		t.Slug = slug.Make(t.Name)
	} else {
		t.Slug = slug.Make(t.Slug)
	}
}

func (t *JobType) PrepareForUpdate() {
	t.UpdatedAt = time.Now().UTC()

	if strings.TrimSpace(t.Slug) == "" && strings.TrimSpace(t.Name) != "" {
		t.Slug = slug.Make(t.Name)
	} else if strings.TrimSpace(t.Slug) != "" {
		t.Slug = slug.Make(t.Slug)
	}
}
