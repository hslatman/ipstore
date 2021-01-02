# ipstore

A fast and simple key-value store using net.IP and net.IPNet as keys

## Description

The `ipstore` package provides a fast and simple in-memory key-value storage for network addresses.
You can store Go `interface{}` types indexed by `net.IP` and `net.IPNet` keys.
You'll have to provide your own type checking logic on top of this library to ensure type safety.

The heavy lifting is done by `github.com/yl2chen/cidranger`, which provides a (compact) prefix tree (or radix tree/trie) to efficiently index IP addresses and CIDR ranges.

## Usage

```bash
go get "github.com/hslatman/ipstore/pkg/ipstore"
```

```go
package main

import (
	"github.com/hslatman/ipstore/pkg/ipstore"
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

## TODO

* Improve README
* Add (more) tests
* Add benchmarking
* Add additional (utility) functions?
* Decide whether to use different fork (i.e. censys/cidranger?) or maintain our own? Or a different library (i.e. https://github.com/k-sone/critbitgo)?
* Add a bit of type safety functionality? (we could block different interface{} types from insertions, but makes it slower).
* Add function for retrieving the first match? I.e. the first interface{} in the slice that matches.