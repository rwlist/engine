package mainlib

import (
	"encoding/json"
	"github.com/rwlist/engine/pkg/auth"
	"github.com/rwlist/engine/pkg/domain"
	"github.com/rwlist/engine/pkg/rwimpl"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var adminUser = &auth.User{IsAdmin: true}

type sampleStruct struct {
	ID    string `json:"_id"`
	Value string `json:"value"`
}

func TestDBSaveLoad(t *testing.T) {
	data := []sampleStruct{
		{
			ID:    "1",
			Value: "abc",
		},
		{
			ID:    "2",
			Value: "dce",
		},
		{
			ID:    "3",
			Value: "{}",
		},
		{
			ID:    "5",
			Value: "+-",
		},
		{
			ID:    "8",
			Value: "\"абв\"",
		},
	}

	dir, err := os.MkdirTemp(os.TempDir(), "rw_save")
	assert.NoError(t, err)

	dbms, err := rwimpl.NewDBMS(&domain.GlobalContext{
		ListFactory: StdFactory(),
		DatabaseDir: dir,
	})
	assert.NoError(t, err)

	dbName := "test"
	db, err := dbms.CreateDatabase(adminUser, dbName)
	assert.NoError(t, err)

	listName := "sample_dict"
	err = createListAndFill(db, listName, data)

	err = dbms.Close()
	assert.NoError(t, err)

	// load everything again
	dbms, err = rwimpl.NewDBMS(&domain.GlobalContext{
		ListFactory: StdFactory(),
		DatabaseDir: dir,
	})
	assert.NoError(t, err)

	db, err = dbms.Database(adminUser, dbName)
	assert.NoError(t, err)

	newData, err := readList(db, listName)
	assert.NoError(t, err)

	assert.Equal(t, data, newData)

	err = dbms.Close()
	assert.NoError(t, err)
}

func readList(db domain.Database, listName string) ([]sampleStruct, error) {
	var data []sampleStruct

	err := db.OpenList(listName, func(list domain.List) error {
		entities, err := list.ReadRange(-1, -1)
		if err != nil {
			return err
		}

		for _, e := range entities {
			var obj sampleStruct
			err = json.Unmarshal(e, &obj)
			if err != nil {
				return err
			}

			data = append(data, obj)
		}

		return nil
	})

	return data, err
}

func createListAndFill(db domain.Database, listName string, data []sampleStruct) error {
	_, err := db.CreateList(&domain.CreateListRequest{
		ListName: listName,
		Engine:   "dict",
	})
	if err != nil {
		return err
	}

	return db.OpenList(listName, func(list domain.List) error {
		for _, e := range data {
			body, err := json.Marshal(&e)
			if err != nil {
				return err
			}

			err = list.Insert(body)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
