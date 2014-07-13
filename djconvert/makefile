#!/usr/bin/make

.PHONY: test

djconvert: main.go
	go build

test: djconvert
	set -e; for t in tests/*; do \
		./test $$t; \
	done
	# All tests pass
