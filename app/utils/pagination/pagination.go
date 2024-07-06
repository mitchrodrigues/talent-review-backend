package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"gorm.io/gorm"
)

// Options to configure pagination
type Options[T any] struct {
	Limit  int
	Model  []T
	Cursor string
	Scopes []func(*gorm.DB) *gorm.DB
}

// CursorPagination struct to encapsulate pagination logic and hold results
type CursorPagination[T any] struct {
	Options[T]

	Edges []Edge[T]

	TotalCount int64
	HasNext    bool
	HasPrev    bool

	NextCursor string
	PrevCursor string
}

// Edge struct to represent an edge in GraphQL pagination
type Edge[T any] struct {
	Cursor string
	Node   T
}

func NewCursorPaginationFromArgs[T any](args map[string]interface{}, model []T, scopes ...func(*gorm.DB) *gorm.DB) *CursorPagination[T] {
	options := PaginationOptionsFromArgs(args, model)

	if len(scopes) > 0 {
		options.Scopes = append(options.Scopes, scopes...)
	}

	return NewCursorPagination(options)
}

// NewCursorPagination function to create a new CursorPagination instance
func NewCursorPagination[T any](opts Options[T]) *CursorPagination[T] {
	return &CursorPagination[T]{Options: opts}
}

func (cp *CursorPagination[T]) SetScopes(scopes ...func(*gorm.DB) *gorm.DB) *CursorPagination[T] {
	if len(scopes) > 0 {
		cp.Scopes = append(cp.Scopes, scopes...)
	}

	return cp
}

// Paginate method to perform pagination and store results in the struct
func (cp *CursorPagination[T]) Paginate(gctx golly.Context) (*CursorPagination[T], error) {
	var totalCount int64

	// Apply scopes
	query := orm.DB(gctx).Model(&cp.Model).Scopes(cp.Scopes...)

	// Count total items
	if err := query.Count(&totalCount).Error; err != nil {
		return cp, err
	}

	// Determine offset from cursor
	offset, err := decodeCursor(cp.Options.Cursor)
	if err != nil {
		return cp, errors.New("invalid cursor")
	}

	// Apply limit and offset
	if err := query.Limit(cp.Limit).Offset(offset).Order("created_at DESC").Find(&cp.Model).Error; err != nil {
		return cp, err
	}

	// Prepare edges
	for i, item := range cp.Model {
		cp.Edges = append(cp.Edges, Edge[T]{
			Cursor: encodeCursor(offset + i),
			Node:   item,
		})
	}

	// Check if there are more items
	cp.HasNext = (offset + cp.Limit) < int(totalCount)
	cp.HasPrev = offset > 0
	cp.TotalCount = totalCount

	// Set next and previous cursors
	if cp.HasNext {
		cp.NextCursor = encodeCursor(offset + cp.Limit)
	}
	if cp.HasPrev {
		cp.PrevCursor = encodeCursor(offset - cp.Limit)
	}

	return cp, nil
}

// encodeCursor encodes an integer offset into a cursor string
func encodeCursor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", offset)))
}

// decodeCursor decodes a cursor string into an integer offset
func decodeCursor(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}
	var offset int
	if _, err := fmt.Sscanf(string(decoded), "%d", &offset); err != nil {
		return 0, err
	}
	return offset, nil
}
