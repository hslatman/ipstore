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

package ipstore_test

import (
	"math/rand"
	"net"
	"sync"
	"testing"

	"github.com/hslatman/ipstore"
)

type value struct {
	v int
}

func newValue() value {
	return value{
		v: rand.Int(),
	}
}

func hosts(cidr string) ([]net.IP, int, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, 0, err
	}

	var ips []net.IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, net.ParseIP(ip.String()))
	}

	return ips, len(ips), nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func TestNew(t *testing.T) {
	n := ipstore.New()
	if n == nil {
		t.Fail()
	}
}

func TestIPWithIPv4(t *testing.T) {
	n := ipstore.New()
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
		t.Errorf("expected 3 items, got %d", n.Len())
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
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], v1)
	}

	r, err = n.Get(ip2)
	if err != nil {
		t.Error(err)
	}

	if r[0] == nil {
		t.Error("expected ip1 to be in store")
	}

	if r[0] != v2 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v2 (%#+v)", r[0], v2)
	}

	r, err = n.Get(ip3)
	if err != nil {
		t.Error(err)
	}

	if r[0] == nil {
		t.Error("expected ip1 to be in store")
	}

	if r[0] != v3 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v3 (%#+v)", r[0], v3)
	}

	r1, err := n.Remove(ip1)
	if err != nil {
		t.Error(err)
	}

	if r1 != v1 {
		t.Errorf("removed r1 (%#+v) does not equal v2 (%#+v)", r1, v2)
	}

	if n.Len() != 2 {
		t.Errorf("expected 2 items, got %d", n.Len())
	}

	r2, err := n.Remove(ip2)
	if err != nil {
		t.Error(err)
	}

	if r2 != v2 {
		t.Errorf("removed r2 (%#+v) does not equal v2 (%#+v)", r2, v2)
	}

	r3, err := n.Remove(ip3)
	if err != nil {
		t.Error(err)
	}

	if r3 != v3 {
		t.Errorf("removed r3 (%#+v) does not equal v3 (%#+v)", r3, v3)
	}

	if n.Len() != 0 {
		t.Errorf("expected 0 items, got %d", n.Len())
	}
}

func TestCIDRWithIPv4(t *testing.T) {
	n := ipstore.New()
	_, cidr1, _ := net.ParseCIDR("192.168.0.1/24")
	v1 := newValue()
	err := n.AddCIDR(*cidr1, v1)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 1 {
		t.Errorf("expected length to be 1; got: %d", n.Len())
	}

	r, err := n.GetCIDR(*cidr1)
	if err != nil {
		t.Error(err)
	}

	if r[0] != v1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], v1)
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
		t.Errorf("retrieved r2[0] (%#+v) does not equal v2 (%#+v)", r2[0], v2)
	}

	_, cidr3, _ := net.ParseCIDR("192.168.0.1/16")
	r3, err := n.GetCIDR(*cidr3)
	if err != nil {
		t.Error(err)
	}

	if len(r3) != 0 {
		t.Errorf("expected length to be 0; got: %d", len(r3))
	}

	v3 := newValue()
	err = n.AddCIDR(*cidr3, v3)
	if err != nil {
		t.Error(err)
	}

	r3, err = n.GetCIDR(*cidr3)
	if err != nil {
		t.Error(err)
	}

	if len(r3) != 1 {
		t.Errorf("expected length to be 1; got: %d", len(r3))
	}

	if r3[0] != v3 {
		t.Errorf("retrieved r3[0] (%#+v) does not equal v3 (%#+v)", r3[0], v3)
	}

	rr3, err := n.RemoveCIDR(*cidr3)
	if err != nil {
		t.Error(err)
	}

	if rr3 != v3 {
		t.Errorf("removed rr3(%#+v) does not equal v3 (%#+v)", rr3, v3)
	}
}

func TestCombinedIPv4(t *testing.T) {
	n := ipstore.New()

	cv3 := "127.0.1.1/16"
	_, cidr3, _ := net.ParseCIDR(cv3)

	err := n.AddCIDR(*cidr3, cv3)
	if err != nil {
		t.Error(err)
	}

	cv4 := "127.1.1.1/8"
	_, cidr4, _ := net.ParseCIDR(cv4)

	err = n.AddCIDR(*cidr4, cv4)
	if err != nil {
		t.Error(err)
	}

	ip1 := net.ParseIP("127.0.0.1")
	iv1 := "127.0.0.1/32"
	err = n.Add(ip1, iv1)
	if err != nil {
		t.Error(err)
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

	if r[0] != iv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], iv1)
	}

	cv1 := "192.168.0.1/24"
	_, cidr1, _ := net.ParseCIDR(cv1)
	err = n.AddCIDR(*cidr1, cv1)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 4 {
		t.Errorf("expected length to be 4; got: %d", n.Len())
	}

	r, err = n.GetCIDR(*cidr1)
	if err != nil {
		t.Error(err)
	}

	if r[0] != cv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], cv1)
	}

	cv2 := "127.0.0.1/24"
	_, cidr2, _ := net.ParseCIDR(cv2)
	err = n.AddCIDR(*cidr2, cv2)
	if err != nil {
		t.Error(err)
	}

	r, err = n.GetCIDR(*cidr2)
	if err != nil {
		t.Error(err)
	}

	if len(r) != 1 {
		t.Errorf("expected length to be 1; got: %d", len(r))
	}

	if r[0] != cv2 {
		t.Errorf("retrieved r[0] (%#+v) does not equal cv2 (%#+v)", r[0], cv2)
	}

	r, err = n.Get(ip1)
	if err != nil {
		t.Error(err)
	}

	// expecting the most specific match to be the first one (i.e. 127.0.0.1/32, the IP of iv1)
	if r[0] != iv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal iv1 (%#+v)", r[0], iv1)
	}
}

