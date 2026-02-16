# ipstore

A fast and simple key-value store using `netip.Addr` and `netip.Prefix` as keys.

## Description

The `ipstore` package provides a fast and simple in-memory key-value storage for network addresses.
You can store Go `any` types indexed by `netip.Addr` and `netip.Prefix` keys.
The storage is generically typed at time of instantiation, meaning that type safety is provided through Go's generics support.

The heavy lifting is done by the awesome [bart](https://github.com/gaissmai/bart) package, which provides a (compact) multibit-trie to efficiently index IP addresses and CIDR ranges.

## Usage

```bash
go get "github.com/hslatman/ipstore"
```

```go
package main

import (
    "fmt"
    "net/netip"
    
    "github.com/hslatman/ipstore"
)

func main() {
    // create new instance taking strings as values
    store := ipstore.New[string]()
    ip := netip.ParseAddr("127.0.0.1")
    
    // add `value` to the store indexed by IP; can result in error
    _ := store.Add(ip, "value")

    // returns the number of entries in the store
    _ := store.Len()

    // returns true
    t, _ := store.Contains(ip)

    // returns all entries (CIDR ranges) matching the IP
    r, err := store.Get(ip)
    for _, e := range r {
        fmt.Println(e)
    }

    // deletes the entry from the store
    _ := store.Remove(ip)
}
```

## Benchmarks

```bash
$ go test -run=XXX -benchmem -bench=. ./...
goos: darwin
goarch: arm64
pkg: github.com/hslatman/ipstore
cpu: Apple M1 Max
BenchmarkInsertions24Bits-10    	   41072	     25289 ns/op	   18336 B/op	     532 allocs/op
BenchmarkInsertions16Bits-10    	     176	   7312219 ns/op	 4587788 B/op	  134277 allocs/op
BenchmarkRetrievals24Bits-10    	   50401	     26365 ns/op	   20480 B/op	     256 allocs/op
BenchmarkRetrievals16Bits-10    	     163	   7240971 ns/op	 5326579 B/op	   66759 allocs/op
BenchmarkDeletes24Bits-10       	   20961	     58683 ns/op	   18481 B/op	     535 allocs/op
BenchmarkDeletes16Bits-10       	      76	  14277057 ns/op	 4668238 B/op	  135025 allocs/op
BenchmarkMixed24Bits-10         	    7713	    144770 ns/op	   59473 B/op	    1048 allocs/op
BenchmarkMixed16Bits-10         	      27	  44661315 ns/op	15696335 B/op	  269838 allocs/op
PASS
ok  	github.com/hslatman/ipstore	14.715s
```

## [cidranger](https://github.com/yl2chen/cidranger) vs. [bart](https://github.com/gaissmai/bart)

This repository used to use a fork of [github.com/yl2chen/cidranger](https://github.com/yl2chen/cidranger).
We've switched to [github.com/gaissmai/bart](https://github.com/gaissmai/bart) for better overall performance, and it contains all the functionality needed exposed in `ipstore`.

## TODO

* Improve README
* Add (more) tests
* Add additional (utility) functions?
* Add function for retrieving the first match? I.e. the first `any` in the slice that matches.