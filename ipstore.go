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

// Store is a simple Key/Value store using IPs and CIDRs as keys.
type Store[T any] struct {
	mu    sync.RWMutex
	table *bart.Table[T]
}

// New returns a new instance of [Store].
func New[T any]() *Store[T] {
	return &Store[T]{
		mu:    sync.RWMutex{},
		table: new(bart.Table[T]),
	}
}

// Add adds a new entry to the store mapped by [netip.Addr].
func (s *Store[T]) Add(key netip.Addr, value T) error {
	prf, err := key.Prefix(key.BitLen())
	if err != nil {
		return err
	}

	return s.AddCIDR(prf, value)
}

// AddCIDR adds a new entry to the store mapped by [netip.Prefix].
func (s *Store[T]) AddCIDR(key netip.Prefix, value T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.table.Insert(key, value)

	return nil
}

// AddIPOrCIDR adds a new entry to the [Store] mapped by an IP or CIDR.
func (s *Store[T]) AddIPOrCIDR(ipOrCIDR string, value T) error {
	prf, err := parsePrefix(ipOrCIDR)
	if err != nil {
		return err
	}

	return s.AddCIDR(prf, value)
}

// Remove removes the entry associated with [netip.Addr] from [Store].
func (s *Store[T]) Remove(key netip.Addr) (T, error) {
	prf, err := key.Prefix(key.BitLen())
	if err != nil {
		return zero[T](), err
	}

	return s.RemoveCIDR(prf)
}

// RemoveCIDR removes the entry associated with [netip.Prefix] from [Store].
func (s *Store[T]) RemoveCIDR(key netip.Prefix) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.table.GetAndDelete(key)
	if !ok {
		return zero[T](), nil
	}

	return value, nil
}

// RemoveIPOrCIDR removes the entry associated with an IP or CIDR from [Store].
func (s *Store[T]) RemoveIPOrCIDR(ipOrCIDR string) (T, error) {
	prf, err := parsePrefix(ipOrCIDR)
	if err != nil {
		return zero[T](), err
	}

	return s.RemoveCIDR(prf)
}

// Contains returns whether an entry is available for the [netip.Addr].
func (s *Store[T]) Contains(ip netip.Addr) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.table.Lookup(ip)

	return ok, nil
}

// Get returns entries from the [Store] based on the [netip.Addr]
// key. Because multiple CIDRs may contain the key, a slice of
// entries is returned instead of a single entry.
func (s *Store[T]) Get(key netip.Addr) ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prf, err := key.Prefix(key.BitLen())
	if err != nil {
		return nil, err
	}

	var result = make([]T, 0, 5)
	supernets := s.table.Supernets(prf)
	supernets(func(p netip.Prefix, t T) bool {
		result = append(result, t)
		return true
	})

	return result, nil
}

// GetOne returns a single entry from the [Store] based on the
// [netip.Addr] key.
func (s *Store[T]) GetOne(key netip.Addr) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.table.Lookup(key)
}

// GetCIDR returns entries from the [Store] by [netip.Prefix].
func (s *Store[T]) GetCIDR(key netip.Prefix) ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result = make([]T, 0, 5)
	supernets := s.table.Supernets(key)
	supernets(func(p netip.Prefix, t T) bool {
		result = append(result, t)
		return true
	})

	return result, nil
}

// GetOneCIDR returns a single entry from the [Store] by [netip.Prefix].
func (s *Store[T]) GetOneCIDR(key netip.Prefix) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.table.LookupPrefix(key)
}

// GetIPOrCIDR returns entries from the [Store] by IP or CIDR.
func (s *Store[T]) GetIPOrCIDR(ipOrCIDR string) ([]T, error) {
	prf, err := parsePrefix(ipOrCIDR)
	if err != nil {
		return nil, err
	}

	return s.GetCIDR(prf)
}

// GetIPOrCIDR returns entries from the [Store] by IP or CIDR.
func (s *Store[T]) GetOneIPOrCIDR(ipOrCIDR string) (T, bool) {
	prf, err := parsePrefix(ipOrCIDR)
	if err != nil {
		return zero[T](), false
	}

	return s.GetOneCIDR(prf)
}

// Len returns the number of entries in the [Store].
func (s *Store[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.table.Size()
}

func zero[T any]() T {
	return *new(T)
}

func parsePrefix(s string) (netip.Prefix, error) {
	ip, err := netip.ParseAddr(s)
	if err != nil || !ip.IsValid() {
		return netip.ParsePrefix(s)
	}

	return ip.Prefix(ip.BitLen())
}
