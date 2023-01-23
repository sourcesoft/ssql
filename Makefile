GOCMD=go
GOTEST=$(GOCMD) test
GOCHECK=$(GOCMD) vet
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOMIGRATIONRUN=migration

doc:
	godoc -http=localhost:6060 & xdg-open http://localhost:6060/pkg/github.com/sourcesoft/ssql/
check:
	@echo "-------- ( ͡° ͜ʖ ͡°) --------"
	@$(GOCHECK) $(shell $(GOTEST_PATH_ALL))

test:
	@echo "-------- ( ͡° ͜ʖ ͡°) --------"
	@$(GOTEST) $(shell $(GOTEST_PATH_ALL))

db-migration:
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN) bootstrap

db-flush:
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN) flush

db-up:
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN) up

db-down:
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN) down

db-seed:
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN)

db-create:
	$(GOMIGRATIONRUN) createdb

db-local-init:
	make install-migration && \
	docker stop postgresql && \
	docker rm postgresql && \
	docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=dev -p 5432:5432 -v /data:/var/lib/postgresql/data -d postgres:14 && \
	sleep 3 && \
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN) createdb && \
	make db-migration

db-local-reset:
	docker stop postgresql && \
	docker rm postgresql && \
	docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=dev -p 5432:5432 -v /data:/var/lib/postgresql/data -d postgres:14 && \
	sleep 2 && \
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN) deletedb && \
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
	SSQL_DDL=scripts/migration/_sample_ddl SSQL_SEEDS=scripts/migration/_sample_seeds/seeds.sql \
	$(GOMIGRATIONRUN) createdb && \
	make db-migration

install-migration:
	cd scripts/migration && go install

example-simple:
	@cd _examples/simple && \
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
    go run *.go

example-find:
	@cd _examples/find && \
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
    go run *.go
	
example-cursor:
	@cd _examples/cursor && \
	SSQL_DB=ssql SSQL_USERNAME=postgres SSQL_PASSWORD=dev SSQL_HOST=localhost SSQL_PORT=5432 \
    go run *.go