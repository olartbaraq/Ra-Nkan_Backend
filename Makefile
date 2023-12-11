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
	docker exec -it ra_nkan_live createdb --username=root --owner=root ra_nkan_db

db_down:
	#delete a database from the db server
	docker exec -it spectrumshelf_postgres dropdb --username=root spectrumshelf_db
	docker exec -it ra_nkan_live dropdb --username=root ra_nkan_db

dock_start:
	#start the docker processes
	docker start spectrumshelf_postgres
	docker start ra_nkan_live

dock_stop:
	#stop the docker processes
	docker stop spectrumshelf_postgres
	docker stop ra_nkan_live

m_up:
	#run a migration to the database
	migrate -path db/migrations -database "postgres://root:testing@localhost:5432/spectrumshelf_db?sslmode=disable" up
	migrate -path db/migrations -database "postgres://root:testing@localhost:5433/ra_nkan_db?sslmode=disable" up

m_down:
	#revert the migration from the database
	migrate -path db/migrations -database "postgres://root:testing@localhost:5432/spectrumshelf_db?sslmode=disable" down
	migrate -path db/migrations -database "postgres://root:testing@localhost:5433/ra_nkan_db?sslmode=disable" down

sqlc:
	#generate the sql queries to golang
	sqlc generate

test:
	#run all tests in test directory
	go test -v -cover ./...
