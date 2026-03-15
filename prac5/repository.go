package main

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPaginatedUsers(page int, pageSize int) (PaginatedResponse, error) {

	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM users`

	var total int

	err := r.db.QueryRow(countQuery).Scan(&total)

	if err != nil {
		return PaginatedResponse{}, err
	}

	rows, err := r.db.Query(
		`SELECT id,name,email,gender,birth_date
		 FROM users
		 ORDER BY id
		 LIMIT $1 OFFSET $2`,
		pageSize,
		offset,
	)

	if err != nil {
		return PaginatedResponse{}, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {

		var u User

		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)

		users = append(users, u)
	}

	return PaginatedResponse{
		Data:       users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(u1 int, u2 int) ([]User, error) {

	query := `
SELECT u.id, u.name, u.email, u.gender, u.birth_date
FROM users u
JOIN user_friends f1 ON u.id = f1.friend_id
JOIN user_friends f2 ON u.id = f2.friend_id
WHERE f1.user_id = $1 AND f2.user_id = $2
`

	rows, err := r.db.Query(query, u1, u2)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {

		var u User

		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)

		users = append(users, u)
	}

	return users, nil
}
