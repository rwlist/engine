package domain

type List interface {
	Info() ListInfo
	Insert(entity Entity) error
	ReadRange(offset, limit int) ([]Entity, error)
}

type DroppableList interface {
	DropList() error
}

type CompleteList interface {
	List
	DroppableList
}
