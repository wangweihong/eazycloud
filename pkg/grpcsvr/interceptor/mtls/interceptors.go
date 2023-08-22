package mtls

import (
	"context"

	"github.com/wangweihong/eazycloud/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/skipper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// UnaryServerInterceptor returns a new unary server interceptor for mtls verify.
func UnaryServerInterceptor(skipperFunc ...skipper.SkipperFunc) grpc.UnaryServerInterceptor {
	name := "mtls"

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.F(ctx).Debugf("Interceptor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)

		if skipper.Skip(info.FullMethod, skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for %s", name, info.FullMethod)

			resp, err := handler(ctx, req)
			return resp, errors.UpdateStack(err)
		}

		peer, ok := peerFromContext(ctx)
		if !ok {
			log.F(ctx).Error("failed to get client peer information")
			return nil, errors.Wrap(code.ErrGRPCClientCertificateError, "failed to get client peer information")
		}

		if peer == nil || peer.AuthInfo == nil {
			log.F(ctx).Error("client is not authenticated")
			return nil, errors.Wrap(code.ErrGRPCClientCertificateError, "client is not authenticated")
		}

		// 获取客户端证书信息
		tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok {
			log.F(ctx).Error("failed to get TLSInfo from client AuthInfo")
			return nil, errors.Wrap(code.ErrGRPCClientCertificateError, "failed to get TLSInfo from client AuthInfo")
		}

		// 获取客户端证书
		certificates := tlsInfo.State.PeerCertificates
		if len(certificates) == 0 {
			log.F(ctx).Error("client certificate is missing")
			return nil, errors.Wrap(code.ErrGRPCClientCertificateError, "client certificate is missing")
		}

		// 验证客户端证书的主体信息
		clientCert := certificates[0]
		if clientCert.Subject.CommonName != "client.example.com" {
			log.F(ctx).Error("invalid client certificate subject")
			return nil, errors.Wrap(code.ErrGRPCClientCertificateError, "invalid client certificate subject")
		}

		resp, err := handler(ctx, req)
		return resp, errors.UpdateStack(err)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for trace.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		peer, ok := peerFromContext(stream.Context())
		if !ok {
			log.Error("failed to get client peer information")
			return errors.Wrap(code.ErrGRPCClientCertificateError, "failed to get client peer information")
		}

		if peer == nil || peer.AuthInfo == nil {
			log.Error("client is not authenticated")
			return errors.Wrap(code.ErrGRPCClientCertificateError, "client is not authenticated")
		}

		// 获取客户端证书信息
		tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok {
			log.Error("failed to get TLSInfo from client AuthInfo")
			return errors.Wrap(code.ErrGRPCClientCertificateError, "failed to get TLSInfo from client AuthInfo")
		}
		// 获取客户端证书
		certificates := tlsInfo.State.PeerCertificates
		if len(certificates) == 0 {
			log.Error("client certificate is missing")
			return errors.Wrap(code.ErrGRPCClientCertificateError, "client certificate is missing")
		}

		// 验证客户端证书的主体信息
		clientCert := certificates[0]
		if clientCert.Subject.CommonName != "client.example.com" {
			log.Error("invalid client certificate subject")
			return errors.Wrap(code.ErrGRPCClientCertificateError, "invalid client certificate subject")
		}

		return handler(srv, stream)
	}
}

// 从 gRPC context 中获取客户端信息.
func peerFromContext(ctx context.Context) (*peer.Peer, bool) {
	p, ok := peer.FromContext(ctx)
	return p, ok
}
