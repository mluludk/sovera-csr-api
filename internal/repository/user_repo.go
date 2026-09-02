package repository

import (
	"context"
	"fmt"

	"sovera-core-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// FindByEmail retrieves a user by their email address (used for login).
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, org_id, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
		LIMIT 1;
	`
	var u model.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.OrgID, &u.Email, &u.PasswordHash,
		&u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

// FindByID retrieves a user by their UUID (used for /auth/me).
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, org_id, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
		LIMIT 1;
	`
	var u model.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.OrgID, &u.Email, &u.PasswordHash,
		&u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

// Create inserts a new user into the database and returns the created user.
func (r *UserRepository) Create(ctx context.Context, u model.User) (*model.User, error) {
	query := `
		INSERT INTO users (org_id, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id, org_id, email, password_hash, full_name, role, is_active, created_at, updated_at;
	`
	var created model.User
	err := r.pool.QueryRow(ctx, query, u.OrgID, u.Email, u.PasswordHash, u.FullName, u.Role).Scan(
		&created.ID, &created.OrgID, &created.Email, &created.PasswordHash,
		&created.FullName, &created.Role, &created.IsActive, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &created, nil
}
