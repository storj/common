#!/usr/bin/env bash
if [ "$#" -ne 2 ]; then
    echo "Usage: ./gen-host-tls.sh <HOSTNAME> <IP>"
    exit 1
fi
HOST=$1
IP=$2
rm $HOST.csr
rm $HOST.crt
EXT="subjectAltName=DNS:$HOST,DNS:localhost,IP:127.0.0.1,IP:$IP"
openssl genrsa -out $HOST.key 4096
openssl req -new -key $HOST.key -out $HOST.csr \
   -subj "/CN=$HOST" \
   -addext "$EXT"
openssl x509 -req \
   -days 3650 \
   -copy_extensions copy \
   -in $HOST.csr \
   -CA ca.crt \
   -CAkey ca.key \
   -CAcreateserial \
   -out $HOST.crt
