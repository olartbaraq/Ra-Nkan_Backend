c_m:
	#create a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)

p_up:
	#create postgres server with docker
	docker compose up -d

p_down:
	#delete postgres server
	docker compose down

db_up:
	#create a database from the db server
	docker exec -it spectrumshelf_postgres createdb --username=root --owner=root spectrumshelf_db
	

db_down:
	#delete a database from the db server
	docker exec -it spectrumshelf_postgres dropdb --username=root spectrumshelf_db

m_up:
	#run a migration to the database
	migrate -path db/migrations -database "postgres://root:testing@localhost:5432/spectrumshelf_db?sslmode=disable" up

m_down:
	#revert the migration from the database
	migrate -path db/migrations -database "postgres://root:testing@localhost:5432/spectrumshelf_db?sslmode=disable" down

sqlc:
	#generate the sql queries to golang
	sqlc generate

sqlc_win:
	#generate the sql queries to golang for windows
	docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate