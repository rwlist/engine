package mainlib

import (
	"github.com/rwlist/engine/pkg/domain"
	"github.com/rwlist/engine/pkg/rwimpl"
	"github.com/rwlist/engine/pkg/rwlists"
)

func StdFactory() domain.ListFactory {
	return rwimpl.NewListFactory(map[string]domain.ListConstructor{
		"dict": rwlists.NewDictionary,
	})
}
