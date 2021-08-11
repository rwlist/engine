package domain

import (
	"go.etcd.io/bbolt"
)

type GlobalContext struct {
	ListFactory ListFactory
	DatabaseDir string
}

type DatabaseContext struct {
	*GlobalContext
	DatabaseName string
	Store        *bbolt.DB
}

type ListContext struct {
	DatabaseContext *DatabaseContext
	ListName        string
	ListBucket      *bbolt.Bucket
	CreatedNew      bool
}
