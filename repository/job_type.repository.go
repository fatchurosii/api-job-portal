package repository

import (
	"JobPortal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

type JobTypeRepository interface {
	GetAllJobTypes(ctx context.Context) ([]*models.JobType, error)
	GetJobTypeById(ctx context.Context, id string) (*models.JobType, error)
	StoreJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error)
	UpdateJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error)
}

type JobTypeRepositoryImpl struct {
	DB    *sql.DB
	table string
}

func NewJobTypeRepository(db *sql.DB) *JobTypeRepositoryImpl {
	return &JobTypeRepositoryImpl{
		DB:    db,
		table: (&models.JobType{}).TableName(),
	}
}

func (repo *JobTypeRepositoryImpl) GetAllJobTypes(ctx context.Context) ([]*models.JobType, error) {

	sel := sqlbuilder.Select(
		"id",
		"name",
		"slug",
		"is_active",
		"created_at",
		"updated_at",
		"deleted_at",
	).From(repo.table).
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
		From(repo.table).
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

func (repo *JobTypeRepositoryImpl) StoreJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error) {
	if jobType == nil {
		return nil, fmt.Errorf("jobType is nil")
	}

	jobType.PrepareForCreate()

	ins := sqlbuilder.NewInsertBuilder()
	ins.InsertInto(repo.table)
	ins.Cols("name", "slug", "is_active")
	ins.Values(jobType.Name, jobType.Slug, jobType.IsActive)

	query, args := ins.BuildWithFlavor(sqlbuilder.PostgreSQL)

	query = query + " RETURNING id, created_at, updated_at"

	fmt.Println("SQL:", query)
	fmt.Printf("Args: %#v\n", args)

	row := repo.DB.QueryRowContext(ctx, query, args...)

	var id string
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&id, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("insert succeeded but no row returned")
		}
		return nil, fmt.Errorf("insert job_type failed: %w", err)
	}

	jobType.Id = id
	jobType.CreatedAt = createdAt
	jobType.UpdatedAt = updatedAt

	return jobType, nil
}

func (repo *JobTypeRepositoryImpl) UpdateJobType(ctx context.Context, jobType *models.JobType) (*models.JobType, error) {
	if jobType == nil {
		return nil, fmt.Errorf("jobType is nil")
	}
	if jobType.Id == "" {
		return nil, fmt.Errorf("jobType id is empty")
	}

	jobType.PrepareForUpdate()

	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(repo.table)

	ub.Set(
		ub.Assign("name", jobType.Name),
		ub.Assign("slug", jobType.Slug),
	)

	ub.Where(ub.Equal("id", jobType.Id))

	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)

	query = query + " RETURNING id, name, slug, is_active, created_at, updated_at"

	fmt.Println("SQL:", query)
	fmt.Printf("Args: %#v\n", args)

	row := repo.DB.QueryRowContext(ctx, query, args...)

	var id string
	var name string
	var slug string
	var isActive bool
	var createdAt time.Time
	var updatedAt time.Time

	if err := row.Scan(&id, &name, &slug, &isActive, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("update succeeded but no row returned")
		}
		return nil, fmt.Errorf("update job_type failed: %w", err)
	}

	jobType.Id = id
	jobType.Name = name
	jobType.Slug = slug
	jobType.IsActive = isActive
	jobType.CreatedAt = createdAt
	jobType.UpdatedAt = updatedAt

	return jobType, nil
}
