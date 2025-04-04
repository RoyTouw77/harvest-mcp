package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"mcp-harvest",
		"1.0.0",
	)

	whoamiTool := mcp.NewTool(
		"whoami",
		mcp.WithDescription("Prompts Harvest API to return the current user"),
	)

	timeEntriesForClientTool := mcp.NewTool(
		"list_time_entries",
		mcp.WithDescription("Lists all time entries for date offset relative to today."),
		mcp.WithString("date_offset",
			mcp.Required(),
			mcp.Description("the date offset relative to today, e.g. -1 for tomorrow, 1 for yesterday, 0 for today"),
		),
	)

	s.AddTool(whoamiTool, whoamiHandler)
	s.AddTool(timeEntriesForClientTool, listTimeEntriesHandler)

	if err := server.ServeStdio(s); err != nil {
		// todo log
	}
}
