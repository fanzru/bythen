include .env

run:
	go run ./cmd/main.go
migrate:
	migrate -database "mysql://$(USERDB):$(PASSDB)@tcp(127.0.0.1:3306)/tix_refund" -path mysqldb/migration up
rollback:
	migrate -database "mysql://$(USERDB):$(PASSDB)@tcp(127.0.0.1:3306)/tix_refund" -path mysqldb/migration down $(N)
create:
	migrate create -ext sql -dir mysqldb/migration -seq $(NAME)

generate-http:
	chmod +x ./script/http.sh
	./script/http.sh

generate-swaggerdoc:
	chmod +x ./script/swaggerdoc.sh
	./script/swaggerdoc.sh

# Database Migration
# create migration file
migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

# migration up (craete all table)
migrate-up:
	migrate -database "mysql://$(MYSQL_DBUSER):$(MYSQL_DBPASS)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DBNAME)" -path mysqldb/migration up

# migration down (drop all table)
migrate-down:
	migrate -database "mysql://$(MYSQL_DBUSER):$(MYSQL_DBPASS)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DBNAME)" -path mysqldb/migration down -all

# rollback migration
migrate-rollback:
	migrate -database "mysql://$(MYSQL_DBUSER):$(MYSQL_DBPASS)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DBNAME)" -path mysqldb/migration down $(N)

# migration force with version (craete all table)
migrate-force:
	migrate -database "mysql://$(MYSQL_DBUSER):$(MYSQL_DBPASS)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DBNAME)" -path mysqldb/migration force $(VERSION)

mock-usecase:
	mockgen -package=mock -source=app/${SUBDOMAIN}/usecase/${FILE}.go -destination=app/${SUBDOMAIN}/usecase/usecase_mock/${FILE}.go

mock-repo:
	mockgen -package=mock -source=app/${SUBDOMAIN}/repo/${FILE}.go  -destination=app/${SUBDOMAIN}/repo/repo_mock/${FILE}.go