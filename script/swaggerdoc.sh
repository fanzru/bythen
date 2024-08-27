#!/bin/bash

########
# Help #
########
Help() {
        # Display Help
        echo "Convert yaml inside api/http from openapi 3 to swagger 2."
        echo
        echo "Usage:"
        echo "  ./swaggerdoc.sh [options]"
        echo
        echo "options:"
        echo "h    Display Help"
        echo "b    Change base file (default: empty)"
        echo
}
files=""

while getopts ":hb:" flag;
do
        case "$flag" in
                h) Help
                   exit;;
                b) files="api/doc/swagger/$OPTARG.json";;
                \?) echo "Illegal option(s)"
                    exit;;
        esac
done

cnt=0

for file in ./api/http/*.yaml; do
        f=$(basename $file .yaml)
        mkdir -p ./api/doc/swagger

        api-spec-converter \
                --from=openapi_3 \
                --to=swagger_2 \
                api/http/$f.yaml > api/doc/swagger/$f.json

        files="$files api/doc/swagger/$f.json"
        cnt=$((cnt+1))
        if (( cnt == 1))
        then
                files="$files api/doc/swagger/$f.json"
        fi
done


mkdir -p ./docs/swagger
swagger -q mixin $files -o docs/swagger/docs.json