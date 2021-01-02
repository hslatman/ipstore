// Copyright 2021 Herman Slatman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipstore

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
)

type value struct {
	v int
}

func newValue() value {
	return value{
		v: rand.Int(),
	}
}

func TestNew(t *testing.T) {
	n := New()
	if n == nil {
		t.Fail()
	}
}

func TestAddIP(t *testing.T) {
	n := New()
	ip1 := net.ParseIP("127.0.0.1")
	v1 := newValue()
	n.Add(ip1, v1)
	ip2 := net.ParseIP("127.0.0.2")
	v2 := newValue()
	n.Add(ip2, v2)
	ip3 := net.ParseIP("192.168.1.0")
	v3 := newValue()
	n.Add(ip3, v3)

	if n.Len() != 3 {
		t.Fail()
	}

	b, err := n.Contains(ip1)
	if err != nil {
		t.Error(err)
	}
	if !b {
		t.Fail()
	}

	// e, err := n.Get(ip1)
	// if err != nil {
	// 	t.Error(err)
	// }

	// ev, ok := e.(entry)

	fmt.Println(v1, v2, v3)

	r, err := n.Get(ip1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(r)

	if r[0] == nil {
		t.Fail()
	}

	// n.Get(ip2)
	// n.Get(ip3)

	//t.Fail()
}

func TestAddCIDR(t *testing.T) {
	n := New()
	_, cidr1, _ := net.ParseCIDR("192.168.0.1/24")
	v1 := newValue()
	n.AddCIDR(*cidr1, v1)

	if n.Len() != 1 {
		t.Fail()
	}

	fmt.Println(cidr1, v1)

	r, err := n.GetCIDR(*cidr1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(r)

	_, cidr2, _ := net.ParseCIDR("192.168.1.1/24")
	fmt.Println(cidr2)

	r2, err := n.GetCIDR(*cidr2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(r2)

	//t.Fail()
}
