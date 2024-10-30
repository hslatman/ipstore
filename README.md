# ipstore

A fast and simple key-value store using `netip.Addr` and `netip.Prefix` as keys.

## Description

The `ipstore` package provides a fast and simple in-memory key-value storage for network addresses.
You can store Go `any` types indexed by `netip.Addr` and `netip.Prefix` keys.
The storage is generically typed at time of instantiation, meaning that type safety is provided through Go's generics support.

The heavy lifting is done by the awesome [https://github.com/gaissmai/bart](github.com/gaissmai/bart) package, which provides a (compact) multibit-trie to efficiently index IP addresses and CIDR ranges.

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
    store := ipstore.New()
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
BenchmarkInsertions24Bits-10    	   45368	     26192 ns/op	   14360 B/op	     290 allocs/op
BenchmarkInsertions16Bits-10    	     170	   7054662 ns/op	 3564928 B/op	   69534 allocs/op
BenchmarkRetrievals24Bits-10    	   35632	     32845 ns/op	   20481 B/op	     256 allocs/op
BenchmarkRetrievals16Bits-10    	     133	   8632705 ns/op	 5337738 B/op	   66549 allocs/op
BenchmarkDeletes24Bits-10       	   21706	     69278 ns/op	   14393 B/op	     294 allocs/op
BenchmarkDeletes16Bits-10       	      76	  14841112 ns/op	 3633203 B/op	   70270 allocs/op
BenchmarkMixed24Bits-10         	    7021	    148739 ns/op	   55494 B/op	     806 allocs/op
BenchmarkMixed16Bits-10         	      21	  55264974 ns/op	14864282 B/op	  206469 allocs/op
PASS
ok  	github.com/hslatman/ipstore	14.057s
```

## [https://github.com/yl2chen/cidranger](cidranger) vs. [https://github.com/gaissmai/bart](bart)

This repository used to use a fork of the [yl2chen/cidranger](https://github.com/yl2chen/cidranger).
We've switched to [github.com/gaissmai/bart](github.com/gaissmai/bart) for better overall performance, and it contains all the functionality needed exposed in `ipstore`.

## TODO

* Improve README
* Add (more) tests
* Add additional (utility) functions?
* Add function for retrieving the first match? I.e. the first `any` in the slice that matches.