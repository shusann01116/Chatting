package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/shusann01116/Chatting/backend/app/errors"
)

type key string

const (
	userLoaderKey    key = "userLoader"
	roomLoaderKey    key = "roomLoader"
	messageLoaderKey key = "messageLoader"
)

type DataBase interface {
	userGetter
	// roomGetter
	// messageGetter
}

func Initialize(db DataBase) Collection {
	return Collection{
		lookup: map[key]dataloader.BatchFunc{
			userLoaderKey: newUserLoader(db),
			// roomLoaderKey:    newRoomLoader(db),
			// messageLoaderKey: newMessageLoader(db),
		},
	}
}

// Collection holds an internal lookup of initialized baatch data load functions.
type Collection struct {
	lookup map[key]dataloader.BatchFunc
}

// Attach creates new instances of dataloder.Loader and attaches the instances on the request context.
func (c Collection) Attach(ctx context.Context) context.Context {
	for k, batchFn := range c.lookup {
		ctx = context.WithValue(ctx, k, dataloader.NewBatchedLoader(batchFn))
	}

	return ctx
}

// extract is a helper function to extract a batch function from the internal lookup.
func extract(ctx context.Context, k key) (*dataloader.Loader, error) {
	ldr, ok := ctx.Value(k).(*dataloader.Loader)
	if !ok {
		return nil, errors.ErrLoaderNotFound
	}

	return ldr, nil
}
