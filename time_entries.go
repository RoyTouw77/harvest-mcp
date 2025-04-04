package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"harvest-mcp/harvestclient"

	"github.com/mark3labs/mcp-go/mcp"
)

// TimeEntriesClient defines the interface for getting time entries
type TimeEntriesClient interface {
	GetTimeEntries(ctx context.Context, date string) ([]byte, error)
}

func listTimeEntriesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dateOffset, ok := request.Params.Arguments["date_offset"].(string)
	if !ok {
		return mcp.NewToolResultError("date offset must be an integer"), nil
	}
	dateOffsetInt, err := strconv.Atoi(dateOffset)
	if err != nil {
		return mcp.NewToolResultError("date offset must be an integer"), nil
	}

	// Adjust the date based on the offset
	targetDate := time.Now().AddDate(0, 0, -dateOffsetInt).Format("20060102")

	// Create Harvest client
	client := harvestclient.NewClient(
		os.Getenv("HARVEST_ACCESS_TOKEN"),
		os.Getenv("HARVEST_ACCOUNT_ID"),
		"MCP-Harvest-Integration (roy.touw@newstory.nl)",
	)

	return handleTimeEntries(ctx, client, targetDate)
}

func handleTimeEntries(ctx context.Context, client TimeEntriesClient, targetDate string) (*mcp.CallToolResult, error) {
	// Make the API call
	body, err := client.GetTimeEntries(ctx, targetDate)
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
