package repository

import (
	"context"

	"github.com/shusann01116/Chatting/backend/app/model"
)

// CreateRoomUser creates a new room_user
func (db Client) CreateRoomUser(ctx context.Context, roomUser *model.RoomUser) error {
	query := `INSERT INTO room_users (id, room_id, user_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err := db.conn.Exec(ctx, query, roomUser.ID, roomUser.RoomID, roomUser.UserID, roomUser.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetRoomUserByID returns a room_user by its ID
func (db Client) GetRoomUserByID(ctx context.Context, userID string) (*model.RoomUser, error) {
	query := `SELECT id, room_id, user_id, created_at FROM room_users WHERE id = $1`
	roomUser := &model.RoomUser{}
	err := db.conn.QueryRow(ctx, query, userID).Scan(&roomUser.ID, &roomUser.RoomID, &roomUser.UserID, &roomUser.CreatedAt)
	if err != nil {
		return nil, err
	}
	return roomUser, nil
}

// GetUsersByRoomID returns users by its roomID
func (db Client) GetUsersByRoomID(ctx context.Context, roomID string) ([]*model.User, error) {
	query := `SELECT users.id, users.name, users.email, users.password, users.created_at FROM room_users JOIN users ON room_users.user_id = users.id WHERE room_users.room_id = $1`
	rows, err := db.conn.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*model.User{}
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetRoomsByUserID returns rooms by its userID
func (db Client) GetRoomsByUserID(ctx context.Context, userID string) ([]*model.Room, error) {
	query := `SELECT rooms.id, rooms.name, rooms.created_at FROM room_users JOIN rooms ON room_users.room_id = rooms.id WHERE room_users.user_id = $1`
	rows, err := db.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := []*model.Room{}
	for rows.Next() {
		room := &model.Room{}
		err := rows.Scan(&room.ID, &room.Name, &room.CreatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// DeleteRoomUser deletes a room_user by its ID
func (db Client) DeleteRoomUser(ctx context.Context, roomUserID string) error {
	query := `DELETE FROM room_users WHERE id = $1`
	_, err := db.conn.Exec(ctx, query, roomUserID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRoomUsersByRoomID deletes room_users by its roomID
func (db Client) DeleteRoomUsersByRoomID(ctx context.Context, roomID string) error {
	query := `DELETE FROM room_users WHERE room_id = $1`
	_, err := db.conn.Exec(ctx, query, roomID)
	if err != nil {
		return err
	}
	return nil
}