func TestIPV6(t *testing.T) {
	n := ipstore.New()

	cv3 := "2001:db8:1234::/48"
	_, cidr3, _ := net.ParseCIDR(cv3)

	err := n.AddCIDR(*cidr3, cv3)
	if err != nil {
		t.Error(err)
	}

	cv4 := "2001:db8::/32"
	_, cidr4, _ := net.ParseCIDR(cv4)

	err = n.AddCIDR(*cidr4, cv4)
	if err != nil {
		t.Error(err)
	}

	ip1 := net.ParseIP("::1")
	iv1 := "::1/128"
	err = n.Add(ip1, iv1)
	if err != nil {
		t.Error(err)
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

	if r[0] != iv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], iv1)
	}

	cv1 := "2001:db8:a::/64"
	_, cidr1, _ := net.ParseCIDR(cv1)
	err = n.AddCIDR(*cidr1, cv1)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 4 {
		t.Errorf("expected length to be 4; got: %d", n.Len())
	}

	r, err = n.GetCIDR(*cidr1)
	if err != nil {
		t.Error(err)
	}

	if r[0] != cv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], cv1)
	}

	cv2 := "3fff::/20"
	_, cidr2, _ := net.ParseCIDR(cv2)
	err = n.AddCIDR(*cidr2, cv2)
	if err != nil {
		t.Error(err)
	}

	r, err = n.GetCIDR(*cidr2)
	if err != nil {
		t.Error(err)
	}

	if len(r) != 1 {
		t.Errorf("expected length to be 1; got: %d", len(r))
	}

	if r[0] != cv2 {
		t.Errorf("retrieved r[0] (%#+v) does not equal cv2 (%#+v)", r[0], cv2)
	}

	r, err = n.Get(ip1)
	if err != nil {
		t.Error(err)
	}

	// expecting the most specific match to be the first one (i.e. 127.0.0.1/32, the IP of iv1)
	if r[0] != iv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal iv1 (%#+v)", r[0], iv1)
	}
}

func BenchmarkInsertions24Bits(b *testing.B) {
	s := ipstore.New()
	ips, _, _ := hosts("192.168.0.1/24")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String)
		}
		s = ipstore.New()
	}
}

func BenchmarkInsertions16Bits(b *testing.B) {
	s := ipstore.New()
	ips, _, _ := hosts("192.168.0.1/16")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String)
		}
		s = ipstore.New()
	}
}

func BenchmarkRetrievals24Bits(b *testing.B) {
	s := ipstore.New()
	ips, _, _ := hosts("192.168.0.1/24")

	for _, ip := range ips {
		s.Add(ip, ip.String)
	}

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Get(ip)
		}
	}
}

func BenchmarkRetrievals16Bits(b *testing.B) {
	s := ipstore.New()
	ips, _, _ := hosts("192.168.0.1/16")

	for _, ip := range ips {
		s.Add(ip, ip.String)
	}

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Get(ip)
		}
	}
}

func BenchmarkDeletes24Bits(b *testing.B) {
	s := ipstore.New()
	ips, _, _ := hosts("192.168.0.0/24")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String)
		}
		for _, ip := range ips {
			s.Remove(ip)
		}
		s = ipstore.New()
	}
}

func BenchmarkDeletes16Bits(b *testing.B) {
	s := ipstore.New()
	ips, _, _ := hosts("192.168.0.1/16")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String)
		}
		for _, ip := range ips {
			s.Remove(ip)
		}
		s = ipstore.New()
	}
}

func BenchmarkMixed24Bits(b *testing.B) {
	s := ipstore.New()
	ips1, _, _ := hosts("192.168.0.1/24")
	ips2, _, _ := hosts("10.0.0.1/24")

	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for _, ip := range ips1 {
				s.Add(ip, ip.String())
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for _, ip := range ips2 {
				s.Get(ip)
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for _, ip := range ips1 {
				s.Get(ip)
			}
			wg.Done()
		}()

		wg.Wait()
		s = ipstore.New()
	}
}

func BenchmarkMixed16Bits(b *testing.B) {
	s := ipstore.New()
	ips1, _, _ := hosts("192.168.0.1/16")
	ips2, _, _ := hosts("10.0.0.1/16")

	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for _, ip := range ips1 {
				s.Add(ip, ip.String())
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for _, ip := range ips2 {
				s.Get(ip)
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for _, ip := range ips1 {
				s.Get(ip)
			}
			wg.Done()
		}()

		wg.Wait()
		s = ipstore.New()
	}
}
