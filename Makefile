test-cov:
	gocov test \
		github.com/panupakm/miniredis/client \
		github.com/panupakm/miniredis/internal/db \
		github.com/panupakm/miniredis/internal/pubsub \
		github.com/panupakm/miniredis/payload \
		github.com/panupakm/miniredis/request \
		github.com/panupakm/miniredis/server \
		github.com/panupakm/miniredis/server/internal/handler \
		| gocov report

test-cov-html:
	go test -coverprofile coverage-html.out \
		github.com/panupakm/miniredis/client \
		github.com/panupakm/miniredis/internal/db \
		github.com/panupakm/miniredis/internal/pubsub \
		github.com/panupakm/miniredis/payload \
		github.com/panupakm/miniredis/request \
		github.com/panupakm/miniredis/server \
		github.com/panupakm/miniredis/server/internal/handler 
	go tool cover -html=coverage-html.out

vul-check:
	govulncheck ./...