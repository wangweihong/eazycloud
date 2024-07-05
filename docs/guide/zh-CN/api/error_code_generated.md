# 错误码

！！系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "EN": "Database error",
  "CN": "数据库出错"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |  中文描述 	|
| ---------- | ---- | --------- | ----------- | ----------- |
| ErrSuccess | 100001 | 200 | Unset error message | 错误信息未设置 |
| ErrUnknown | 100002 | 500 | Unset error message | 错误信息未设置 |
| ErrBind | 100003 | 400 | Unset error message | 错误信息未设置 |
| ErrValidation | 100004 | 400 | Unset error message | 错误信息未设置 |
| ErrTokenInvalid | 100005 | 401 | Unset error message | 错误信息未设置 |
| ErrPageNotFound | 100006 | 404 | Unset error message | 错误信息未设置 |
| ErrOperationBatchExecute | 100007 | 200 | Unset error message | 错误信息未设置 |
| ErrDatabase | 100101 | 500 | Unset error message | 错误信息未设置 |
| ErrEncrypt | 100201 | 401 | Unset error message | 错误信息未设置 |
| ErrSignatureInvalid | 100202 | 401 | Unset error message | 错误信息未设置 |
| ErrExpired | 100203 | 401 | Unset error message | 错误信息未设置 |
| ErrInvalidAuthHeader | 100204 | 401 | Unset error message | 错误信息未设置 |
| ErrMissingHeader | 100205 | 401 | Unset error message | 错误信息未设置 |
| ErrPasswordIncorrect | 100206 | 401 | Unset error message | 错误信息未设置 |
| ErrPermissionDenied | 100207 | 403 | Unset error message | 错误信息未设置 |
| ErrEncodingFailed | 100301 | 500 | Unset error message | 错误信息未设置 |
| ErrDecodingFailed | 100302 | 500 | Unset error message | 错误信息未设置 |
| ErrInvalidJSON | 100303 | 500 | Unset error message | 错误信息未设置 |
| ErrEncodingJSON | 100304 | 500 | Unset error message | 错误信息未设置 |
| ErrDecodingJSON | 100305 | 500 | Unset error message | 错误信息未设置 |
| ErrInvalidYaml | 100306 | 500 | Unset error message | 错误信息未设置 |
| ErrEncodingYaml | 100307 | 500 | Unset error message | 错误信息未设置 |
| ErrDecodingYaml | 100308 | 500 | Unset error message | 错误信息未设置 |
| ErrHTTPError | 100501 | 500 | Unset error message | 错误信息未设置 |
| ErrHTTPResponseDataParseError | 100502 | 500 | Unset error message | 错误信息未设置 |
| ErrHTTPClientGenerateError | 100503 | 500 | Unset error message | 错误信息未设置 |
| ErrGRPCClientGenerateError | 100701 | 500 | Unset error message | 错误信息未设置 |
| ErrGRPCClientCertificateError | 100702 | 500 | Unset error message | 错误信息未设置 |
| ErrGRPCClientDialError | 100703 | 500 | Unset error message | 错误信息未设置 |
| ErrGRPCClientInvokeServiceError | 100704 | 500 | Unset error message | 错误信息未设置 |
| ErrGRPCResponseDataParseError | 100705 | 500 | Unset error message | 错误信息未设置 |
| ErrUserNotFound | 110001 | 500 | Unset error message | 错误信息未设置 |

