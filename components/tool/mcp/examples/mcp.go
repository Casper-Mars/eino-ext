/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"log"
	"time"

	"github.com/cloudwego/eino/components/tool"

	mcpp "github.com/Casper-Mars/eino-ext/components/tool/mcp"
)

func main() {
	startMCPServer()
	time.Sleep(1 * time.Second)
	ctx := context.Background()

	mcpTools := getMCPTool(ctx)

	for i, mcpTool := range mcpTools {
		fmt.Println(i, ":")
		info, err := mcpTool.Info(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Name:", info.Name)
		fmt.Println("Desc:", info.Desc)
		fmt.Println()
	}
}

func getMCPTool(ctx context.Context) []tool.BaseTool {
	// Create transport client (using SSE in this example)
	transportClient, err := transport.NewSSEClientTransport("http://127.0.0.1:8080/sse")
	if err != nil {
		log.Fatalf("Failed to create transport client: %v", err)
	}

	// Create MCP client using transport
	mcpClient, err := client.NewClient(transportClient, client.WithClientInfo(protocol.Implementation{
		Name:    "example MCP client",
		Version: "1.0.0",
	}))
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: mcpClient})
	if err != nil {
		log.Fatal(err)
	}

	return tools
}

type currentTimeReq struct {
	Timezone string `json:"timezone" description:"current time timezone"`
}

func startMCPServer() {
	// Create transport server (using SSE in this example)
	transportServer, err := transport.NewSSEServerTransport("127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to create transport server: %v", err)
	}

	// Create MCP server using transport
	mcpServer, err := server.NewServer(transportServer,
		// Set server implementation information
		server.WithServerInfo(protocol.Implementation{
			Name:    "Example MCP Server",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}
	tool, err := protocol.NewTool("current time", "Get current time with timezone, Asia/Shanghai is default", currentTimeReq{})
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
		return
	}

	// Register tool handler
	mcpServer.RegisterTool(tool, func(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		req := new(currentTimeReq)
		if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
			return nil, err
		}

		loc, err := time.LoadLocation(req.Timezone)
		if err != nil {
			return nil, fmt.Errorf("parse timezone with error: %v", err)
		}
		text := fmt.Sprintf(`current time is %s`, time.Now().In(loc))

		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: text,
				},
			},
		}, nil
	})

	go func() {
		defer func() {
			e := recover()
			if e != nil {
				fmt.Println(e)
			}
		}()
		if err = mcpServer.Run(); err != nil {
			log.Fatalf("Failed to start MCP server: %v", err)
			return
		}
	}()
}
