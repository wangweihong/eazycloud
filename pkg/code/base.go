package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ../../../docs/guide/zh-CN/api/error_code_generated.md

// Common: basic errors.
// Code must start with 1xxxxx.
const (

	// @HTTP 200
	// @MessageCN 请求成功
	// @MessageEN Success.
	ErrSuccess int = iota + 100001

	// @HTTP 500
	// @MessageCN 服务器出错
	// @MessageEN Internal server error.
	ErrUnknown

	// @HTTP 400
	// @MessageCN 解析结构体出错
	// @MessageEN Error occurred while binding the request body to the struct.
	ErrBind

	// @HTTP 400
	// @MessageCN  参数校验失败
	// @MessageEN  Validation failed.
	ErrValidation

	// @HTTP 401
	// @MessageCN  令牌无效
	// @MessageEN  Token invalid.
	ErrTokenInvalid

	// @HTTP 404
	// @MessageCN  请求路由不存在
	// @MessageEN  Page not found.
	ErrPageNotFound

	// ErrOperationBatchExecute 表明当前操作为批量操作,需解析结构确认批量结果
	// @HTTP 200
	// @MessageCN  批量执行操作
	// @MessageEN  Operation batch execute.
	ErrOperationBatchExecute
)

// common: database errors.
const (
	// @HTTP 500
	// @MessageCN  数据库出错
	// @MessageEN  Database error.
	ErrDatabase int = iota + 100101
)

// common: authorization and authentication errors.
const (
	// @HTTP 401
	// @MessageCN  用户密码加密失败
	// @MessageEN  Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201

	// @HTTP 401
	// @MessageCN  签名无效
	// @MessageEN  Signature is invalid.
	ErrSignatureInvalid

	// @HTTP 401
	// @MessageCN  令牌
	// @MessageEN  Token expired.
	ErrExpired

	// @HTTP 401
	// @MessageCN  无效的请求授权头部
	// @MessageEN  Invalid authorization header.
	ErrInvalidAuthHeader

	// @HTTP 401
	// @MessageCN  请求授权头部为空
	// @MessageEN  The `Authorization` header was empty.
	ErrMissingHeader

	// @HTTP 401
	// @MessageCN  密码验证失败
	// @MessageEN  Password was incorrect.
	ErrPasswordIncorrect

	// @HTTP 403
	// @MessageCN  请求无权限执行
	// @MessageEN  Permission denied.
	ErrPermissionDenied
)

// common: encode/decode errors.
const (
	// @HTTP 500
	// @MessageCN  数据编码出错
	// @MessageEN  Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301

	// @HTTP 500
	// @MessageCN  数据解码出错
	// @MessageEN  Decoding failed due to an error with the data.
	ErrDecodingFailed

	// @HTTP 500
	// @MessageCN  数据非有效JSON结构
	// @MessageEN   Data is not valid JSON.
	ErrInvalidJSON

	// @HTTP 500
	// @MessageCN  JSON数据编码失败
	// @MessageEN  JSON data could not be encoded.
	ErrEncodingJSON

	// @HTTP 500
	// @MessageCN  JSON数据解码失败
	// @MessageEN  JSON data could not be decoded.
	ErrDecodingJSON

	// @HTTP 500
	// @MessageCN  数据非有效YAML结构
	// @MessageEN  Data is not valid Yaml.
	ErrInvalidYaml

	// @HTTP 500
	// @MessageCN  YAML数据编码失败
	// @MessageEN  Yaml data could not be encoded.
	ErrEncodingYaml

	// @HTTP 500
	// @MessageCN  YAML数据编码失败
	// @MessageEN  Yaml data could not be decoded.
	ErrDecodingYaml
)

// common: Http  server error.
const ()

// common: Http  client error.
const (
	// @HTTP 500
	// @MessageCN  HTTP请求失败
	// @MessageEN  HTTP request error.
	ErrHTTPError int = iota + 100501

	// @HTTP 500
	// @MessageCN  解析HTTP服务返回数据失败
	// @MessageEN  Decode data from http response error.
	ErrHTTPResponseDataParseError

	// @HTTP 500
	// @MessageCN  生成HTTP客户端失败
	// @MessageEN  Generate HTTP client error.
	ErrHTTPClientGenerateError
)

// common: gRPC  server error.
const ()

// common: gRPC  client error.
const (
	// @HTTP 500
	// @MessageCN  生成gRPC客户端失败
	// @MessageEN  Generate gRPC client error.
	ErrGRPCClientGenerateError int = iota + 100701

	// @HTTP 500
	// @MessageCN  gRPC客户端证书错误
	// @MessageEN   Validate gRPC client certificate error.
	ErrGRPCClientCertificateError

	// @HTTP 500
	// @MessageCN  gRPC客户端连接失败
	// @MessageEN   Dial to gRPC server error.
	ErrGRPCClientDialError

	// @HTTP 500
	// @MessageCN  gRPC客户端访问服务接口失败
	// @MessageEN   Invoke gRPC server service function error.
	ErrGRPCClientInvokeServiceError

	// @HTTP 500
	// @MessageCN  解析gRPC服务返回数据失败
	// @MessageEN  Decode data from gRPC service error.
	ErrGRPCResponseDataParseError
)
