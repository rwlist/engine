package rwimpl

import (
	"fmt"
	"github.com/rwlist/engine/pkg/domain"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"os"
)

const engineKey = "_engine"

type Database struct {
	ctx   *domain.DatabaseContext
	store *bbolt.DB
}

func NewDatabase(ctx *domain.DatabaseContext) *Database {
	return &Database{
		ctx:   ctx,
		store: ctx.Store,
	}
}

func (d *Database) AllLists() ([]domain.ListInfo, error) {
	var res []domain.ListInfo

	err := d.store.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			list, err := d.listFromBucket(name, b)
			if err != nil {
				log.WithError(err).WithField("name", string(name)).Error("failed to init list from bucket")
			}

			res = append(res, list.Info())
			return nil
		})
	})

	return res, err
}

func (d *Database) CreateList(req *domain.CreateListRequest) (domain.ListInfo, error) {
	var info domain.ListInfo

	err := d.store.Update(func(tx *bbolt.Tx) error {
		list, err := d.createList(tx, req)
		if err != nil {
			return err
		}

		info = list.Info()
		return nil
	})

	return info, err
}

func (d *Database) DropList(req *domain.DropListRequest) error {
	err := d.store.Update(func(tx *bbolt.Tx) error {
		list, err := d.listByName(tx, req.ListName)
		if err != nil {
			return err
		}

		droppable, ok := list.(domain.DroppableList)
		if ok {
			err = droppable.DropList()
			if err != nil {
				return err
			}
		}

		return tx.DeleteBucket([]byte(req.ListName))
	})

	return err
}

func (d *Database) listByName(tx *bbolt.Tx, nameStr string) (domain.List, error) {
	name := []byte(nameStr)
	b := tx.Bucket(name)
	if b == nil {
		return nil, domain.ErrListNotFound
	}

	return d.listFromBucket(name, b)
}

func (d *Database) listFromBucket(name []byte, b *bbolt.Bucket) (domain.List, error) {
	engine := b.Get([]byte(engineKey))
	if engine == nil {
		return nil, fmt.Errorf("engine field missing: %w", domain.ErrInvalidList)
	}

	constructor := d.ctx.ListFactory.Get(string(engine))
	if constructor == nil {
		return nil, fmt.Errorf("unknown engine %s: %w", engine, domain.ErrInvalidList)
	}

	ctx := &domain.ListContext{
		DatabaseContext: d.ctx,
		ListName:        string(name),
		ListBucket:      b, // TODO: possible to embed deeper
		CreatedNew:      false,
	}

	return constructor(ctx)
}

func (d *Database) createList(tx *bbolt.Tx, req *domain.CreateListRequest) (domain.List, error) {
	constructor := d.ctx.ListFactory.Get(req.Engine)
	if constructor == nil {
		return nil, fmt.Errorf("unknown engine %s: %w", req.Engine, domain.ErrInvalidList)
	}

	name := []byte(req.ListName)
	b, err := tx.CreateBucket(name)
	if err != nil {
		return nil, err
	}

	err = b.Put([]byte(engineKey), []byte(req.Engine))
	if err != nil {
		return nil, err
	}

	ctx := &domain.ListContext{
		DatabaseContext: d.ctx,
		ListName:        string(name),
		ListBucket:      b, // TODO: possible to embed deeper
		CreatedNew:      true,
	}

	return constructor(ctx)
}

func (d *Database) OpenList(name string, f func(domain.List) error) error {
	return d.store.Update(func(tx *bbolt.Tx) error {
		list, err := d.listByName(tx, name)
		if err != nil {
			return err
		}

		return f(list)
	})
}

func (d *Database) DropDatabase() error {
	path := d.store.Path()

	err := d.Close()
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func (d *Database) Close() error {
	return d.store.Close()
}
