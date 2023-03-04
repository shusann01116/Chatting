package repository

import (
	"context"

	"github.com/shusann01116/Chatting/backend/app/model"
)

type UserPage struct {
	Count int64        `json:"count"`
	Users []model.User `json:"results"`
}

func (user UserPage) IDs() []string {
	ids := make([]string, len(user.Users))
	for i, user := range user.Users {
		ids[i] = user.ID
	}
	return ids
}

func (db Client) User(ctx context.Context, id string) (model.User, error) {
	if id == "" {
		return model.User{}, nil
	}

	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	user := model.User{}

	// Query and return the result
	err := db.Exec(
		QueryArg{Sql: query, Args: []interface{}{id}},
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (db Client) SearchUsers(ctx context.Context, id string) (UserPage, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	user := model.User{}
	users := []model.User{}

	// Query and return the result
	err := db.Exec(
		QueryArg{Sql: query, Args: []interface{}{id}},
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	users = append(users, user)

	if err != nil {
		return UserPage{}, err
	}

	result := UserPage{
		Count: int64(len(users)),
		Users: users,
	}

	return result, nil
}
