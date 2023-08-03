
# GRPC 服务
GRPC服务同时支持tcp和unix两种协议，可以同时启用

## 添加服务
1. internal/pkg/proto目录


# 测试
## 手动测试
手动测试可以通过grpc客户端工具如`grpcurl`来测试. 

注意事项 **前提是grpc服务需要启动reflect服务。**
否则执行命令时会报`server does not support the reflection API`

参考自[服务的方法列表](https://chai2010.cn/advanced-go-programming-book/ch4-rpc/ch4-08-grpcurl.html)

###  gRPC服务启动反射服务
默认reflect服务是关闭的, 可以通过参数`--service.reflect`来开启。

`example-grpc --service.reflect=true`

如果正常启动，可以看到对应的启动日志
```
2023-07-25 07:02:05.986 INFO    grpcserver/config.go:59 gRPC service run with reflect service
``` 

### 查看gRPC服务
#### 查看gRPC安装的服务
查看gRPC服务器安装的服务`grpcurl -plaintext localhost:8081 list`
* `-plaintext`相当于curl的-k参数。即忽略gRPC服务tls证书检测。即使**gRPC**并没有启动TLS,
    也需要带上这个参数。否则会报`first record does not \
                   look like a TLS handshake`

```
root@wwhvw:~/go/src/eazycloud# grpcurl -plaintext localhost:8081 list
grpc.reflection.v1alpha.ServerReflection
version.VersionService
```

继续使用 list 子命令还可以查看 HelloService 服务的方法列表
```
root@wwhvw:~/go/src/eazycloud# grpcurl -plaintext localhost:8081 list version.VersionService
version.VersionService.Version
```

#### 查看gRPC方法的具体信息
如果还想了解方法的细节，可以使用 grpcurl 提供的 describe 子命令查看更详细的描述信息
```
version.VersionService is a service:
service VersionService {
  rpc Version ( .version.VersionRequest ) returns ( .version.VersionResponse );
}
```

也可以通过describe查看请求参数信息
```

root@wwhvw:~/go/src/eazycloud# grpcurl -plaintext localhost:8081 describe  .version.VersionResponse
version.VersionResponse is a message:
message VersionResponse {
  string GitVersion = 1;
  string GitCommit = 2;
  string GitTreeState = 3;
  string BuildDate = 4;
  string GoVersion = 5;
  string Compiler = 6;
  string Platform = 7;
}

root@wwhvw:~/go/src/eazycloud# grpcurl -plaintext localhost:8081 describe  .version.VersionRequest
version.VersionRequest is a message:
message VersionRequest {
}

```

### 调用方法进行测试
上面的命令可以看到version方法的参数和返回值。version不需要传参，因此
```
root@wwhvw:~/go/src/eazycloud# grpcurl -plaintext -d '{}' localhost:8081 version.VersionService/Version
{
  "GitVersion": "f9f71a2",
  "GitCommit": "f9f71a24d71abfb295b2912bee529fb95bd62299",
  "GitTreeState": "dirty",
  "BuildDate": "2023-07-25T07:29:42Z",
  "GoVersion": "go1.17.13",
  "Compiler": "gc",
  "Platform": "linux/amd64"
}
```

如果需要参数`$ grpcurl -plaintext -d '{"value":"gopher"}' 
          localhost:1234 HelloService.HelloService/Hello`


## TLS服务测试
