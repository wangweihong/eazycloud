#!/usr/bin/env bash
# gRPC  server must enable reflect and debug service
# 可使用-plaintext 忽略证书检测
grpcurl --cacert=./_output/cert/ca.crt debug.DebugService/Panic