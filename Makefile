
.PHONY: build clean test install
PREFIX?=/usr/local
BINARY_PATH=${PWD}/bin
export

build: ${BINARY}/runn

${BINARY}/runn:
	@$(MAKE) -C runn build
	
clean:
	@$(MAKE) -C runn clean
	
install:
	@$(MAKE) -C runn install

uninstall:
	@$(MAKE) -C runn uninstall