# ipstore

A fast and simple key-value store using net.IP and net.IPNet as keys

## Description

The `ipstore` package provides a fast and simple in-memory key-value storage for network addresses.
You can store Go `any` types indexed by `net.IP` and `net.IPNet` keys.
You'll have to provide your own type checking logic on top of this library to ensure type safety.

The heavy lifting is done by `github.com/yl2chen/cidranger`, which provides a (compact) prefix tree (or radix tree/trie) to efficiently index IP addresses and CIDR ranges.

## Usage

```bash
go get "github.com/hslatman/ipstore"
```

```go
package main

import (
	"github.com/hslatman/ipstore"
)

func main() {
    
    store := ipstore.New()
    ip := net.ParseIP("127.0.0.1")
    
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
$ go test -run=XXX -bench=. ./...

pkg: github.com/hslatman/ipstore/pkg/ipstore
BenchmarkInsertions24Bits-8      4935	    233720 ns/op	   61463 B/op	    3840 allocs/op
BenchmarkInsertions16Bits-8        18	  64520756 ns/op	16252053 B/op	  983078 allocs/op
BenchmarkRetrievals24Bits-8     27381	     42750 ns/op	    9218 B/op	     768 allocs/op
BenchmarkRetrievals16Bits-8        94	  12280622 ns/op	 2625165 B/op	  207067 allocs/op
BenchmarkMixed24Bits-8         155709	     35556 ns/op	   10912 B/op	    1005 allocs/op
BenchmarkMixed16Bits-8          51330	     58402 ns/op	   11549 B/op	    1293 allocs/op
PASS
ok  	github.com/hslatman/ipstore/pkg/ipstore	29.489s
```

## cidranger

Currently this repository uses a fork of the [yl2chen/cidranger](https://github.com/yl2chen/cidranger), because we use a function to exactly match a CIDR network that is not yet available in the original.
The fork is [hslatman/cidranger](https://github.com/hslatman/cidranger).
We might switch back to the original, a different fork or even a different library (like [critbitgo](https://github.com/k-sone/critbitgo)) for handling the radix tree in the future, but for now we're OK with the current fork.

## TODO

* Improve README
* Add (more) tests
* Add benchmarking
* Add additional (utility) functions?
* Add a bit of type safety functionality? (we could block different any types from insertions, but makes it slower).
* Add function for retrieving the first match? I.e. the first `any` in the slice that matches.