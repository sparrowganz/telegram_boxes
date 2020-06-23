#!/usr/bin/env bash

#----------------------------------------------------------------------------------
#   Core
#----------------------------------------------------------------------------------

echo "Генерация сервера Core"

echo "Генерация сервера Servers"
protoc -I$GOPATH/src/telegram_boxes \
    services/core/protobuf/servers.proto --go_out=plugins=grpc:.

echo "Генерация сервера Tasks"
protoc -I$GOPATH/src/telegram_boxes \
    services/core/protobuf/tasks.proto --go_out=plugins=grpc:.

echo "------- Генерация клиента сервиса Logger"
protoc -I$GOPATH/src/telegram_boxes \
    services/logs/protobuf/logs.proto --go_out=plugins=grpc:./services/core/protobuf

#----------------------------------------------------------------------------------
#   Box
#----------------------------------------------------------------------------------
echo "Генерация сервера Box"

echo "------- Генерация клиента сервиса Tasks"
protoc -I$GOPATH/src/telegram_boxes \
    services/core/protobuf/tasks.proto --go_out=plugins=grpc:./services/box/protobuf

echo "------- Генерация клиента сервиса Logger"
protoc -I$GOPATH/src/telegram_boxes \
    services/logs/protobuf/logs.proto --go_out=plugins=grpc:./services/box/protobuf

echo "------- Генерация клиента сервиса Servers"
protoc -I$GOPATH/src/telegram_boxes \
    services/core/protobuf/servers.proto --go_out=plugins=grpc:./services/box/protobuf



#----------------------------------------------------------------------------------
#   Admin
#----------------------------------------------------------------------------------
echo "Генерация сервера Admin"

echo "------- Генерация клиента сервиса Logger"
protoc -I$GOPATH/src/telegram_boxes \
    services/logs/protobuf/logs.proto --go_out=plugins=grpc:./services/admin/protobuf

#----------------------------------------------------------------------------------
#   Logger
#----------------------------------------------------------------------------------

echo "Генерация сервера Log"
protoc -I$GOPATH/src/telegram_boxes \
    services/logs/protobuf/logs.proto --go_out=plugins=grpc:.