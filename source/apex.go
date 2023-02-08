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

package source

import "context"

// Getter specifies how to get data from a datasource
type Getter interface {
	Get(ctx context.Context, key string) ([]byte, error)
}

// GetterFunc uses the same concept as http.HandlerFunc
type GetterFunc func(ctx context.Context, key string) ([]byte, error)

// Get dat from datasource
func (f GetterFunc) Get(ctx context.Context, key string) ([]byte, error) {
	return f(ctx, key)
}
