# phzd
# Makefile generated by 'gomaker' 

NAME ?= phzd
VERSION ?= 

# with a common install location, scripts should work on most platforms.
#
# default shebang: #!/usr/local/bin/phz
PREFIX ?= /usr/local/bin
VER ?= X
COMMIT != git rev-parse --verify --short HEAD
COMMIT ?= ${VER}
RELEASE ?= ${VERSION}${COMMIT}
all:	${NAME} phz-cli

# build web server
build: phzd
phzd:
	@echo 'Building ${NAME} version ${RELEASE}'
	go get -d -x -v .
	go build -o ${NAME} -x --ldflags "-s -extldflags='-static' -X main.version=${RELEASE}" ./cmd/phzd
	@echo 'Successfully built ${NAME}'
phz-cli: build
	go build -o phz-cli -x --ldflags "-s -extldflags='-static' -X main.version=${RELEASE}" ./cmd/phz-cli

# install web server
install-phzd:
	@echo 'PREFIX=${PREFIX}'
	@mkdir -p ${PREFIX}
	@install ${NAME} ${PREFIX}/${NAME}
	@echo 'Successfully installed ${NAME} to ${PREFIX}'
clean:
	@rm -fv ${NAME} phz-cli

test: build
	go test -v ./phz/...
	go test -v ./cmd/...
	@echo Running the phz command line interpreter:
	./phz-cli test.phz

# for heroku
run-demo: build
	./phzd -conf config.toml.default -addr 0.0.0.0:80 -v

run: build
	env -i myapikey=secret11223344 PATH=/usr/bin:/bin:/usr/local/bin ./phzd -conf config.toml.default -addr 0.0.0.0:8080 -v

# install phz-cli to /usr/local/bin/phz
install: install-phzd
	@echo 'PREFIX=${PREFIX}'
	@mkdir -p ${PREFIX}
	@install phz-cli ${PREFIX}/phz
	@echo "Successfully installed phz interpreter"
	@echo "Now you can use '#!/usr/local/bin/phz' in your scripts!"
	@echo "Try some of the executable scripts in testdata, such as testdata/phzzr"
	
