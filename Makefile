generate:
	go generate ./...
	jet -source=sqlite -dsn="database.db" -schema=main -path=./internal/repo/jet
	go mod tidy

install-tools:
	go install github.com/go-jet/jet/v2/cmd/jet@latest
	go install github.com/jj-style/gonstructor/cmd/gonstructor@f908b05
	go install github.com/jmattheis/goverter/cmd/goverter@v1.4.0
