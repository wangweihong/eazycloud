#!/usr/bin/env bash

# 项目源码根目录
SOURCE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# 生成文件存放目录
LOCAL_OUTPUT_ROOT="${SOURCE_ROOT}/${OUT_DIR:-_output}"


# 设置安装目录
readonly INSTALL_DIR=${INSTALL_DIR:-/tmp/installation}
mkdir -p ${INSTALL_DIR}
readonly ENV_FILE=${SOURCE_ROOT}/scripts/install/environment.sh

# eazycloud 配置
readonly EAZYCLOUD_DATA_DIR=${EAZYCLOUD_DATA_DIR:-/var/lib/eazycloud/data} # eazycloud 各组件数据目录
readonly EAZYCLOUD_INSTALL_DIR=${EAZYCLOUD_INSTALL_DIR:-/var/lib/eazycloud/bin} # eazycloud 安装文件存放目录
readonly EAZYCLOUD_CONFIG_DIR=${EAZYCLOUD_CONFIG_DIR:-/var/lib/eazycloud/conf} # eazycloud 配置文件存放目录
readonly EAZYCLOUD_LOG_DIR=${EAZYCLOUD_LOG_DIR:-/var/log/eazycloud} # eazycloud 日志文件存放目录
readonly CA_FILE=${CA_FILE:-${EAZYCLOUD_CONFIG_DIR}/cert/ca.pem} # ca

# example-server 配置
readonly EXAMPLE_SERVER_HOST=${EXAMPLE_SERVER_HOST:-127.0.0.1} # eazycloud-apiserver 部署机器 IP 地址
readonly EXAMPLE_SERVER_GRPC_BIND_ADDRESS=${EXAMPLE_SERVER_GRPC_BIND_ADDRESS:-0.0.0.0}
readonly EXAMPLE_SERVER_GRPC_BIND_PORT=${EXAMPLE_SERVER_GRPC_BIND_PORT:-8081}
readonly EXAMPLE_SERVER_INSECURE_BIND_ADDRESS=${EXAMPLE_SERVER_INSECURE_BIND_ADDRESS:-127.0.0.1}
readonly EXAMPLE_SERVER_INSECURE_BIND_PORT=${EXAMPLE_SERVER_INSECURE_BIND_PORT:-8080}
readonly EXAMPLE_SERVER_SECURE_BIND_ADDRESS=${EXAMPLE_SERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly EXAMPLE_SERVER_SECURE_BIND_PORT=${EXAMPLE_SERVER_SECURE_BIND_PORT:-8443}
readonly EXAMPLE_SERVER_SECURE_TLS_CERT_FILE=${EXAMPLE_SERVER_SECURE_TLS_CERT_FILE:-${EAZYCLOUD_CONFIG_DIR}/cert/example-server.crt}
readonly EXAMPLE_SERVER_SECURE_TLS_CERT_KEY=${EXAMPLE_SERVER_SECURE_TLS_CERT_KEY:-${EAZYCLOUD_CONFIG_DIR}/cert/example-server.key}

