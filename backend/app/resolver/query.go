package resolver

import (
	"context"

	"github.com/shusann01116/Chatting/backend/app/errors"
	"github.com/shusann01116/Chatting/backend/app/repository"
)

// The QueryResolver is the top level resolver for all top-level read operations.
type QueryResolver struct {
	client *repository.Client
}

func NewRoot(client *repository.Client) (*QueryResolver, error) {
	if client == nil {
		return nil, errors.ErrUnableToResolve
	}
	return &QueryResolver{client: client}, nil
}

// UsersQueryArgs is the arguments for the "user" query.
type UsersQueryArgs struct {
	// ID of the users to fetch. When nil, all users are fetched.
	ID *string
}

// Users resolvs a list of users. If no ID is provided, all users are returned.
func (r *QueryResolver) Users(ctx context.Context, args UsersQueryArgs) (*[]*UserResolver, error) {
	user, err := r.client.SearchUsers(ctx, strValue(args.ID))
	if err != nil {
		return nil, err
	}

	return NewUsers(ctx, NewUsersArgs{Page: user})
}
