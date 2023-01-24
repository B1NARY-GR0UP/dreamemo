package distributed

type Distributed interface {
	Add(nodes ...string)
	Get(key string) string
	Name() string
}
