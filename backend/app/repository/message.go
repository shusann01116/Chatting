package repository

import (
	"context"

	"github.com/shusann01116/Chatting/backend/app/model"
)

// CreateMessage creates a new message in the database
func (db Client) CreateMessage(ctx context.Context, message *model.Message) error {
	query := `INSERT INTO messages (id, room_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := db.conn.QueryRow(ctx, query, message.ID, message.RoomID, message.UserID, message.Content, message.CreatedAt).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetMessagesByRoomID returns all messages in a room
func (db Client) GetMessagesByRoomID(ctx context.Context, roomID string) ([]*model.Message, error) {
	query := `SELECT id, room_id, user_id, content, created_at FROM messages WHERE room_id = $1`
	rows, err := db.conn.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*model.Message{}
	for rows.Next() {
		message := &model.Message{}
		err := rows.Scan(&message.ID, &message.RoomID, &message.UserID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// GetMessageByID get a message in the database
func (db *Client) GetMessageByID(ctx context.Context, messageID string) (*model.Message, error) {
	query := `SELECT id, room_id, user_id, content, created_at FROM messages WHERE id = $1`

	message := &model.Message{}
	err := db.conn.QueryRow(ctx, query, messageID).Scan(&message.ID, &message.RoomID, &message.UserID, &message.Content, &message.CreatedAt)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// UpdateMessage updates a message in the database
func (db *Client) UpdateMessage(ctx context.Context, message *model.Message) error {
	query := `UPDATE messages SET message = $1 WHERE id = $2`
	_, err := db.conn.Exec(ctx, query, message.Content, message.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMessage deletes a message from the database
func (db *Client) DeleteMessage(ctx context.Context, messageID string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := db.conn.Exec(ctx, query, messageID)
	if err != nil {
		return err
	}
	return nil
}
