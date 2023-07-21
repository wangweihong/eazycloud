#!/usr/bin/env bash

components="$1"
if [ $# -ne 1 ];then
    lib::log::error "Usage: generate_default_config.sh example-server example-cli"
    exit 1
fi
# 这段脚本的作用是调用./genconfig.sh将指定组件../configs/目录上各个配置文件模板变量值
# 替换成install/environment.sh中定义的默认值

# 代码根目录
SOURCE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${SOURCE_ROOT}/scripts/common.sh"

# 配置文件生成路径. 默认是 ${SOURCE_ROOT}/_output/configs目录下
readonly LOCAL_OUTPUT_CONFIGPATH="${LOCAL_OUTPUT_ROOT}/configs"
mkdir -p ${LOCAL_OUTPUT_CONFIGPATH}

cd ${SOURCE_ROOT}/scripts

for comp in $components
do
  lib::log::info "generate ${LOCAL_OUTPUT_CONFIGPATH}/${comp}.yaml"
  ./genconfig.sh install/environment.sh ../configs/${comp}.yaml > ${LOCAL_OUTPUT_CONFIGPATH}/${comp}.yaml
done
