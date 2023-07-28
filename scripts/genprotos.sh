#!/usr/bin/env bash

SOURCE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${SOURCE_ROOT}/scripts/lib/init.sh"

readonly  exclude_files=()
# Args:
#   $1 (the directory that proto files )
function generate_protos()
{
 local proto_dirs=${1}

 pushd ${proto_dirs}
 # remove old pb file
 find ./ -name "*.pb.go" |xargs rm -rf

 # 生成新的pb.go文件
 for pbfile in `find  -name "*.proto"` ; do
   echo $pbfile
   # 注意生成的pb包路径为 --go_out路径+.pb包中go_package变量指定的路径
   # 如version.pb文件中go_package值为`option go_package = "apis/version";`
   # 其生成的version.pb.go路径为../apis/version/version.pb.go
   protoc --go_out=plugins=grpc:../ $pbfile
 done

 popd
}

# 加上这句后, 就允许外部通过<脚本名> <函数名> 来直接调用脚本内的函数
$*