package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"harvest-mcp/harvestclient"

	"github.com/mark3labs/mcp-go/mcp"
)

// WhoAmIClient defines the interface for getting user information
type WhoAmIClient interface {
	GetWhoAmI(ctx context.Context) ([]byte, error)
}

func whoamiHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create Harvest client
	client := harvestclient.NewClient(
		os.Getenv("HARVEST_ACCESS_TOKEN"),
		os.Getenv("HARVEST_ACCOUNT_ID"),
		"MCP-Harvest-Integration (roy.touw@newstory.nl)",
	)

	return handleWhoAmI(ctx, client)
}

func handleWhoAmI(ctx context.Context, client WhoAmIClient) (*mcp.CallToolResult, error) {
	// Make the API call
	body, err := client.GetWhoAmI(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Pretty print the JSON response
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(prettyJSON.String()), nil
}
