package user

import (
	"context"
	"fmt"
	"pet-store/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (username, first_name, last_name, email, password, phone, user_status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := ur.db.QueryRow(ctx, query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Phone,
		user.UserStatus,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, password, phone, user_status
		FROM users 
		WHERE username = $1
	`
	
	user := &entities.User{}
	err := ur.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.UserStatus,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	
	return user, nil
}


func (ur *UserRepository) Update(ctx context.Context, username string, user *entities.User) error {
	query := `
		UPDATE users 
		SET username = $1, first_name = $2, last_name = $3, email = $4, password = $5, phone = $6, user_status = $7
		WHERE  username = $8
	`

	_, err := ur.db.Exec(ctx, query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Phone,
		user.UserStatus,
		username,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, username string) error {
	query := `DELETE FROM users WHERE username = $1`
	
	_, err := ur.db.Exec(ctx, query, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}

func (ur *UserRepository) CreateMultiple(ctx context.Context, users []*entities.User) error {
	for i, user := range users {
		if err := ur.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user at index %d: %w", i, err)
		}
	}
	return nil
}