default: run

.PHONY: default run test import lib



lib: test
	GOPATH=`pwd` go install littlereader

import: lib
	GOPATH=`pwd` go run src/progs/importer.go

run: lib
	GOPATH=`pwd` go run src/progs/reader.go

test:
	GOPATH=`pwd` go test littlereader
