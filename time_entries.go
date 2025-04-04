package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

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

	// Create HTTP client and request
	client := &http.Client{}
	url := fmt.Sprintf("https://api.harvestapp.com/v2/time_entries?from=%s&to=%s", targetDate, targetDate)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
