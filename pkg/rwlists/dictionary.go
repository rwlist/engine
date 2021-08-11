package rwlists

import (
	"encoding/json"
	"errors"
	"github.com/rwlist/engine/pkg/domain"
	"go.etcd.io/bbolt"
)

var (
	errFinished = errors.New("iteration is finished")
	dataKey     = []byte("data")
)

type idStruct struct {
	ID string `json:"_id"`
}

type Dictionary struct {
	ctx  *domain.ListContext
	data *bbolt.Bucket
}

func NewDictionary(ctx *domain.ListContext) (domain.List, error) {
	dict := &Dictionary{
		ctx: ctx,
	}

	var err error

	if ctx.CreatedNew {
		dict.data, err = ctx.ListBucket.CreateBucket(dataKey)
	} else {
		dict.data = ctx.ListBucket.Bucket(dataKey)
		if dict.data == nil {
			err = domain.ErrInvalidList
		}
	}

	if err != nil {
		return nil, err
	}

	return dict, nil
}

func (d *Dictionary) Info() domain.ListInfo {
	stats := d.data.Stats()

	return domain.ListInfo{
		Name:      d.ctx.ListName,
		KVCount:   stats.KeyN,
		TreeDepth: stats.Depth,
	}
}

func (d *Dictionary) Insert(entity domain.Entity) error {
	var id idStruct
	err := json.Unmarshal(entity, &id)
	if err != nil {
		return err
	}

	return d.data.Put([]byte(id.ID), entity)
}

func (d *Dictionary) ReadRange(offset, limit int) ([]domain.Entity, error) {
	var entities []domain.Entity

	err := d.data.ForEach(func(k, v []byte) error {
		if offset > 0 {
			offset--
			return nil
		}
		if limit == 0 {
			return errFinished
		}

		limit--
		entities = append(entities, v)
		return nil
	})
	if err == errFinished {
		err = nil
	}

	return entities, err
}
