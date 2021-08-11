package domain

type CreateDatabaseRequest struct {
	Database string
}

type DropDatabaseRequest struct {
	Database string
}

type CreateListRequest struct {
	ListName string
	Engine   string
}

type DropListRequest struct {
	ListName string
}
