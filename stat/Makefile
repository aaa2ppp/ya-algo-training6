bin/stat: ./cmd/stat/*.go
	go test ./cmd/stat && go build -o bin/stat ./cmd/stat

bin/parse: ./cmd/parse/*.go
	go test ./cmd/parse && go build -o bin/parse ./cmd/parse

build: bin/stat bin/parse

stat: build
	./stat.sh

.PHONE: build stat
