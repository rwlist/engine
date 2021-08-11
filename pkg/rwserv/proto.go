package rwserv

import (
	"encoding/json"
	"github.com/rwlist/engine/pkg/domain"
)

// databases.getAll
type AllDatabasesResponse struct {
	Databases []domain.DatabaseInfo
}

// databases.create
type CreateDatabaseRequest struct {
	Database string
}

type CreateDatabaseResponse domain.DatabaseInfo

// databases.drop
type DropDatabaseRequest struct {
	Database string
}

type DropDatabaseResponse struct{}

// lists.getAll
type GetAllListsRequest struct {
	Database string
}

type GetAllListsResponse struct {
	Lists []domain.ListInfo
}

// lists.create
type CreateListRequest struct {
	Database string
	ListName string
	Engine   string
}

type CreateListResponse domain.ListInfo

// lists.drop
type DropListRequest struct {
	Database string
	ListName string
}

type DropListResponse struct{}

// list.insertMany
type InsertManyRequest struct {
	Database string
	ListName string
	Entries  []json.RawMessage
}

type InsertManyResponse struct{}

// list.readRange
type ReadRangeRequest struct {
	Database string
	ListName string
	Offset   int
	Limit    int
}

type ReadRangeResponse struct {
	Entries []json.RawMessage
}
