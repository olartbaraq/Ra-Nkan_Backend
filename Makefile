c_m:
	#create a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)

p_up:
	#create all services listed in docker compose file with docker
	docker compose up -d

p_down:
	#delete all services listed in docker compose file
	docker compose down

db_up:
	#create a database from the db server
	docker exec -it spectrumshelf_postgres createdb --username=root --owner=root spectrumshelf_db
	# docker exec -it ra_nkan_live createdb --username=root --owner=root ra_nkan_db

db_down:
	#delete a database from the db server
	docker exec -it spectrumshelf_postgres dropdb --username=root spectrumshelf_db
	# docker exec -it ra_nkan_live dropdb --username=root ra_nkan_db

dock_start:
	#start the docker processes
	docker start spectrumshelf_postgres
	docker start ra_nkan_live
	docker start ra_nkan_api
	docker start redis_live

dock_stop:
	#stop the docker processes
	docker stop spectrumshelf_postgres
	docker stop ra_nkan_live
	docker stop ra_nkan_api
	docker stop redis_live

m_up:
	#run a migration to the database
	migrate -path db/migrations -database "postgres://root:testing@localhost:5432/spectrumshelf_db?sslmode=disable" up
	# migrate -path db/migrations -database "postgres://root:testing@localhost:5433/ra_nkan_db?sslmode=disable" up

m_down:
	#revert the migration from the database
	migrate -path db/migrations -database "postgres://root:testing@localhost:5432/spectrumshelf_db?sslmode=disable" down
	# migrate -path db/migrations -database "postgres://root:testing@localhost:5433/ra_nkan_db?sslmode=disable" down

sqlc:
	#generate the sql queries to golang
	sqlc generate

build_image:
	#build project file to a docker image
	docker build -t ra_nkan:latest .

run_image:
	#docker command to run the project docker image
	# docker rm ra_nkan_api_image
	docker run --name ra_nkan_api_image -p 8000:8000 -e DB_SOURCE_LIVE="postgres://root:testing@172.17.0.1:5433/ra_nkan_db?sslmode=disable" ra_nkan:latest

# run_image_prod:
# 	#docker command to run the project docker image
# 	docker rm ra_nkan
# 	docker run --name ra_nkan --network ra_nkan_network -p 8000:8000 -e GIN_MODE=release -e DB_SOURCE_LIVE="postgres://root:testing@172.20.0.1:5433/ra_nkan_db?sslmode=disable" ra_nkan:latest

test:
	#run all tests in test directory
	go test -v -cover ./...
