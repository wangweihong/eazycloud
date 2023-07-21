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

	// ErrValidation - 400: Validation failed.
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

// example: policy errors.
const (

	// @HTTP 404
	// @MessageCN 策略未找到
	// @MessageEN Policy not found.
	ErrPolicyNotFound int = iota + 110201
)

// example: user errors.
const (

	// @HTTP 404
	// @MessageCN 用户未找到
	// @MessageEN User not found.
	ErrUserNotFound int = iota + 110202
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
	// @MessageCN  YAML数据编码
	// @MessageEN  Yaml data could not be decoded.
	ErrDecodingYaml
)
