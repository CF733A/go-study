package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"repo/internal/entities/user"
	"repo/internal/storage"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *user.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6)	
	`
	_, err := u.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsDeleted)

	log.Println("user successfully created")
	return err
}

func (u *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var usr user.User
	query := `SELECT * FROM users WHERE id = $1 AND is_deleted = false`

	err := u.db.GetContext(ctx, &usr, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	return &usr, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error){
	var usr user.User
	query := `SELECT * FROM users WHERE email = $1 AND is_deleted = false`

	err := u.db.GetContext(ctx, &usr, query, email)
	if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    return &usr, nil
}

func (u *UserRepository) Update(ctx context.Context, user *user.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = $3
		WHERE id = $4 AND is_deleted = false
	`

	result, err := u.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.UpdatedAt,
		user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found or is deleted")
	}

	return nil
}

func (u *UserRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET deleted_at = $1, is_deleted = true, updated_at = $2
		WHERE id = $3 AND is_deleted = false
	`

	now := time.Now()
	result, err := u.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	log.Println("user successfully deleted")
	return nil
}

func (u *UserRepository) List(ctx context.Context, c storage.Conditions) ([]*user.User, int, error) {
	var users []*user.User
	var total int

	countQuery := `SELECT COUNT(*) FROM users WHERE is_deleted = false`
	err := u.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT * FROM users
		WHERE is_deleted = false
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	err = u.db.SelectContext(ctx, &users, query, c.Limit, c.Offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
