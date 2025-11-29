package grpcinternal

import (
	"context"
	"fmt"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func GetRequestInfo(log app.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)

		peerInfo, _ := peer.FromContext(ctx)
		ip := peerInfo.Addr.String()

		md, _ := metadata.FromIncomingContext(ctx)
		userAgent := md.Get("user-agent")[0]

		reqStatus, _ := status.FromError(err)
		statusCode := reqStatus.Code()

		log.Info(fmt.Sprintf("%s  %s  %d  %s  %v",
			ip,
			info.FullMethod,
			statusCode,
			userAgent,
			time.Since(start)))

		return resp, err
	}
}
