// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/B1NARY-GR0UP/dreamemo/app"
	"github.com/B1NARY-GR0UP/dreamemo/loadbalance"
	"github.com/B1NARY-GR0UP/dreamemo/protocol"
	"github.com/B1NARY-GR0UP/dreamemo/protocol/protobuf"
	pthrift "github.com/B1NARY-GR0UP/dreamemo/protocol/thrift"
	"github.com/apache/thrift/lib/go/thrift"
	"google.golang.org/protobuf/proto"
)

var _ loadbalance.Instance = (*Client)(nil)

const HTTPRequestMethod = http.MethodGet

type Client struct {
	Options   *app.Options
	BasePath  string
	Transport func(context.Context) http.RoundTripper
}

var defaultBufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func (c *Client) Get(ctx context.Context, in protocol.GetRequest, out protocol.GetResponse) error {
	requestURL := fmt.Sprintf("%v%v/%v", c.BasePath, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	req, err := http.NewRequest(HTTPRequestMethod, requestURL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	tpt := http.DefaultTransport
	if c.Transport != nil {
		tpt = c.Transport(ctx)
	}
	resp, err := tpt.RoundTrip(req)
	// defer resp.Body.Close() // nolint:errcheck
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error server response status: %v", resp.Status)
	}
	b := defaultBufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer defaultBufferPool.Put(b)
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	if c.Options.Thrift {
		if _, ok := in.(*pthrift.GetRequest); ok {
			deserializer := thrift.NewTDeserializer()
			err = deserializer.Read(ctx, out.(*pthrift.GetResponse), b.Bytes())
			if err != nil {
				return fmt.Errorf("err decoding thrift response body: %v", err)
			}
		}
	} else {
		if _, ok := in.(*protobuf.GetRequest); ok {
			err = proto.Unmarshal(b.Bytes(), out.(*protobuf.GetResponse))
			if err != nil {
				return fmt.Errorf("error decoding protobuf response body: %v", err)
			}
		}
	}
	return nil
}
