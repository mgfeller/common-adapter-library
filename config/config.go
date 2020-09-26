package config

type Handler interface {
	SetKey(key string, value string)

	GetKey(key string) string

	Server(result interface{}) error

	MeshSpec(result interface{}) error

	MeshInstance(result interface{}) error

	Operations(result interface{}) error
}
