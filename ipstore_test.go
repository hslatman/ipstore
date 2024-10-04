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
	"net/netip"
	"sync"
	"testing"

	"github.com/hslatman/ipstore"
)

type value struct {
	v int
}

func newValue() *value {
	return &value{
		v: rand.Int(),
	}
}

func hosts(tb testing.TB, cidr string) ([]netip.Addr, int) {
	tb.Helper()

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		tb.Fatal(err)
	}

	var ips []netip.Addr
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, netip.MustParseAddr(ip.String()))
	}

	return ips, len(ips)
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
	n := ipstore.New[*value]()
	if n == nil {
		t.Fail()
	}
}

func TestIPWithIPv4(t *testing.T) {
	n := ipstore.New[*value]()
	ip1 := netip.MustParseAddr("127.0.0.1")
	v1 := newValue()
	err := n.Add(ip1, v1)
	if err != nil {
		t.Error(err)
	}

	ip2 := netip.MustParseAddr("127.0.0.2")
	v2 := newValue()
	err = n.Add(ip2, v2)
	if err != nil {
		t.Error(err)
	}

	ip3 := netip.MustParseAddr("192.168.1.0")
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
		t.Error("expected ip2 to be in store")
	}

	if r[0] != v2 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v2 (%#+v)", r[0], v2)
	}

	r, err = n.Get(ip3)
	if err != nil {
		t.Error(err)
	}

	if r[0] == nil {
		t.Error("expected ip3 to be in store")
	}

	if r[0] != v3 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v3 (%#+v)", r[0], v3)
	}

	v, ok := n.GetOne(ip1)
	if !ok {
		t.Error("expected ip1 to be in store")
	}
	if v != v1 {
		t.Errorf("retrieved v (%#+v) does not equal v1 (%#+v)", v, v1)
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
	n := ipstore.New[*value]()

	cidr1 := netip.MustParsePrefix("192.168.0.1/24")
	v1 := newValue()
	err := n.AddCIDR(cidr1, v1)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 1 {
		t.Errorf("expected length to be 1; got: %d", n.Len())
	}

	r, err := n.GetCIDR(cidr1)
	if err != nil {
		t.Error(err)
	}

	if r[0] != v1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], v1)
	}

	cidr2 := netip.MustParsePrefix("192.168.1.1/24")
	r2, err := n.GetCIDR(cidr2)
	if err != nil {
		t.Error(err)
	}

	if len(r2) != 0 {
		t.Error("retrieved CIDR that should not be retrieved")
	}

	v2 := newValue()
	err = n.AddCIDR(cidr2, v2)
	if err != nil {
		t.Error(err)
	}

	r2, err = n.GetCIDR(cidr2)
	if err != nil {
		t.Error(err)
	}

	if len(r2) == 0 {
		t.Error("did  not retrieve CIDR that should not be retrieved")
	}

	if r2[0] != v2 {
		t.Errorf("retrieved r2[0] (%#+v) does not equal v2 (%#+v)", r2[0], v2)
	}

	cidr3 := netip.MustParsePrefix("192.168.0.1/16")
	r3, err := n.GetCIDR(cidr3)
	if err != nil {
		t.Error(err)
	}

	if len(r3) != 0 {
		t.Errorf("expected length to be 0; got: %d", len(r3))
	}

	v3 := newValue()
	err = n.AddCIDR(cidr3, v3)
	if err != nil {
		t.Error(err)
	}

	r3, err = n.GetCIDR(cidr3)
	if err != nil {
		t.Error(err)
	}

	if len(r3) != 1 {
		t.Errorf("expected length to be 1; got: %d", len(r3))
	}

	if r3[0] != v3 {
		t.Errorf("retrieved r3[0] (%#+v) does not equal v3 (%#+v)", r3[0], v3)
	}

	rr3, err := n.RemoveCIDR(cidr3)
	if err != nil {
		t.Error(err)
	}

	if rr3 != v3 {
		t.Errorf("removed rr3(%#+v) does not equal v3 (%#+v)", rr3, v3)
	}
}

