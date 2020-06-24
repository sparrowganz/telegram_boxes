
core-compose = botes/core/docker-compose.yaml
core-env = botes/core/.env

box-compose= botes/boxes/docker-compose.yaml
box-env = botes/boxes/Test/.env

mongo-compose = botes/mongo/docker-compose.yaml
mongo-env = botes/mongo/.env

build = build
up = up
start = up -d
stop = down

#-----------------------------------------------------------------------------------------------------------------------
#MONGO
#-----------------------------------------------------------------------------------------------------------------------

.PHONY: mongo-run run-mongo
mongo-run run-mongo:
	docker-compose -f $(mongo-compose) --env $(mongo-env) $(up)

.PHONY: mongo-run-d run-mongo-d
mongo-run-d run-mongo-d:
	docker-compose -f $(mongo-compose) --env $(mongo-env) $(start)

.PHONY: mongo-stop stop-mongo
mongo-stop stop-mongo:
	docker-compose -f $(mongo-compose) --env $(mongo-env) $(stop)

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
	docker-compose -f $(core-compose) --env $(core-env) $(stop)

#-----------------------------------------------------------------------------------------------------------------------
#BOX
#-----------------------------------------------------------------------------------------------------------------------

.PHONY: box-build build-box
box-build build-box:
	docker-compose -f $(box-compose) --env $(box-env) $(build)

.PHONY: box-run run-box
box-run run-box:
	docker-compose -f $(box-compose) --env $(box-env) $(up)

.PHONY: box-run-d run-box-d
box-run-d run-box-d:
	docker-compose -f $(box-compose) --env $(box-env) $(start)

.PHONY: box-stop stop-box
box-stop stop-box:
	docker-compose -f $(box-compose) --env $(box-env) $(stop)


#-----------------------------------------------------------------------------------------------------------------------
#ALL
#-----------------------------------------------------------------------------------------------------------------------

.PHONY: build-all all-build
build-all all-build: core-build box-build

.PHONY: run-all all-run
run-all all-run : core-run-d box-run-d

.PHONY: stop-all all-stop
stop-all all-stop: core-stop  box-stop

#-----------------------------------------------------------------------------------------------------------------------
#Clean
#-----------------------------------------------------------------------------------------------------------------------

PHONY: clean
clean: stop-all remove prune

remove:
	docker rmi $$(docker images -a)

prune:
	 docker system prune