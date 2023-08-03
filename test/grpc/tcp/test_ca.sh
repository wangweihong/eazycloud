#!/usr/bin/env bash

grpcurl --cacert=./_output/cert/ca.crt  localhost:8081 version.VersionService/Version