package rwimpl

import "github.com/rwlist/engine/pkg/domain"

type ListFactory struct {
	engines map[string]domain.ListConstructor
}

func NewListFactory(engines map[string]domain.ListConstructor) *ListFactory {
	return &ListFactory{
		engines: engines,
	}
}

func (f *ListFactory) Get(engine string) domain.ListConstructor {
	return f.engines[engine]
}