func TestCombinedIPv4(t *testing.T) {
	n := ipstore.New[string]()

	cv3 := "127.0.1.1/16"
	cidr3 := netip.MustParsePrefix(cv3)

	err := n.AddCIDR(cidr3, cv3)
	if err != nil {
		t.Error(err)
	}

	cv4 := "127.1.1.1/8"
	cidr4 := netip.MustParsePrefix(cv4)
	err = n.AddCIDR(cidr4, cv4)
	if err != nil {
		t.Error(err)
	}

	ip1 := netip.MustParseAddr("127.0.0.1")
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

	if r[0] == "" {
		t.Error("expected ip1 to be in store")
	}

	if r[0] != iv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], iv1)
	}

	cv1 := "192.168.0.1/24"
	cidr1 := netip.MustParsePrefix(cv1)
	err = n.AddCIDR(cidr1, cv1)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 4 {
		t.Errorf("expected length to be 4; got: %d", n.Len())
	}

	r, err = n.GetCIDR(cidr1)
	if err != nil {
		t.Error(err)
	}

	if r[0] != cv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], cv1)
	}

	cv2 := "127.0.0.1/24"
	cidr2 := netip.MustParsePrefix(cv2)
	err = n.AddCIDR(cidr2, cv2)
	if err != nil {
		t.Error(err)
	}

	r, err = n.GetCIDR(cidr2)
	if err != nil {
		t.Error(err)
	}

	if len(r) != 3 {
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
	n := ipstore.New[string]()

	cv3 := "2001:db8:1234::/48"
	cidr3 := netip.MustParsePrefix(cv3)
	err := n.AddCIDR(cidr3, cv3)
	if err != nil {
		t.Error(err)
	}

	cv4 := "2001:db8::/32"
	cidr4 := netip.MustParsePrefix(cv4)
	err = n.AddCIDR(cidr4, cv4)
	if err != nil {
		t.Error(err)
	}

	ip1 := netip.MustParseAddr("::1")
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

	if r[0] == "" {
		t.Error("expected ip1 to be in store")
	}

	if r[0] != iv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], iv1)
	}

	cv1 := "2001:db8:a::/64"
	cidr1 := netip.MustParsePrefix(cv1)
	err = n.AddCIDR(cidr1, cv1)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 4 {
		t.Errorf("expected length to be 4; got: %d", n.Len())
	}

	r, err = n.GetCIDR(cidr1)
	if err != nil {
		t.Error(err)
	}

	if r[0] != cv1 {
		t.Errorf("retrieved r[0] (%#+v) does not equal v1 (%#+v)", r[0], cv1)
	}

	cv2 := "3fff::/20"
	cidr2 := netip.MustParsePrefix(cv2)
	err = n.AddCIDR(cidr2, cv2)
	if err != nil {
		t.Error(err)
	}

	r, err = n.GetCIDR(cidr2)
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

func TestDuplicateInsertion(t *testing.T) {
	s := ipstore.New[string]()
	ip := netip.MustParseAddr("127.0.0.1")

	err := s.Add(ip, "value1")
	if err != nil {
		t.Error(err)
	}

	err = s.Add(ip, "value2")
	if err != nil {
		t.Error(err)
	}

	r, err := s.Get(ip)
	if err != nil {
		t.Error(err)
	}

	if len(r) != 1 {
		t.Errorf("expected 1 result; got %d", len(r))
	}

	if r[0] != "value2" {
		t.Errorf(`expected result to be "value2"; got %q`, r[0])
	}
}

func TestIPOrCIDR(t *testing.T) {
	s := ipstore.New[string]()
	ip1 := "127.0.0.1"
	ip2 := "127.0.0.1/32"
	ip3 := "127.0.0.2"
	range1 := "127.0.0.1/24"

	err := s.AddIPOrCIDR(ip1, ip1)
	if err != nil {
		t.Error(err)
	}

	r, err := s.GetIPOrCIDR(ip1)
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("expected 1 result; got %d", len(r))
	}
	if r[0] != ip1 {
		t.Errorf("expected %q; got %q", ip1, r[0])
	}

	err = s.AddIPOrCIDR(ip2, ip2)
	if err != nil {
		t.Error(err)
	}

	r, err = s.GetIPOrCIDR(ip1) // result for ip1 should be same as for ip2
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("expected 1 result; got %d", len(r))
	}
	if r[0] != ip2 {
		t.Errorf("expected %q; got %q", ip2, r[0]) // result overwritten
	}

	err = s.AddIPOrCIDR(ip3, ip3)
	if err != nil {
		t.Error(err)
	}

	err = s.AddIPOrCIDR(range1, range1)
	if err != nil {
		t.Error(err)
	}

	r, err = s.GetIPOrCIDR("127.0.0.100") // within 127.0.0.1/24 range
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("expected 1 result; got %d", len(r))
	}
	if r[0] != range1 {
		t.Errorf("expected %q; got %q", range1, r[0])
	}

	r, err = s.GetIPOrCIDR(ip1) // result for ip1 should be same as for ip2
	if err != nil {
		t.Error(err)
	}
	if len(r) != 2 {
		t.Errorf("expected 1 result; got %d", len(r))
	}
	if r[0] != ip2 {
		t.Errorf("expected %q; got %q", ip2, r[0])
	}

	// remove the specific IP, look it up again; should return the
	// range in which it was contained before.
	v, err := s.RemoveIPOrCIDR(ip1)
	if err != nil {
		t.Error(err)
	}
	if v != ip2 {
		t.Errorf("expected %q; got %q", ip2, r[0])
	}

	r, err = s.GetIPOrCIDR(ip1)
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("expected 1 result; got %d", len(r))
	}
	if r[0] != range1 {
		t.Errorf("expected %q; got %q", range1, r[0])
	}
}

