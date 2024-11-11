# jocp

- docker compose exec -it postgres psql -U root -d test_db

# goose commands

- go install github.com/pressly/goose/v3/cmd/goose@v3

cd migrations
- > goose create users sql
- > goose postgres "host=localhost port=5432 user=root password=secret dbname=test_db sslmode=disable" status
- > goose postgres "host=localhost port=5432 user=root password=secret dbname=test_db sslmode=disable" up
- > goose postgres "host=localhost port=5432 user=root password=secret dbname=test_db sslmode=disable" down
- > goose postgres "host=localhost port=5432 user=root password=secret dbname=test_db sslmode=disable" reset

# Set go path

- go env GOPATH
- export PATH=$PATH:$(go env GOPATH)/bin
- source ~/.bashrc


ON 137 - signout# jocp.xyz
# jocp.xyz
