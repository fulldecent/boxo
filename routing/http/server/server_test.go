package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ipfs/boxo/routing/http/types"
	"github.com/ipfs/boxo/routing/http/types/iter"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	router := &mockContentRouter{}
	server := httptest.NewServer(Handler(router))
	t.Cleanup(server.Close)
	serverAddr := "http://" + server.Listener.Addr().String()

	results := iter.FromSlice([]iter.Result[types.ProviderResponse]{
		{Val: &types.ReadBitswapProviderRecord{
			Protocol: "transport-bitswap",
			Schema:   types.SchemaBitswap,
		}}},
	)

	c := "baeabep4vu3ceru7nerjjbk37sxb7wmftteve4hcosmyolsbsiubw2vr6pqzj6mw7kv6tbn6nqkkldnklbjgm5tzbi4hkpkled4xlcr7xz4bq"
	cb, err := cid.Decode(c)
	require.NoError(t, err)

	router.On("FindProviders", mock.Anything, cb).
		Return(results, nil)

	resp, err := http.Get(serverAddr + ProvidePath + c)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	header := resp.Header.Get("Content-Type")
	require.Equal(t, mediaTypeJSON, header)

	resp, err = http.Get(serverAddr + ProvidePath + "BAD_CID")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 400, resp.StatusCode)
	header = resp.Header.Get("Content-Type")
	require.Equal(t, "text/plain; charset=utf-8", header)
}

type mockContentRouter struct{ mock.Mock }

func (m *mockContentRouter) FindProviders(ctx context.Context, key cid.Cid) (iter.ResultIter[types.ProviderResponse], error) {
	args := m.Called(ctx, key)
	return args.Get(0).(iter.ResultIter[types.ProviderResponse]), args.Error(1)
}
func (m *mockContentRouter) ProvideBitswap(ctx context.Context, req *BitswapWriteProvideRequest) (time.Duration, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *mockContentRouter) Provide(ctx context.Context, req *WriteProvideRequest) (types.ProviderResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(types.ProviderResponse), args.Error(1)
}
