# phzd
# Makefile generated by 'gomaker' 

NAME ?= phzd
VERSION ?= 
PREFIX ?= /usr/local/bin
VER ?= X
COMMIT=$(shell git rev-parse --verify --short HEAD)
COMMIT ?= ${VER}
RELEASE ?= ${VERSION}${COMMIT}


all:	${NAME}


build:
	@echo 'Building ${NAME} version ${RELEASE}'

	go get -d -x -v .
	go build -o ${NAME} -x --ldflags "-s -extldflags='-static' -X main.version=${RELEASE}" .
	@echo 'Successfully built ${NAME}'



${NAME}: build


install:
	@echo 'PREFIX=${PREFIX}'

	@mkdir -p ${PREFIX}
	@mv ${NAME} ${PREFIX}/${NAME}
	@echo 'Successfully installed ${NAME} to ${PREFIX}'

clean:
	@rm -v ${NAME}

run:
	./phzd -conf config.toml.default -addr 0.0.0.0:8080
	