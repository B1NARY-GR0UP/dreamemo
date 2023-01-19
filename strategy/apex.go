package strategy

// Memo defines allowed options to strategy
type Memo interface {
	Add(Key, Value)
	Get(Key) (Value, bool)
	Remove(Key)
	Clear()
}

type (
	Core struct {
		// MaxSize is the max numbers of entries the strategy can take
		// Zero means no limit
		MaxSize   int
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
