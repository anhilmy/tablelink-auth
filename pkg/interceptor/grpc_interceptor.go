package interceptor

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RequestHeader struct {
	AccessToken string
}

const AccessTokenContextKey = "request_header"

func IncomingRequest() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Fatal("cannot get metadata")
		}

		linkService := md.Get("x-link-service")
		if len(linkService) == 0 {
			return nil, errors.New("invalid header: no x-link-service header")
		} else if linkService[0] != "be" {
			return nil, errors.New("invalid header: invalid x-link-service value")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, errors.New("invalid header: no authorization token")
		}

		logData := RequestHeader{
			AccessToken: tokens[0],
		}

		ctx = context.WithValue(ctx, AccessTokenContextKey, logData)

		// Handle the request
		res, err := handler(ctx, req)
		return res, err
	}
}
