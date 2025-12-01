package repository

import (
	"JobPortal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
)

type JobTypeRepository interface {
	GetAllJobTypes(ctx context.Context) ([]*models.JobType, error)
	GetJobTypeById(ctx context.Context, id string) (*models.JobType, error)
	//storeJobTypes(jobTypes *models.JobTypes) error
}

type JobTypeRepositoryImpl struct {
	DB *sql.DB
}

func NewJobTypeRepository(db *sql.DB) JobTypeRepository {
	return &JobTypeRepositoryImpl{DB: db}
}

var table = (&models.JobType{}).TableName()

func (repo *JobTypeRepositoryImpl) GetAllJobTypes(ctx context.Context) ([]*models.JobType, error) {

	sel := sqlbuilder.Select(
		"id",
		"name",
		"slug",
		"is_active",
		"created_at",
		"updated_at",
		"deleted_at",
	).From(table).
		Where("deleted_at IS NULL").
		Limit(10)

	query, args := sel.BuildWithFlavor(sqlbuilder.PostgreSQL)

	fmt.Println("SQL:", query)
	fmt.Println("Args:", args)

	rows, err := repo.DB.QueryContext(ctx, query, args...)
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

	if result == nil {
		result = []*models.JobType{}
	}

	return result, nil
}

func (repo *JobTypeRepositoryImpl) GetJobTypeById(ctx context.Context, id string) (*models.JobType, error) {

	sel := sqlbuilder.NewSelectBuilder()
	sel.Select("*").
		From(table).
		Where(sel.Equal("id", id)).
		Limit(1)

	query, args := sel.BuildWithFlavor(sqlbuilder.PostgreSQL)

	fmt.Println("SQL:", query)
	fmt.Println("Args:", args)

	row := repo.DB.QueryRowContext(ctx, query, args...)
	var jt models.JobType
	var deletedAt sql.NullTime

	if err := row.Scan(
		&jt.Id,
		&jt.Name,
		&jt.Slug,
		&jt.IsActive,
		&jt.CreatedAt,
		&jt.UpdatedAt,
		&deletedAt,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if deletedAt.Valid {
		ts := deletedAt.Time
		jt.DeletedAt = &ts
	} else {
		jt.DeletedAt = nil
	}

	return &jt, nil
}
