package rwimpl

import (
	"github.com/rwlist/engine/pkg/auth"
	"github.com/rwlist/engine/pkg/domain"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	dbSuffix = ".rwl"
	dbMode   = 0600
)

type DBMS struct {
	globalCtx *domain.GlobalContext
	dirPath   string
	loaded    map[string]domain.DatabaseInternals
	mutex     sync.RWMutex
}

func NewDBMS(globalCtx *domain.GlobalContext) (*DBMS, error) {
	d := &DBMS{
		globalCtx: globalCtx,
		dirPath:   globalCtx.DatabaseDir,
		loaded:    map[string]domain.DatabaseInternals{},
	}

	err := d.preloadAll()
	if err != nil {
		_ = d.Close()
		return nil, err
	}

	return d, nil
}

func (d *DBMS) preloadAll() error {
	entries, err := os.ReadDir(d.dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		filename := entry.Name()
		if !strings.HasSuffix(filename, dbSuffix) {
			// not a database
			continue
		}

		name := strings.TrimSuffix(filename, dbSuffix)
		_, err = d.loadDB(name)
		if err != nil {
			log.WithError(err).Error("db not loaded on start")
			continue
		}

		log.WithField("db", name).Info("db loaded on start")
	}

	return nil
}

func (d *DBMS) loadDB(name string) (domain.Database, error) {
	filename := name + dbSuffix
	boltDB, err := bbolt.Open(filepath.Join(d.dirPath, filename), dbMode, nil)
	if err != nil {
		return nil, err
	}

	ctx := &domain.DatabaseContext{
		GlobalContext: d.globalCtx,
		DatabaseName:  name,
		Store:         boltDB,
	}

	db := NewDatabase(ctx)
	d.loaded[name] = db
	return db, nil
}

func (d *DBMS) AllDatabases(user *auth.User) ([]domain.Database, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var res []domain.Database
	if !user.IsAdmin {
		db := d.loaded[user.Username]
		if db != nil {
			res = append(res, db)
		}

		return res, nil
	}

	for _, db := range d.loaded {
		res = append(res, db)
	}

	return res, nil
}

func (d *DBMS) CreateDatabase(user *auth.User, dbName string) (domain.Database, error) {
	if !user.IsAdmin {
		return nil, domain.ErrAccessDenied
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if db := d.loaded[dbName]; db != nil {
		return nil, domain.ErrDatabaseExists
	}

	db, err := d.loadDB(dbName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DBMS) DropDatabase(user *auth.User, dbName string) error {
	if !user.IsAdmin {
		return domain.ErrAccessDenied
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	db := d.loaded[dbName]
	if db == nil {
		return domain.ErrDatabaseNotFound
	}

	delete(d.loaded, dbName)

	return db.DropDatabase()
}

func (d *DBMS) Database(user *auth.User, name string) (domain.Database, error) {
	if !user.IsAdmin && name != user.Username {
		return nil, domain.ErrAccessDenied
	}

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	db := d.loaded[name]
	if db == nil {
		return nil, domain.ErrDatabaseNotFound
	}

	return db, nil
}

func (d *DBMS) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for name, db := range d.loaded {
		err := db.Close()
		if err != nil {
			log.WithError(err).WithField("db", name).Warn("failed to close")
		}
	}

	d.loaded = nil
	return nil
}
