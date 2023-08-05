package repository

type Access interface {
	GetKeys() ([]string, error)
	DelKeys([]string) error
	Save()
}
