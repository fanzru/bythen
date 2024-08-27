#!/bin/bash
for file in ./api/http/*.yaml; do
        f=$(basename $file .yaml)
        mkdir -p ./internal/app/$f/port/genhttp


        oapi-codegen -generate types  \
        -o ./internal/app/$f/port/genhttp/openapi_types.gen.go \
        -package genhttp \
        api/http/$f.yaml 

        oapi-codegen -generate std-http \
        -o ./internal/app/$f/port/genhttp/openapi_server.gen.go \
        -package genhttp \
        api/http/$f.yaml 
done