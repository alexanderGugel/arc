[![Build Status](https://travis-ci.org/alexanderGugel/arc.svg?branch=master)](https://travis-ci.org/alexanderGugel/arc)

arc
===

An [Adaptive Replacement Cache (ARC)](http://web.archive.org/web/20150405221102/https://www.usenix.org/legacy/event/fast03/tech/full_papers/megiddo/megiddo.pdf) written in [Go](http://golang.org/).

[GoDoc](https://godoc.org/github.com/alexanderGugel/arc)

This project implements "ARC", a self-tuning, low overhead replacement cache. The goal of this project is to expose an interface compareable to common LRU cache management systems. ARC uses a learning rule to adaptively and continually revise its assumptions about the workload in order to adjust the internal LRU and LFU cache sizes.

This implementation is based on Nimrod Megiddo and Dharmendra S. Modha's ["ARC: A SELF-TUNING, LOW OVERHEAD REPLACEMENT CACHE"](http://web.archive.org/web/20150405221102/https://www.usenix.org/legacy/event/fast03/tech/full_papers/megiddo/megiddo.pdf), while definitely useable and thread safe, this is still an experiment and shouldn't be considered production-ready.

```
<------- cache size c ------>
+-----------------+----------
| LFU             | LRU     |
+-----------------+----------
                  ^
                  |
                  p (dynamically adjusted by learning rule)

B1 [...]
B2 [...]
```

The cache is implemented using two internal caching systems L1 and L2. The cache size c defines the maximum number of entries stored (excluding ghost entries). Ghost entries are being stored in two "ghost registries" B1 and B1. Ghost entries no longer have a value associated with them.

Ghost entries are being used in order to keep track of expelled pages. They no longer have a value associated with them, but can be promoted into the internal LRU cache.

Frequently requested pages are being promoted into the LFU.
