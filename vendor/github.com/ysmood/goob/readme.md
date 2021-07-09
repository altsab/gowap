# Overview

A lightweight observable lib. Go channel doesn't support unlimited buffer size,
it's a pain to decide what size to use, this lib will handle it dynamically.

- unlimited buffer size
- one publisher to multiple subscribers
- thread-safe
- subscribers never block each other
- stable event order

## Examples

See [examples_test.go](examples_test.go).
