.PHONY: all
all: goboardfs

SRCS=$(filter-out %_test.go, $(wildcard *.go */*.go))
TESTSRCS=$(wildcard *_test.go */*_test.go) 

goboardfs: $(SRCS)
	go build

.PHONY: test
test: test.out
	cat test.out

test.out: goboardfs $(TESTSRCS)
	go test -v ./board >$@ 2>&1 || cat $@

.PHONY: sloc
sloc:
	wc -l `ls *go */*go | grep -v test`

.PHONY: sloc-board
sloc-board:
	wc -l `ls board/*go | grep -v test`


.PHONY: sloc-all
sloc-all:
	wc -l *go */*go

.PHONY: vet
vet: vet.out
	cat vet.out

vet.out: goboardfs
	go tool vet -v ./board ./boardfs >$@ 2>&1 || cat $@