func TestGetMultipleResults(t *testing.T) {
	s := ipstore.New[string]()
	ip1 := netip.MustParseAddr("127.0.0.1")
	range1 := netip.MustParsePrefix("127.0.0.1/24")
	range2 := netip.MustParsePrefix("127.0.0.1/23")
	err := s.Add(ip1, ip1.String())
	if err != nil {
		t.Error(err)
	}

	err = s.AddCIDR(range1, range1.String())
	if err != nil {
		t.Error(err)
	}

	err = s.AddCIDR(range2, range2.String())
	if err != nil {
		t.Error(err)
	}

	if s.Len() != 3 {
		t.Errorf("expected 2 entries; got %d entries", s.Len())
	}

	r, err := s.Get(ip1)
	if err != nil {
		t.Error(err)
	}
	if len(r) != 3 {
		t.Errorf("expected 3 results; got %d", len(r))
	}
	if r[0] != ip1.String() {
		t.Errorf("expected %q; got %q", ip1.String(), r[0])
	}

	r, err = s.GetCIDR(range1)
	if err != nil {
		t.Error(err)
	}
	if len(r) != 2 {
		t.Errorf("expected 2 results; got %d", len(r))
	}
	if r[0] != range1.String() {
		t.Errorf("expected %q; got %q", range1.String(), r[0])
	}
}

func TestLen(t *testing.T) {
	n := ipstore.New[*value]()
	if n.Len() != 0 {
		t.Errorf("expected store to be empty; got %d entries", n.Len())
	}

	ip1 := netip.MustParseAddr("127.0.0.1")
	err := n.Add(ip1, newValue())
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 1 {
		t.Errorf("expected 1 entry; got %d entries", n.Len())
	}

	ip2 := netip.MustParseAddr("127.0.0.2")
	err = n.Add(ip2, newValue())
	if err != nil {
		t.Error(err)
	}
	if n.Len() != 2 {
		t.Errorf("expected 2 entries; got %d entries", n.Len())
	}

	// empty the store again
	_, err = n.Remove(ip1)
	if err != nil {
		t.Error(err)
	}
	_, err = n.Remove(ip2)
	if err != nil {
		t.Error(err)
	}

	if n.Len() != 0 {
		t.Errorf("expected store to be empty again; got %d entries", n.Len())
	}
}

func BenchmarkInsertions24Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips, _ := hosts(b, "192.168.0.1/24")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String())
		}
		s = ipstore.New[string]()
	}
}

func BenchmarkInsertions16Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips, _ := hosts(b, "192.168.0.1/16")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String())
		}
		s = ipstore.New[string]()
	}
}

func BenchmarkRetrievals24Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips, _ := hosts(b, "192.168.0.1/24")

	for _, ip := range ips {
		s.Add(ip, ip.String())
	}

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Get(ip)
		}
	}
}

func BenchmarkRetrievals16Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips, _ := hosts(b, "192.168.0.1/16")

	for _, ip := range ips {
		s.Add(ip, ip.String())
	}

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Get(ip)
		}
	}
}

func BenchmarkDeletes24Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips, _ := hosts(b, "192.168.0.0/24")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String())
		}
		for _, ip := range ips {
			s.Remove(ip)
		}
		s = ipstore.New[string]()
	}
}

func BenchmarkDeletes16Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips, _ := hosts(b, "192.168.0.1/16")

	for n := 0; n < b.N; n++ {
		for _, ip := range ips {
			s.Add(ip, ip.String())
		}
		for _, ip := range ips {
			s.Remove(ip)
		}
		s = ipstore.New[string]()
	}
}

func BenchmarkMixed24Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips1, _ := hosts(b, "192.168.0.1/24")
	ips2, _ := hosts(b, "10.0.0.1/24")

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
		s = ipstore.New[string]()
	}
}

func BenchmarkMixed16Bits(b *testing.B) {
	s := ipstore.New[string]()
	ips1, _ := hosts(b, "192.168.0.1/16")
	ips2, _ := hosts(b, "10.0.0.1/16")

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
		s = ipstore.New[string]()
	}
}
