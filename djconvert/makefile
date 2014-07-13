#!/usr/bin/make

.PHONY: test

djconvert: main.go
	go build

test: djconvert
	for t in tests/*; do \
		./test $$t; \
	done
	# All tests pass
