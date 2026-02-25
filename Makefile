migrate-create:
	migrate create -ext=sql -dir=./internal/postgres/migrations $(name)
