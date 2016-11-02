package database

import (
	"errors"
	"fmt"

	"github.com/graphql-go/relay"
	"gopkg.in/mgutz/dat.v1"
)

const OrderOnCreatedAt = "created_at"

var (
	DefaultLimit                  = 50
	ErrScopeInvalidBeforeAndAfter = errors.New("You cannot use before and after in the same query")
	ErrScopeInvalidFirstAndLast   = errors.New("You cannot use first and last in the same query")
)

type GraphQLCursor interface {
	Cursor() relay.ConnectionCursor
}

type GraphQLScope struct {
	relay.ConnectionArguments
	relay.ArraySliceMetaInfo
	Args    map[string]interface{}
	Limit   int
	Order   string
	OrderBy string
}

func NewGraphQLScope() GraphQLScope {
	filters := map[string]interface{}{}
	return NewGraphQLScopeWithFilters(filters)
}

func NewGraphQLScopeWithFilters(args map[string]interface{}) GraphQLScope {
	scope := GraphQLScope{
		Args:                args,
		ConnectionArguments: relay.NewConnectionArguments(args),
		Limit:               DefaultLimit,
		Order:               "ASC",
		OrderBy:             "id",
	}

	if scope.Args["order"] != nil {
		scope.Order = scope.Args["order"].(string)
	}

	return scope
}

func ApplyGraphQLScope(builder *dat.SelectBuilder, scope GraphQLScope) (*dat.SelectBuilder, error) {
	// Strongly discouraged from using both
	if scope.Before != "" && scope.After != "" {
		return nil, ErrScopeInvalidBeforeAndAfter
	}

	if scope.First != -1 && scope.Last != -1 {
		return nil, ErrScopeInvalidFirstAndLast
	}

	if scope.After != "" {
		after, err := relay.CursorToOffset(scope.After)
		if err != nil {
			return nil, err
		}

		if scope.Order == "DESC" {
			builder = builder.Where("id <= $1", after)
		} else {
			builder = builder.Where("id >= $1", after)
		}
	}

	if scope.First != -1 {
		scope.Limit = scope.First + 1
	}

	if scope.Before != "" {
		before, err := relay.CursorToOffset(scope.Before)
		if err != nil {
			return nil, err
		}

		if scope.Order == "DESC" {
			builder = builder.Where("id >= $1", before)
		} else {
			builder = builder.Where("id <= $1", before)
		}
	}

	if scope.Last != -1 {
		scope.Limit = scope.Last + 1

		// Flip the order of the results as we want to go in the reverse direction
		if scope.Order == "ASC" {
			scope.Order = "DESC"
		} else {
			scope.Order = "ASC"
		}
	}

	fmt.Printf("%+v", scope)
	if scope.Limit == 0 {
		builder = builder.Limit(uint64(DefaultLimit))
	} else if scope.Limit != -1 {
		builder = builder.Limit(uint64(scope.Limit))
	}

	builder = ApplyOrder(builder, scope)

	return builder, nil
}

func ApplyOrder(builder *dat.SelectBuilder, scope GraphQLScope) *dat.SelectBuilder {
	if scope.OrderBy != "" && scope.Order != "" {
		sql := fmt.Sprintf("%s %s", scope.OrderBy, scope.Order)
		builder = builder.OrderBy(sql)
	}

	return builder
}

func GraphQLConnection(arraySlice []GraphQLCursor, scope GraphQLScope) *relay.Connection {
	if scope.Limit == -1 {
		conn := relay.NewConnection()
		conn.PageInfo = relay.PageInfo{
			StartCursor:     "",
			EndCursor:       "",
			HasPreviousPage: false,
			HasNextPage:     false,
		}

		edges := []*relay.Edge{}
		for _, value := range arraySlice {
			edges = append(edges, &relay.Edge{
				Cursor: value.Cursor(),
				Node:   value,
			})
		}

		conn.Edges = edges

		return conn
	}

	var startCursor, endCursor relay.ConnectionCursor
	args := scope.ConnectionArguments
	// Make sure we're within the bounds of
	limit := min(DefaultLimit, len(arraySlice))
	begin, end := 0, limit

	if args.First == -1 && args.Last == -1 {
		// There are more pages
		if len(arraySlice) == limit+1 {
			// We don't want to grab the last edge that was used for pagination
			end = end - 1
			endCursor = arraySlice[end].Cursor()
		}
	} else if args.First != -1 {
		// There are more pages
		if len(arraySlice) == args.First+1 {
			// We don't want to grab the last edge that was used for pagination
			end = end - 1
			endCursor = arraySlice[args.First].Cursor()
		}
	} else if args.Last != -1 {
		// There are more pages
		if len(arraySlice) == args.Last+1 {
			end = args.Last
			startCursor = arraySlice[args.Last].Cursor()
		}
	}

	slice := arraySlice[begin:end]

	edges := []*relay.Edge{}
	for _, value := range slice {
		edges = append(edges, &relay.Edge{
			Cursor: value.Cursor(),
			Node:   value,
		})
	}

	conn := relay.NewConnection()
	conn.Edges = edges
	conn.PageInfo = relay.PageInfo{
		StartCursor:     startCursor,
		EndCursor:       endCursor,
		HasPreviousPage: startCursor != "",
		HasNextPage:     endCursor != "",
	}

	return conn
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
