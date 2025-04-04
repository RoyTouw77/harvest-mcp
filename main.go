package main

import (
	"os"

	"harvest-mcp/harvestclient"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	client := harvestclient.NewClient(
		os.Getenv("HARVEST_ACCESS_TOKEN"),
		os.Getenv("HARVEST_ACCOUNT_ID"),
		"MCP-Harvest-Integration (roy.touw@newstory.nl)",
	)

	s := server.NewMCPServer(
		"mcp-harvest",
		"1.0.0",
	)

	s.AddTool(whoamiTool, whoamiHandler(client))
	s.AddTool(timeEntriesForClientTool, listTimeEntriesHandler(client))

	if err := server.ServeStdio(s); err != nil {
		// todo log
	}
}
