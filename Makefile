
#CORE
.PHONY: core-build
core-build:
	docker-compose -f botes/core/docker-compose.yaml --env botes/core/.env build

.PHONY: core-run
core-run:
	docker-compose -f botes/core/docker-compose.yaml --env botes/core/.env up

.PHONY: core-run-d
core-run-d:
	docker-compose -f botes/core/docker-compose.yaml --env botes/core/.env up -d

.PHONY: core-stop
core-stop:
	docker-compose -f botes/core/docker-compose.yaml --env botes/core/.env down

#ADMIN
.PHONY: admin-build
admin-build:
	docker-compose -f botes/admin/docker-compose.yaml --env botes/admin/.env build

.PHONY: admin-run
admin-run:
	docker-compose -f botes/admin/docker-compose.yaml --env botes/admin/.env up

.PHONY: admin-run-d
admin-run-d:
	docker-compose -f botes/admin/docker-compose.yaml --env botes/admin/.env up -d

.PHONY: admin-stop
admin-stop:
	docker-compose -f botes/admin/docker-compose.yaml --env botes/admin/.env down

#ALL
.PHONY: build-all
build-all : core-build admin-build

.PHONY: run-all
run-all : core-run-d admin-run-d

.PHONY: stop-all
stop-all : core-stop admin-stop

