package repository

import (
	"JobPortal/models"
	"context"
	"database/sql"
)

type JobTypeRepository interface {
	GetAllJobTypes(ctx context.Context) ([]*models.JobType, error)
	//getJobTypeById(id string) (*models.JobTypes, error)
	//storeJobTypes(jobTypes *models.JobTypes) error
}

type JobTypeRepositoryImpl struct {
	DB *sql.DB
}

func NewJobTypeRepository(db *sql.DB) JobTypeRepository {
	return &JobTypeRepositoryImpl{DB: db}
}

func (repo *JobTypeRepositoryImpl) GetAllJobTypes(ctx context.Context) ([]*models.JobType, error) {
	const q = `
		SELECT id, name, slug, is_active, created_at, updated_at, deleted_at
		FROM job_types
		WHERE deleted_at IS NULL
	`

	rows, err := repo.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.JobType
	for rows.Next() {
		var jt models.JobType
		var deletedAt sql.NullTime

		if err := rows.Scan(
			&jt.Id,
			&jt.Name,
			&jt.Slug,
			&jt.IsActive,
			&jt.CreatedAt,
			&jt.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			ts := deletedAt.Time
			jt.DeletedAt = &ts
		} else {
			jt.DeletedAt = nil
		}

		result = append(result, &jt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		result = []*models.JobType{}
	}

	return result, nil
}
