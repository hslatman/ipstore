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
	"net/netip"
	"sync"

	"github.com/gaissmai/bart"
)

// Store is a (simple) Key/Value store using IPs and CIDRs as keys.
type Store[V any] struct {
	sync.RWMutex
	table *bart.Table[entry[V]]
}

type entry[T any] struct {
	value T
}

// New returns a new instance of Store.
func New[V any]() *Store[V] {
	return &Store[V]{
		RWMutex: sync.RWMutex{},
		table:   new(bart.Table[entry[V]]),
	}
}

// Add adds a new entry to the store mapped by [netip.Addr].
func (s *Store[V]) Add(key netip.Addr, value V) error {
	prf, err := key.Prefix(key.BitLen())
	if err != nil {
		return err
	}

	return s.AddCIDR(prf, value)
}

// AddCIDR adds a new entry to the store mapped by [netip.Prefix].
func (s *Store[V]) AddCIDR(key netip.Prefix, value V) error {
	s.Lock()
	defer s.Unlock()

	entry := entry[V]{value: value}

	s.table.Insert(key, entry)

	return nil
}

// AddIPOrCIDR adds a new entry to the store mapped by an IP or CIDR.
func (s *Store[V]) AddIPOrCIDR(ipOrCIDR string, value V) error {
	// TODO: implementation
	return nil
}

// Remove removes the entry associated with [netip.Addr] from [Store].
func (s *Store[V]) Remove(key netip.Addr) (V, error) {
	prf, err := key.Prefix(key.BitLen())
	if err != nil {
		return zero[V](), err
	}

	return s.RemoveCIDR(prf)
}

// RemoveCIDR removes the entry associated with [netip.Prefix] from [Store].
func (s *Store[V]) RemoveCIDR(key netip.Prefix) (V, error) {
	s.Lock()
	defer s.Unlock()

	re, ok := s.table.GetAndDelete(key)
	if !ok {
		return zero[V](), nil
	}

	return re.value, nil
}

// RemoveIPOrCIDR removes the entry associated with an IP or CIDR.
func (s *Store[V]) RemoveIPOrCIDR(ipOrCIDR string, value any) (any, error) {
	// TODO: implementation
	return nil, nil
}

// Contains returns whether an entry is available for the [netip.Addr].
func (s *Store[V]) Contains(ip netip.Addr) (bool, error) {
	s.RLock()
	defer s.RUnlock()

	_, ok := s.table.Lookup(ip)

	return ok, nil
}

// Get returns entries from the [Store] based on the [netip.Addr]
// key. Because multiple CIDRs may contain the key, a slice of
// entries is returned instead of a single entry.
func (s *Store[V]) Get(key netip.Addr) ([]V, error) {
	s.RLock()
	defer s.RUnlock()

	e, ok := s.table.Lookup(key)
	if !ok {
		return nil, nil
	}

	result := []V{e.value}

	// TODO: return multiple results (again); supernets?

	// r, err := s.trie.ContainingNetworks(net.IP(key.AsSlice()))
	// if err != nil {
	// 	return nil, err
	// }

	// // return all networks that this IP is part by reverse looping through the result
	// // haven't fully deduced it yet, but it seems that the order of the entries from ContainingNetworks
	// // are from biggest CIDR to smallest CIDR. I think the most logical thing to do is to return the
	// // most specific CIDR that the net.IP is part of first instead of last, so that's why the
	// // returned slice of any is reversed.
	// // TODO: verify that this is correct?
	// var result []any
	// for i := len(r) - 1; i >= 0; i-- {
	// 	e, _ := r[i].(entry) // type is guarded by Add/AddCIDR
	// 	result = append(result, e.value)
	// }

	return result, nil
}

// GetCIDR returns entries from the [Store] by [netip.Prefix].
func (s *Store[V]) GetCIDR(key netip.Prefix) ([]V, error) {
	s.RLock()
	defer s.RUnlock()

	e, ok := s.table.LookupPrefix(key)
	if !ok {
		return nil, nil
	}

	result := []V{e.value}

	// TODO: return multiple results (again); subnets?

	// subnets := s.table.Subnets(key)
	// subnets(func(p netip.Prefix, e entry) bool {

	// })

	// TODO: decide if we only want to return a single any, because a specific CIDR should only exist once now

	// // first perform exact match of the network
	// t, err := s.trie.ContainsNetwork(key)
	// if err != nil {
	// 	return nil, err
	// }

	// TODO: decide if we want to keep the check above; we could also call CoveredNetworks and just loop.
	// The additional call to ContainsNetwork was the reason I forked the original library, so we might
	// be able to return to using the original instead of the fork at github.com/hslatman/cidranger.
	// There are also changes in other forks that may be of interest, though ...

	// // return with empty result if there's no exact match
	// if !t {
	// 	return nil, nil
	// }

	// // get all covered networks, including the exact match (if it exists) and smaller CIDR ranges
	// r, err := s.trie.CoveredNetworks(key)
	// if err != nil {
	// 	return nil, err
	// }

	// // loop through the results and do a full equality check on the IP and IPMask
	// var result []any
	// for _, re := range r {
	// 	e, _ := re.(entry)                            // type is guarded by Add/AddCIDR
	// 	keyMaskOnes, keyMaskZeroes := key.Mask.Size() // TODO: improve the equality check? Is what we do here correct?
	// 	entryMaskOnes, entryMaskZeroes := e.net.Mask.Size()
	// 	if key.IP.Equal(e.net.IP) && keyMaskOnes == entryMaskOnes && keyMaskZeroes == entryMaskZeroes {
	// 		result = append(result, e.value)
	// 	}
	// }

	return result, nil
}

// Len returns the number of entries in the [Store].
func (s *Store[V]) Len() int {
	s.RLock()
	defer s.RUnlock()

	return s.table.Size()
}

func zero[V any]() V {
	return *new(V)
}
