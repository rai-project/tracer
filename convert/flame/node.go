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

package types

import "encoding/json"

type Node struct {
	Name     string
	Value    int
	Children map[string]*Node
}

func (n *Node) Add(stackPtr *[]string, index int, value int) {
	n.Value += value
	if index >= 0 {
		head := (*stackPtr)[index]
		childPtr, ok := n.Children[head]
		if !ok {
			childPtr = &(Node{head, 0, make(map[string]*Node)})
			n.Children[head] = childPtr
		}
		childPtr.Add(stackPtr, index-1, value)
	}
}

func (n *Node) MarshalJSON() ([]byte, error) {
	v := make([]Node, 0, len(n.Children))
	for _, value := range n.Children {
		v = append(v, *value)
	}

	return json.Marshal(&struct {
		Name     string `json:"name"`
		Value    int    `json:"value"`
		Children []Node `json:"children"`
	}{
		Name:     n.Name,
		Value:    n.Value,
		Children: v,
	})
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



type Profile struct {
	RootNode Node
	Stack    []string
}

func (p *Profile) OpenStack() {
	p.Stack = []string{}
}

func (p *Profile) CloseStack() {
	p.RootNode.Add(&p.Stack, len(p.Stack)-1, 1)
	p.Stack = []string{}
}

func (p *Profile) AddFrame(name string) {
	re, _ := regexp.Compile(`^\(`) // Skip process names
	if !re.MatchString(name) {
		name = strings.Replace(name, ";", ":", -1) // replace ; with :
		name = strings.Replace(name, "<", "", -1)  // remove '<'
		name = strings.Replace(name, ">", "", -1)  // remove '>'
		name = strings.Replace(name, "\\", "", -1) // remove '\'
		name = strings.Replace(name, "\"", "", -1) // remove '"'
		if index := strings.Index(name, "("); index != -1 {
			name = name[:index] // delete everything after '('
		}
		p.Stack = append(p.Stack, name)
	}
}
