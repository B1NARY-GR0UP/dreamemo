package eliminate

// ICore defines allowed operations of a memo
type ICore interface {
	Add(Key, Value)
	Get(Key) (Value, bool)
	Remove(Key)
	Clear()
	Name() string
}

type (
	Core struct {
		MaxSize   int // 0 => no limit
		UsedSize  int
		OnEvicted EvictFunc
	}
	EvictFunc func(key Key, value Value)
)

// Entity
type (
	Entity struct {
		Key   Key
		Value Value
	}
	Key   string
	Value any
)
