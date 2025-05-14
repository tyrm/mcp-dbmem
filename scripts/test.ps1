go test -coverprofile coverage.out -cover ./...
go tool cover -html coverage.out -o coverage.html

Start-Process "./coverage.html"
