package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

var ErrUserNotFound = errors.New("user with this id does not exist")

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var users []modules.User
	err := r.db.DB.SelectContext(ctx, &users, "SELECT id, name, email, age, city FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CreateUser(u modules.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var id int
	err := r.db.DB.QueryRowxContext(
		ctx,
		`INSERT INTO users (name, email, age, city)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		u.Name, u.Email, u.Age, u.City,
	).Scan(&id)

	if err != nil {
		// пример обработки частого кейса: уникальный email
		if pgErr, ok := err.(*pq.Error); ok {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				return 0, errors.New("user with this email already exists")
			}
		}
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateUser(id int, u modules.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(
		ctx,
		`UPDATE users
		 SET name = $1, email = $2, age = $3, city = $4
		 WHERE id = $5`,
		u.Name, u.Email, u.Age, u.City, id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	var user modules.User
	err := r.db.DB.GetContext(ctx, &user, "SELECT id, name, email, age, city FROM users WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) DeleteUserByID(id int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rows == 0 {
		return 0, ErrUserNotFound
	}

	return rows, nil
}
