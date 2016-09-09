
VERSION="0.0.1"
BUILD_TIME=`date +%FT%T%z`
GO=`which go`

PREFIX?=/usr/local
BINARY_PATH?=${PWD}/bin

BINARY=${BINARY_PATH}/runn

UNAME := $(shell uname)

LDFLAGS=-ldflags "-s -w -X github.com/kildevaeld/runn/runn/cmd.VERSION=${VERSION} -X github.com/kildevaeld/runn/runn/cmd.BuildDate=${BUILD_TIME}"

.PHONY: build clean run

build: ${BINARY}

${BINARY}:
	@${GO} build -x ${LDFLAGS} -o $@ main.go



pack: build
	@upx --lzma -o runnp ${BINARY}
	

clean:
	rm -f ${BINARY}

install: ${BINARY}
	install -m 0755 ${BINARY} ${PREFIX}/bin/runn

uninstall:
	rm ${PREFIX}/bin/runn