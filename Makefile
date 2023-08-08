test-cov:
	go test -v -coverprofile coverage.out \
		github.com/panupakm/miniredis/payload \
		github.com/panupakm/miniredis/client \
		github.com/panupakm/miniredis/server \
		github.com/panupakm/miniredis/request
	go tool cover -html=coverage.out
