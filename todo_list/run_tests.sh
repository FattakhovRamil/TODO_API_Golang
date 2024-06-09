#!/bin/bash

# Запуск тестов и сохранение результатов
go test ./tests/db -coverpkg=./db > tests/test_results.txt
go test ./tests/handlers -coverpkg=./handlers >> tests/test_results.txt
go test ./tests/models -coverpkg=./models >> tests/test_results.txt

