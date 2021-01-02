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

func TestIPWithIPv4(t *testing.T) {

	n := New()
	ip1 := net.ParseIP("127.0.0.1")
	v1 := newValue()
	err := n.Add(ip1, v1)
	if err != nil {
		t.Error(err)
	}

	ip2 := net.ParseIP("127.0.0.2")
	v2 := newValue()
	err = n.Add(ip2, v2)
	if err != nil {
		t.Error(err)
	}

	ip3 := net.ParseIP("192.168.1.0")
	v3 := newValue()
	err = n.Add(ip3, v3)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 3 {
		t.Error(fmt.Sprintf("expected 3 items, got %d", n.Len()))
	}

	b, err := n.Contains(ip1)
	if err != nil {
		t.Error(err)
	}

	if !b {
		t.Error("expected ip1 to be in store")
	}

	r, err := n.Get(ip1)
	if err != nil {
		t.Error(err)
	}

	if r[0] == nil {
		t.Error("expected ip1 to be in store")
	}

	if r[0] != v1 {
		t.Error(fmt.Sprintf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], v1))
	}

	r, err = n.Get(ip2)
	if err != nil {
		t.Error(err)
	}

	if r[0] == nil {
		t.Error("expected ip1 to be in store")
	}

	if r[0] != v2 {
		t.Error(fmt.Sprintf("retrieved r[0] (%#+v) does not equal v2 (%#+v)", r[0], v2))
	}

	r, err = n.Get(ip3)
	if err != nil {
		t.Error(err)
	}

	if r[0] == nil {
		t.Error("expected ip1 to be in store")
	}

	if r[0] != v3 {
		t.Error(fmt.Sprintf("retrieved r[0] (%#+v) does not equal v3 (%#+v)", r[0], v3))
	}

	r1, err := n.Remove(ip1)
	if err != nil {
		t.Error(err)
	}

	if r1 != v1 {
		t.Error(fmt.Sprintf("removed r1 (%#+v) does not equal v2 (%#+v)", r1, v2))
	}

	if n.Len() != 2 {
		t.Error(fmt.Sprintf("expected 2 items, got %d", n.Len()))
	}

	r2, err := n.Remove(ip2)
	if err != nil {
		t.Error(err)
	}

	if r2 != v2 {
		t.Error(fmt.Sprintf("removed r2 (%#+v) does not equal v2 (%#+v)", r2, v2))
	}

	r3, err := n.Remove(ip3)
	if err != nil {
		t.Error(err)
	}

	if r3 != v3 {
		t.Error(fmt.Sprintf("removed r3 (%#+v) does not equal v3 (%#+v)", r3, v3))
	}

	if n.Len() != 0 {
		t.Error(fmt.Sprintf("expected 0 items, got %d", n.Len()))
	}
}

func TestCIDRWithIPv4(t *testing.T) {
	n := New()
	_, cidr1, _ := net.ParseCIDR("192.168.0.1/24")
	v1 := newValue()
	err := n.AddCIDR(*cidr1, v1)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 1 {
		t.Error(fmt.Sprintf("expected length to be 1; got: %d", n.Len()))
	}

	r, err := n.GetCIDR(*cidr1)
	if err != nil {
		t.Error(err)
	}

	if r[0] != v1 {
		t.Error(fmt.Sprintf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], v1))
	}

	_, cidr2, _ := net.ParseCIDR("192.168.1.1/24")

	r2, err := n.GetCIDR(*cidr2)
	if err != nil {
		t.Error(err)
	}

	if len(r2) != 0 {
		t.Error("retrieved CIDR that should not be retrieved")
	}

	v2 := newValue()
	err = n.AddCIDR(*cidr2, v2)
	if err != nil {
		t.Error(err)
	}

	r2, err = n.GetCIDR(*cidr2)
	if err != nil {
		t.Error(err)
	}

	if len(r2) == 0 {
		t.Error("did  not retrieve CIDR that should not be retrieved")
	}

	if r2[0] != v2 {
		t.Error(fmt.Sprintf("retrieved r2[0] (%#+v) does not equal v2 (%#+v)", r2[0], v2))
	}

	// TODO: this returns both of the networks before, because it's currently not an exact match, but including
	// the smaller CIDRs that are below (i.e. 192.168.0.1/24 and 192.168.1.1/24). It is useful to have, but not
	// the right thing for lookup cases that act upon the actual key.
	_, cidr3, _ := net.ParseCIDR("192.168.0.1/16")
	r3, err := n.GetCIDR(*cidr3)
	if err != nil {
		t.Error(err)
	}

	if len(r3) != 0 {
		t.Error(fmt.Sprintf("expected length to be 0; got: %d", len(r3)))
	}

	v3 := newValue()
	err = n.AddCIDR(*cidr3, v3)

	r3, err = n.GetCIDR(*cidr3)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(r3)

	if len(r3) != 1 {
		t.Error(fmt.Sprintf("expected length to be 1; got: %d", len(r3)))
	}

	if r3[0] != v3 {
		t.Error(fmt.Sprintf("retrieved r3[0] (%#+v) does not equal v3 (%#+v)", r3[0], v3))
	}

	rr3, err := n.RemoveCIDR(*cidr3)
	if err != nil {
		t.Error(err)
	}

	if rr3 != v3 {
		t.Error(fmt.Sprintf("removed rr3(%#+v) does not equal v3 (%#+v)", rr3, v3))
	}
}
