test-cov:
	gocov test \
		github.com/panupakm/miniredis/client \
		github.com/panupakm/miniredis/server/storage \
		github.com/panupakm/miniredis/server/pubsub \
		github.com/panupakm/miniredis/payload \
		github.com/panupakm/miniredis/request \
		github.com/panupakm/miniredis/server \
		github.com/panupakm/miniredis/server/internal/handler \
		| gocov report

test-cov-html:
	go test -coverprofile coverage-html.out \
		github.com/panupakm/miniredis/client \
		github.com/panupakm/miniredis/server/storage \
		github.com/panupakm/miniredis/server/pubsub \
		github.com/panupakm/miniredis/payload \
		github.com/panupakm/miniredis/request \
		github.com/panupakm/miniredis/server \
		github.com/panupakm/miniredis/server/internal/handler
	go tool cover -html=coverage-html.out

test-integration:
	go test ./tests

build: 
	go build -o bin/server ./cmd/server/main.go
	go build -o bin/client ./cmd/client/main.go

vul-check:
	govulncheck ./...