package resolver

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/shusann01116/Chatting/backend/app/errors"
	"github.com/shusann01116/Chatting/backend/app/loader"
	"github.com/shusann01116/Chatting/backend/app/model"
	"github.com/shusann01116/Chatting/backend/app/repository"
)

// UserResolver resolves the User type
type UserResolver struct {
	user model.User
}

type NewUsersArgs struct {
	Page repository.UserPage
	IDs  []string
}

type NewUserArgs struct {
	User model.User
	ID   string
}

func NewUser(ctx context.Context, args NewUserArgs) (*UserResolver, error) {
	var user model.User
	var err error

	switch {
	case args.User.ID != "":
		user = args.User
	case args.ID != "":
		user, err = loader.LoadUser(ctx, args.ID)
	default:
		err = errors.ErrUnableToResolve
	}

	if err != nil {
		return nil, err
	}

	return &UserResolver{user: user}, err
}

func NewUsers(ctx context.Context, args NewUsersArgs) (*[]*UserResolver, error) {
	err := loader.PrimeUsers(ctx, args.Page)
	if err != nil {
		return nil, err
	}

	results, err := loader.LoadUsers(ctx, append(args.IDs, args.Page.IDs()...))
	if err != nil {
		return nil, err
	}

	var users = results.WithoutErrors()
	var resolvers = make([]*UserResolver, 0, len(users))
	var errs errors.Errors

	for i, u := range users {
		r, err := NewUser(ctx, NewUserArgs{User: u})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		resolvers = append(resolvers, r)
	}

	return &resolvers, errs.Err()
}

// ID resolves the user's unique identifier.
func (r *UserResolver) ID() graphql.ID {
	return graphql.ID(r.user.ID)
}

// Name resolves the user's name.
func (r *UserResolver) Name() string {
	return r.user.Name
}

// ProfilePhoto resolves the user's profile photo.
func (r *UserResolver) ProfilePhoto() string {
	return r.user.ProfilePhoto
}

// Email resolves the user's profile photo.
func (r *UserResolver) Email() string {
	return r.user.Email
}

// CreatedAt resolves the user's profile photo.
func (r *UserResolver) CreatedAt() string {
	return r.user.CreatedAt.String()
}
