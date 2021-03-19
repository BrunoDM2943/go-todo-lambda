#!/bin/bash
GOOS=linux CGO_ENABLED=0 go build
zip function.zip go-todo-lambda