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