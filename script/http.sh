#!/bin/bash
for file in ./api/http/*.yaml; do
    f=$(basename $file .yaml)
    mkdir -p ./internal/app/$f/port/genhttp

    # Generate types with skip-prune option
    oapi-codegen -generate types,skip-prune \
        -o ./internal/app/$f/port/genhttp/openapi_types.gen.go \
        -package genhttp \
        api/http/$f.yaml 

    # Generate std-http with strict-server option
    oapi-codegen -generate std-http \
        -o ./internal/app/$f/port/genhttp/openapi_server.gen.go \
        -package genhttp \
        api/http/$f.yaml 
done
