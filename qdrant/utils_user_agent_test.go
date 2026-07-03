package qdrant

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
)

func TestGrpcClientDialOptions_UserAgent(t *testing.T) {
	old := providerVersion
	t.Cleanup(func() { providerVersion = old })
	providerVersion = "1.2.3"

	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	var captured metadata.MD
	srv := grpc.NewServer(grpc.UnaryInterceptor(func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		captured, _ = metadata.FromIncomingContext(ctx)
		return handler(ctx, req)
	}))
	qcCluster.RegisterClusterServiceServer(srv, &qcCluster.UnimplementedClusterServiceServer{})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(func() { srv.Stop() })

	dialOpts := []grpc.DialOption{
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent(providerUserAgent()),
	}
	conn, err := grpc.NewClient("passthrough:///bufnet", dialOpts...)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client := qcCluster.NewClusterServiceClient(conn)
	_, _ = client.ListClusters(context.Background(), &qcCluster.ListClustersRequest{
		AccountId: "00000000-0000-0000-0000-000000000001",
	})

	userAgent := captured.Get("user-agent")
	require.NotEmpty(t, userAgent)
	assert.Contains(t, userAgent[0], "terraform-provider-qdrant-cloud/1.2.3")
}
