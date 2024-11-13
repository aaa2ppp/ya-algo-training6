bin/stat: ./cmd/stat/main.go
	go build -o bin/stat ./cmd/stat/main.go

bin/parse: ./cmd/parse/main.go
	go build -o bin/parse ./cmd/parse/main.go

build: bin/stat bin/parse

stat: build
	./stat.sh

.PHONE: build stat
