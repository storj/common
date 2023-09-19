#!/usr/bin/env bash
NAME=$1
openssl genrsa -out $1.key 4096
openssl req -new -x509 -days 3650 -key $1.key -out $1.crt
