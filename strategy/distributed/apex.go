package distributed

type Instance interface {
	Add(nodes ...string)
	Get(key string) string
	Name() string
}
