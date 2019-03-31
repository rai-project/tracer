package flame

// Copyright © 2017 Martin Spier <spiermar@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"encoding/json"
	"strings"

	"github.com/rai-project/tracer/convert"
)

type Node struct {
	ID       string            `json:"-,omitempty"`
	Name     string            `json:"name,omitempty"`
	Value    int               `json:"value,omitempty"`
	Interval *convert.Interval `json:"-"`
	Children []*Node           `json:"children,omitempty"`
}

func (n *Node) MarshalIndentJSON() ([]byte, error) {
	v := make([]Node, 0, len(n.Children))
	for _, value := range n.Children {
		v = append(v, *value)
	}

	return json.MarshalIndent(&struct {
		Name     string `json:"name"`
		Value    int    `json:"value"`
		Children []Node `json:"children"`
	}{
		Name:     n.Name,
		Value:    n.Value,
		Children: v,
	}, "", "  ")
}

func cleanName(name string) string {
	name = strings.Replace(name, ";", ":", -1) // replace ; with :
	name = strings.Replace(name, "<", "", -1)  // remove '<'
	name = strings.Replace(name, ">", "", -1)  // remove '>'
	name = strings.Replace(name, "\\", "", -1) // remove '\'
	name = strings.Replace(name, "\"", "", -1) // remove '"'

	return name
}
