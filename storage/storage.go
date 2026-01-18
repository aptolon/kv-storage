package storage

type Storage interface {
	Set(key string, value []byte) error
	Get(key string) (value []byte, err error)
	Delete(key string) error
}
