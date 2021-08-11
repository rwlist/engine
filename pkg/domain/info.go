package domain

type DatabaseInfo struct {
	Name string
}

type ListInfo struct {
	Name      string
	KVCount   int
	TreeDepth int

	// TODO: Open method, maybe?
}
