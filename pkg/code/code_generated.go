// Code generated by "codegen -type=int /root/go/src/eazycloud/pkg/code"; DO NOT EDIT.

package code

// init register error codes defines in this source code to errors package
func init() {
	register(ErrSuccess, 200, map[string]string{"MessageCN": "请求成功", "MessageEN": "Success."})
	register(ErrUnknown, 500, map[string]string{"MessageCN": "服务器出错", "MessageEN": "Internal server error."})
	register(ErrBind, 400, map[string]string{"MessageCN": "解析结构体出错", "MessageEN": "Error occurred while binding the request body to the struct."})
	register(ErrValidation, 400, map[string]string{"MessageCN": "参数校验失败", "MessageEN": "Validation failed."})
	register(ErrTokenInvalid, 401, map[string]string{"MessageCN": "令牌无效", "MessageEN": "Token invalid."})
	register(ErrPageNotFound, 404, map[string]string{"MessageCN": "请求路由不存在", "MessageEN": "Page not found."})
	register(ErrOperationBatchExecute, 200, map[string]string{"MessageCN": "批量执行操作", "MessageEN": "Operation batch execute."})
	register(ErrDatabase, 500, map[string]string{"MessageCN": "数据库出错", "MessageEN": "Database error."})
	register(ErrEncrypt, 401, map[string]string{"MessageCN": "用户密码加密失败", "MessageEN": "Error occurred while encrypting the user password."})
	register(ErrSignatureInvalid, 401, map[string]string{"MessageCN": "签名无效", "MessageEN": "Signature is invalid."})
	register(ErrExpired, 401, map[string]string{"MessageCN": "令牌", "MessageEN": "Token expired."})
	register(ErrInvalidAuthHeader, 401, map[string]string{"MessageCN": "无效的请求授权头部", "MessageEN": "Invalid authorization header."})
	register(ErrMissingHeader, 401, map[string]string{"MessageCN": "请求授权头部为空", "MessageEN": "The `Authorization` header was empty."})
	register(ErrPasswordIncorrect, 401, map[string]string{"MessageCN": "密码验证失败", "MessageEN": "Password was incorrect."})
	register(ErrPermissionDenied, 403, map[string]string{"MessageCN": "请求无权限执行", "MessageEN": "Permission denied."})
	register(ErrEncodingFailed, 500, map[string]string{"MessageCN": "数据编码出错", "MessageEN": "Encoding failed due to an error with the data."})
	register(ErrDecodingFailed, 500, map[string]string{"MessageCN": "数据解码出错", "MessageEN": "Decoding failed due to an error with the data."})
	register(ErrInvalidJSON, 500, map[string]string{"MessageCN": "数据非有效JSON结构", "MessageEN": "Data is not valid JSON."})
	register(ErrEncodingJSON, 500, map[string]string{"MessageCN": "JSON数据编码失败", "MessageEN": "JSON data could not be encoded."})
	register(ErrDecodingJSON, 500, map[string]string{"MessageCN": "JSON数据解码失败", "MessageEN": "JSON data could not be decoded."})
	register(ErrInvalidYaml, 500, map[string]string{"MessageCN": "数据非有效YAML结构", "MessageEN": "Data is not valid Yaml."})
	register(ErrEncodingYaml, 500, map[string]string{"MessageCN": "YAML数据编码失败", "MessageEN": "Yaml data could not be encoded."})
	register(ErrDecodingYaml, 500, map[string]string{"MessageCN": "YAML数据编码失败", "MessageEN": "Yaml data could not be decoded."})
	register(ErrHTTPError, 500, map[string]string{"MessageCN": "HTTP请求失败", "MessageEN": "HTTP request error."})
	register(ErrHTTPResponseDataParseError, 500, map[string]string{"MessageCN": "解析HTTP服务返回数据失败", "MessageEN": "Decode data from http response error."})
	register(ErrHTTPClientGenerateError, 500, map[string]string{"MessageCN": "生成HTTP客户端失败", "MessageEN": "Generate HTTP client error."})
	register(ErrGRPCClientGenerateError, 500, map[string]string{"MessageCN": "生成gRPC客户端失败", "MessageEN": "Generate gRPC client error."})
	register(ErrGRPCClientCertificateError, 500, map[string]string{"MessageCN": "gRPC客户端证书错误", "MessageEN": "Validate gRPC client certificate error."})
	register(ErrGRPCClientDialError, 500, map[string]string{"MessageCN": "gRPC客户端连接失败", "MessageEN": "Dial to gRPC server error."})
	register(ErrGRPCClientInvokeServiceError, 500, map[string]string{"MessageCN": "gRPC客户端访问服务接口失败", "MessageEN": "Invoke gRPC server service function error."})
	register(ErrGRPCResponseDataParseError, 500, map[string]string{"MessageCN": "解析gRPC服务返回数据失败", "MessageEN": "Decode data from gRPC service error."})
}
