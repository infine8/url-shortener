package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"url-shortener/internal/config"

	ssov1 "github.com/infine8/go-sso-proto/gen/go/sso"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
    api ssov1.AuthClient
    log *slog.Logger
}

func New(
    ctx context.Context,
    log *slog.Logger,
    grpcConfig *config.GrpcClient,
) (*Client, error) {
    const op = "grpc.New"

    retryOpts := []grpcretry.CallOption{
        grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
        grpcretry.WithMax(uint(grpcConfig.Retries)),
        grpcretry.WithPerRetryTimeout(grpcConfig.Timeout),
    }

	logOpts := []grpclog.Option{
        grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
    }

	cc, err := grpc.DialContext(ctx, grpcConfig.Address,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithChainUnaryInterceptor(
            grpclog.UnaryClientInterceptor(interceptorLogger(log), logOpts...),
            grpcretry.UnaryClientInterceptor(retryOpts...),
        ))

    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	grpcClient := ssov1.NewAuthClient(cc)

    return &Client{
        api: grpcClient,
    }, nil
}


func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
    const op = "grpc.IsAdmin"

    resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
        UserId: userID,
    })
	
    if err != nil {
        return false, fmt.Errorf("%s: %w", op, err)
    }

    return resp.IsAdmin, nil
}

