#!/usr/bin/env bash

SOURCE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${SOURCE_ROOT}/scripts/lib/init.sh"

# OUT_DIR can come in from the Makefile, so honor it.
readonly LOCAL_OUTPUT_ROOT="${SOURCE_ROOT}/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_CAPATH="${LOCAL_OUTPUT_ROOT}/cert"

readonly CERT_HOSTNAME="${CERT_HOSTNAME:-example-server.exzycloud.com},127.0.0.1,localhost"


# Args:
#   $1 (the directory that certificate files to save)
#   $2 (the prefix of the certificate filename)
#   $3 (cert alert name)
#   $4 (cert subject)
function generate_certificate()
{
 local cert_dir=${1}
 local prefix=${2}
 local cert_hostname=${3}
 local cert_subject=${4}

 mkdir -p "${cert_dir}"

 # 确认openssl是否安装
 lib::util::test_openssl_installed
 # 将当前路径入栈,并跳转到证书目录
 pushd "${cert_dir}"

 if [ ! -r "ca.crt" ]; then
   lib::log::info "ca.crt not exist, trying to generate ca.art"
   ${OPENSSL_BIN} genrsa -out ca.key 4096
   ${OPENSSL_BIN} req -x509 -new -nodes -sha512 -days 3650 \
      -subj "$cert_subject" \
      -key ca.key \
      -out ca.crt
 fi

 lib::log::info "Generate "${prefix}" certificates...."
 # server cert
 ${OPENSSL_BIN} genrsa -out ${prefix}.key 4096
 ${OPENSSL_BIN} req -sha512 -new \
     -subj "$cert_subject" \
     -key ${prefix}.key \
     -out ${prefix}.csr

 v3ExtFILE=${prefix}_v3.ext

cat > ${v3ExtFILE} <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
EOF

# 按,切割证书主体
IFS=',' read -ra elements <<< "${cert_hostname}"

# 使用循环遍历主体，生成证书主体
 j=0
 for (( i=0; i<${#elements[@]}; i++ )); do
  element="${elements[$i]}"
  if [[ $element =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # 如果是IP地址，则给 IP.* 赋值
    echo "IP.$((j=j+1)) = $element" >> ${v3ExtFILE}
  fi
 done

 j=0
 for (( i=0; i<${#elements[@]}; i++ )); do
  element="${elements[$i]}"
  if [[ ! $element =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # 否则，给 DNS.* 赋值
    echo "DNS.$((j=j+1)) = $element" >> ${v3ExtFILE}
  fi
 done

 ${OPENSSL_BIN} x509 -req -sha512 -days 3650 \
     -extfile ${v3ExtFILE}  \
     -CA ca.crt -CAkey ca.key -CAcreateserial \
     -in ${prefix}.csr \
     -out ${prefix}.crt

 # 跳回到上一次入栈的路径
 popd
}

# 加上这句后, 就允许外部通过<脚本名> <函数名> 来直接调用脚本内的函数
$*