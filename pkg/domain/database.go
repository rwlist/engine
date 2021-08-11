package domain

import "github.com/rwlist/engine/pkg/auth"

type DBMS interface {
	AllDatabases(user *auth.User) ([]Database, error)
	CreateDatabase(user *auth.User, req *CreateDatabaseRequest) (Database, error)
	DropDatabase(user *auth.User, req *DropDatabaseRequest) error
	Database(user *auth.User, id string) (Database, error)
	Close() error
}

type Database interface {
	AllLists() ([]ListInfo, error)
	CreateList(req *CreateListRequest) (ListInfo, error)
	DropList(req *DropListRequest) error
	OpenList(id string, f func(List) error) error
}

type DatabaseInternals interface {
	Database

	DropDatabase() error
	Close() error
}
