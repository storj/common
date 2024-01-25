#!/usr/bin/env bash
if [ "$#" -ne 1 ]; then
    echo "Usage: ./gen-client.tls.sh <NAME>"
    exit 1
fi

NAME=$1
IP=$2
rm $NAME.key || true
rm $NAME.csr || true
rm $NAME.crt || true
openssl genrsa -out $NAME.key 4096
openssl req -new -key $NAME.key -out $NAME.csr -subj "/CN=$NAME"
openssl x509 -req \
   -days 3650 \
   -copy_extensions copy \
   -in $NAME.csr \
   -CA ca-client.crt \
   -CAkey ca-client.key \
   -CAcreateserial \
   -out $NAME.crt
