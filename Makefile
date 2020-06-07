
core-compose = botes/core/docker-compose.yaml
core-env = botes/core/.env

admin-compose = botes/admin/docker-compose.yaml
admin-env = botes/admin/.env

build = build
up = up
start = up -d
stop = down

#-----------------------------------------------------------------------------------------------------------------------
#CORE
#-----------------------------------------------------------------------------------------------------------------------

.PHONY: core-build build-core
core-build build-core:
	docker-compose -f $(core-compose) --env $(core-env) $(build)

.PHONY: core-run run-core
core-run run-core:
	docker-compose -f $(core-compose) --env $(core-env) $(up)

.PHONY: core-run-d run-core-d
core-run-d run-core-d:
	docker-compose -f $(core-compose) --env $(core-env) $(start)

.PHONY: core-stop stop-core
core-stop stop-core:
	docker-compose -f $(core-compose) --env $(core-env) $(down)

#-----------------------------------------------------------------------------------------------------------------------
#ADMIN
#-----------------------------------------------------------------------------------------------------------------------

.PHONY: admin-build build-admin
admin-build build-admin:
	docker-compose -f $(admin-compose) --env $(admin-env) $(build)

.PHONY: admin-run run-admin
admin-run run-admin:
	docker-compose -f $(admin-compose) --env $(admin-env) $(up)

.PHONY: admin-run-d run-admin-d
admin-run-d run-admin-dd:
	docker-compose -f $(admin-compose) --env $(admin-env) $(start)

.PHONY: admin-stop stop-admin
admin-stop stop-admin:
	docker-compose -f $(admin-compose) --env $(admin-env) $(stop)

#-----------------------------------------------------------------------------------------------------------------------
#ALL
#-----------------------------------------------------------------------------------------------------------------------

.PHONY: build-all all-build
build-all all-build: core-build admin-build

.PHONY: run-all all-run
run-all all-run : core-run-d admin-run-d

.PHONY: stop-all all-stop
stop-all all-stop: core-stop admin-stop
