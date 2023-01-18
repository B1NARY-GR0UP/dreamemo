package cache

// Operation defines allowed options to cache
type Operation interface {
	Add(Key, Value)
	Get(Key) (Value, bool)
	Remove(Key)
	Clear()
}

type (
	Cache struct {
		// MaxSize is the max numbers of entries the cache can take
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
