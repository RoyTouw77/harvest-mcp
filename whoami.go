package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
)

func whoamiHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.harvestapp.com/api/v2/users/me.json", nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create request: %v", err)), nil
	}

	// Add required headers
	req.Header.Add("Authorization", "Bearer "+os.Getenv("HARVEST_ACCESS_TOKEN"))
	req.Header.Add("Harvest-Account-ID", os.Getenv("HARVEST_ACCOUNT_ID"))
	req.Header.Add("User-Agent", "MCP-Harvest-Integration (roy.touw@newstory.nl)")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to make request: %v", err)), nil
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read response: %v", err)), nil
	}

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(body))), nil
	}

	// Pretty print the JSON response
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(prettyJSON.String()), nil
}
