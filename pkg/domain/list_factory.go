package domain

type ListConstructor func(ctx *ListContext) (List, error)

type ListFactory interface {
	Get(engine string) ListConstructor
}
