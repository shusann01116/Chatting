package resolver_test

import (
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/shusann01116/Chatting/backend/app/resolver"
	"github.com/shusann01116/Chatting/backend/app/schema"
	"github.com/stretchr/testify/require"
)

func TestResolversSatisfySchema(t *testing.T) {
	s, err := schema.String()
	require.NoError(t, err)
	require.NotEmpty(t, s)

	rootResolver := &resolver.QueryResolver{}

	_, err = graphql.ParseSchema(s, rootResolver)
	require.NoError(t, err)
}
