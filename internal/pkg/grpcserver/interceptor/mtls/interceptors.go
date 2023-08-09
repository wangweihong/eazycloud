package mtls

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// UnaryServerInterceptor returns a new unary server interceptor for mtls verify.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		peer, ok := peerFromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("failed to get client peer information")
		}

		if peer == nil || peer.AuthInfo == nil {
			return nil, fmt.Errorf("client is not authenticated")
		}

		// 获取客户端证书信息
		tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok {
			return nil, fmt.Errorf("failed to get TLSInfo from client AuthInfo")
		}

		// 获取客户端证书
		certificates := tlsInfo.State.PeerCertificates
		if len(certificates) == 0 {
			return nil, fmt.Errorf("client certificate is missing")
		}

		// 验证客户端证书的主体信息
		clientCert := certificates[0]
		if clientCert.Subject.CommonName != "client.example.com" {
			return nil, fmt.Errorf("invalid client certificate subject")
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for trace.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		peer, ok := peerFromContext(stream.Context())
		if !ok {
			return fmt.Errorf("failed to get client peer information")
		}

		if peer == nil || peer.AuthInfo == nil {
			return fmt.Errorf("client is not authenticated")
		}

		// 获取客户端证书信息
		tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok {
			return fmt.Errorf("failed to get TLSInfo from client AuthInfo")
		}
		// 获取客户端证书
		certificates := tlsInfo.State.PeerCertificates
		if len(certificates) == 0 {
			return fmt.Errorf("client certificate is missing")
		}

		// 验证客户端证书的主体信息
		clientCert := certificates[0]
		if clientCert.Subject.CommonName != "client.example.com" {
			return fmt.Errorf("invalid client certificate subject")
		}

		return handler(srv, stream)
	}
}

// 从 gRPC context 中获取客户端信息.
func peerFromContext(ctx context.Context) (*peer.Peer, bool) {
	p, ok := peer.FromContext(ctx)
	return p, ok
}
