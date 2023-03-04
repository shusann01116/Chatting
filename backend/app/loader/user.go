package loader

import (
	"context"
	"sync"

	"github.com/graph-gophers/dataloader"
	"github.com/shusann01116/Chatting/backend/app/errors"
	"github.com/shusann01116/Chatting/backend/app/model"
	"github.com/shusann01116/Chatting/backend/app/repository"
)

type UserResult struct {
	User  model.User
	Error error
}

type UserResults []UserResult

// WithoutErrors returns the Users does not have error
func (results UserResults) WithoutErrors() []model.User {
	users := make([]model.User, 0, len(results))

	for _, r := range results {
		if r.Error != nil {
			continue
		}

		users = append(users, r.User)
	}

	return users
}

func LoadUser(ctx context.Context, id string) (model.User, error) {
	var user model.User

	loader, err := extract(ctx, userLoaderKey)
	if err != nil {
		return user, err
	}

	data, err := loader.Load(ctx, dataloader.StringKey(id))()
	if err != nil {
		return user, err
	}

	user, ok := data.(model.User)
	if !ok {
		return user, errors.WrongType(user, data)
	}

	return user, nil
}

func LoadUsers(ctx context.Context, ids []string) (UserResults, error) {
	var results []UserResult
	loader, err := extract(ctx, userLoaderKey)
	if err != nil {
		return results, err
	}

	// Resolves the keys into values using the batch function
	data, errs := loader.LoadMany(ctx, dataloader.NewKeysFromStrings(ids))()
	results = make([]UserResult, 0, len(ids))

	for i, d := range data {
		var e error
		if errs != nil {
			e = errs[i]
		}

		user, ok := d.(model.User)
		if !ok {
			e = errors.WrongType(user, d)
		}

		results = append(results, UserResult{User: user, Error: e})
	}

	return results, nil
}

// Prime adds the provided key and value to the cache.
// If the key already exists, no change is made. Returns self for method chaining
func PrimeUsers(ctx context.Context, users repository.UserPage) error {
	loader, err := extract(ctx, userLoaderKey)
	if err != nil {
		return err
	}

	for _, u := range users.Users {
		loader.Prime(ctx, dataloader.StringKey(u.ID), u)
	}

	return nil
}

type userGetter interface {
	User(ctx context.Context, id string) (model.User, error)
}

// UserLoader contains the database required to load user resources.
type userLoader struct {
	get userGetter
}

func newUserLoader(db userGetter) dataloader.BatchFunc {
	return userLoader{get: db}.loadBatch
}

func (l userLoader) loadBatch(ctx context.Context, ids dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(ids)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, id := range ids {
		go func(i int, id dataloader.Key) {
			defer wg.Done()

			resp, err := l.get.User(ctx, id.String())
			results[i] = &dataloader.Result{Data: resp, Error: err}
		}(i, id)
	}

	wg.Wait()

	return results
}
