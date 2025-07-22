#!/bin/bash
go run ../cmd
docker build -t gomdhtml-nginx .
docker run -d -p 80:80 gomdhtml-nginx