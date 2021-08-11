package domain

import "github.com/rwlist/engine/pkg/auth"

type DBMS interface {
	AllDatabases(user *auth.User) ([]Database, error)
	CreateDatabase(user *auth.User, dbName string) (Database, error)
	DropDatabase(user *auth.User, dbName string) error
	Database(user *auth.User, dbName string) (Database, error)
	Close() error
}

type Database interface {
	AllLists() ([]ListInfo, error)
	CreateList(req *CreateListRequest) (*ListInfo, error)
	DropList(listName string) error
	OpenList(id string, f func(List) error) error

	Info() (*DatabaseInfo, error)
}

type DatabaseInternals interface {
	Database

	DropDatabase() error
	Close() error
}
