migrate:
	go run ./cmd/migrate/main.go --migrate

seed:
	go run ./cmd/migrate/main.go --seed

fresh:
	go run ./cmd/migrate/main.go --fresh

consume:
	go run ./cmd/consume/main.go
