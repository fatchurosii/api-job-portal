package models

import (
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type Types struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type TypeResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (t *Types) TableName() string {
	return "types"
}

type CreateTypesInput struct {
	Name     string `json:"name" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
	IsActive bool   `json:"is_active"`
}

func NewTypesResponse(types *Types) *TypeResponse {
	return &TypeResponse{
		Id:        types.Id,
		Name:      types.Name,
		Slug:      types.Slug,
		IsActive:  types.IsActive,
		CreatedAt: types.CreatedAt,
		UpdatedAt: types.UpdatedAt,
		DeletedAt: types.DeletedAt,
	}
}

func NewUserResponses(types []Types) []*TypeResponse {
	responses := make([]*TypeResponse, len(types))
	for i, user := range types {
		responses[i] = NewTypesResponse(&user)
	}
	return responses
}

func (t *Types) beforeCreate() error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	if strings.TrimSpace(t.Name) == "" {
		t.Slug = slug.Make(t.Name)
	} else {
		t.Slug = slug.Make(t.Name)
	}

	return nil
}

func (t *Types) beforeUpdate() error {
	t.UpdatedAt = time.Now()

	if strings.TrimSpace(t.Slug) == "" && strings.TrimSpace(t.Name) != "" {
		t.Slug = slug.Make(t.Name)
	} else {
		t.Slug = slug.Make(t.Slug)
	}
	return nil
}
