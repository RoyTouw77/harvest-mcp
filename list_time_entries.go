package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

type TimeEntriesClient interface {
	GetTimeEntries(ctx context.Context, date string) ([]byte, error)
}

var timeEntriesForClientTool = mcp.NewTool(
	"list_time_entries",
	mcp.WithDescription("Lists all time entries for date offset relative to today."),
	mcp.WithString("date_offset",
		mcp.Required(),
		mcp.Description("the date offset relative to today, e.g. -1 for tomorrow, 1 for yesterday, 0 for today"),
	),
)

func listTimeEntriesHandler(client TimeEntriesClient) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		return handleTimeEntries(ctx, client, targetDate)
	}
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
