#!/usr/bin/env bash 

openssl genrsa -out ../certs/server.key 2048

# Common Name should be ip address or hostname tied to the certificate 
openssl req -new -x509 -sha256 -key ../certs/server.key -out ../certs/server.crt -days 3650

openssl req -new -sha256 -key ../certs/server.key -out ../certs/server.csr

openssl x509 -req -sha256 -in ../certs/server.csr -signkey ../certs/server.key -out ../certs/server.crt -days 3650