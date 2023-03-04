package repository

import (
	"context"

	"github.com/shusann01116/Chatting/backend/app/model"
)

// CreateRoom creates a new room
func (db Client) CreateRoom(ctx context.Context, room *model.Room) error {
	query := `INSERT INTO rooms (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := db.conn.Exec(ctx, query, room.ID, room.Name, room.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetRoomByID returns a room by its ID
func (db Client) GetRoomByID(ctx context.Context, roomID string) (*model.Room, error) {
	query := `SELECT id, name, created_at FROM rooms WHERE id = $1`
	room := &model.Room{}
	err := db.conn.QueryRow(ctx, query, roomID).Scan(&room.ID, &room.Name, &room.CreatedAt)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetAllRooms returns all rooms
func (db Client) GetAllRooms(ctx context.Context) ([]*model.Room, error) {
	query := `SELECT id, name, created_at FROM rooms`
	rows, err := db.conn.Query(ctx, query)
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

// UpdateRoom updates a room
func (db Client) UpdateRoom(ctx context.Context, room *model.Room) error {
	query := `UPDATE rooms SET name = $1 WHERE id = $2 RETURNING id, name, created_at`
	err := db.conn.QueryRow(ctx, query, room.Name, room.ID).Scan(&room.ID, &room.Name, &room.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRoom deletes a room
func (db Client) DeleteRoom(ctx context.Context, roomID string) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := db.conn.Exec(ctx, query, roomID)
	if err != nil {
		return err
	}
	return nil
}

// AddUserToRoom adds a user to a room
func (db Client) AddUserToRoom(ctx context.Context, roomID, userID string) error {
	query := `INSERT INTO room_users (room_id, user_id) VALUES ($1, $2)`
	_, err := db.conn.Exec(ctx, query, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserFromRoom deletes a user from a room
func (db Client) DeleteUserFromRoom(ctx context.Context, roomID, userID string) error {
	query := `DELETE FROM room_users WHERE room_id = $1 AND user_id = $2`
	_, err := db.conn.Exec(ctx, query, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}
