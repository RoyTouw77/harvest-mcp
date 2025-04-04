package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

type WhoAmIClient interface {
	GetWhoAmI(ctx context.Context) ([]byte, error)
}

var whoamiTool = mcp.NewTool(
	"whoami",
	mcp.WithDescription("Prompts Harvest API to return the current user"),
)

func whoamiHandler(client WhoAmIClient) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		body, err := client.GetWhoAmI(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(prettyJSON.String()), nil
	}
}
