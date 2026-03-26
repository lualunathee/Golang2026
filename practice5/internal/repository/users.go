package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"assignment5/pkg/modules"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetPaginatedUsers(page, pageSize int, filters map[string]string, orderBy string) (modules.PaginatedResponse, error) {
	offset := (page - 1) * pageSize

	allowedOrder := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"gender":     true,
		"birth_date": true,
	}

	if !allowedOrder[orderBy] {
		orderBy = "id"
	}

	baseWhere := []string{}
	args := []interface{}{}
	argPos := 1

	if v := filters["id"]; v != "" {
		baseWhere = append(baseWhere, fmt.Sprintf("id = $%d", argPos))
		args = append(args, v)
		argPos++
	}
	if v := filters["name"]; v != "" {
		baseWhere = append(baseWhere, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+v+"%")
		argPos++
	}
	if v := filters["email"]; v != "" {
		baseWhere = append(baseWhere, fmt.Sprintf("email ILIKE $%d", argPos))
		args = append(args, "%"+v+"%")
		argPos++
	}
	if v := filters["gender"]; v != "" {
		baseWhere = append(baseWhere, fmt.Sprintf("gender = $%d", argPos))
		args = append(args, v)
		argPos++
	}
	if v := filters["birth_date"]; v != "" {
		baseWhere = append(baseWhere, fmt.Sprintf("birth_date = $%d", argPos))
		args = append(args, v)
		argPos++
	}

	whereSQL := ""
	if len(baseWhere) > 0 {
		whereSQL = " WHERE " + strings.Join(baseWhere, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM users" + whereSQL
	var totalCount int
	if err := r.db.Get(&totalCount, countQuery, args...); err != nil {
		return modules.PaginatedResponse{}, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date
		FROM users
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, argPos, argPos+1)

	args = append(args, pageSize, offset)

	var users []modules.User
	if err := r.db.Select(&users, query, args...); err != nil {
		return modules.PaginatedResponse{}, err
	}

	return modules.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *UserRepository) GetCommonFriends(user1, user2 int) ([]modules.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends f1 ON u.id = f1.friend_id
		JOIN user_friends f2 ON u.id = f2.friend_id
		WHERE f1.user_id = $1 AND f2.user_id = $2
		ORDER BY u.id
	`

	var users []modules.User
	if err := r.db.Select(&users, query, user1, user2); err != nil {
		return nil, err
	}

	return users, nil
}