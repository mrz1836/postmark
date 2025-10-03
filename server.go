package postmark

import (
	"context"
)

// GetCurrentServer gets details for the server associated
// with the currently in-use server API Key
func (client *Client) GetCurrentServer(ctx context.Context) (Server, error) {
	res := Server{}
	err := client.get(ctx, "server", &res)
	return res, err
}

// EditCurrentServer updates details for the server associated
// with the currently in-use server API Key
func (client *Client) EditCurrentServer(ctx context.Context, server Server) (Server, error) {
	res := Server{}
	err := client.put(ctx, "server", server, &res)
	return res, err
}