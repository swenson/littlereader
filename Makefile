default: run

.PHONY: default run test import lib requirements

requirements:
	./requirements.sh

reader: lib
	GOPATH=`pwd` go build src/progs/reader.go

lib: test
	GOPATH=`pwd` go install littlereader

import: lib
	GOPATH=`pwd` go run src/progs/importer.go

run: lib
	GOPATH=`pwd` go run src/progs/reader.go

test: requirements
	GOPATH=`pwd` go test littlereader
