package loadbalance

import (
	"context"
	"github.com/B1NARY-GR0UP/dreamemo/protocol/protobuf"
)

type LoadBalancer interface {
	Pick(key string) (Instance, bool)
}

type Instance interface {
	Get(ctx context.Context, in *protobuf.GetRequest, out *protobuf.GetResponse) error
}
