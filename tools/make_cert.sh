#!/bin/bash
#
#

KEYNAME=dev-server.key
CERTNAME=dev-server.crt
FQDN=localhost

openssl genrsa -out $KEYNAME 2048
openssl ecparam -genkey -name secp384r1 -out $KEYNAME
openssl req -new -x509 -sha256 -key $KEYNAME -out $CERTNAME -days 3650 -subj "/C=US/ST=Washington/L=Seattle/O=int-main-void/OU=na/CN=$FQDN"

