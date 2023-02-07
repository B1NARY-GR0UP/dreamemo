package protocol

type GetRequest interface {
	GetGroup() string
	GetKey() string
	String() string
}

type GetResponse interface {
	GetValue() []byte
	String() string
}
