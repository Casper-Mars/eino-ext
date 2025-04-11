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

package mcp

import (
	"context"
	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"testing"

	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestTool(t *testing.T) {
	// Create transport client (using SSE in this example)
	transportClient, err := transport.NewSSEClientTransport("http://127.0.0.1:8080/sse")
	if err != nil {
		t.Fatalf("Failed to create transport client: %v", err)
	}

	// Create MCP client using transport
	mcpClient, err := client.NewClient(transportClient, client.WithClientInfo(protocol.Implementation{
		Name:    "example MCP client",
		Version: "1.0.0",
	}))
	if err != nil {
		t.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()
	ctx := context.Background()

	tools, err := GetTools(ctx, &Config{Cli: mcpClient, ToolNameList: []string{"name"}})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tools))
	info, err := tools[0].Info(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "name", info.Name)

	result, err := tools[0].(tool.InvokableTool).InvokableRun(ctx, "{\"input\": \"123\"}")
	assert.NoError(t, err)
	assert.Equal(t, "{\"content\":[{\"type\":\"text\",\"text\":\"hello\"}]}", result)
}

type mockMCPClient struct{}

func (m *mockMCPClient) Initialize(ctx context.Context, request mcp.InitializeRequest) (*mcp.InitializeResult, error) {
	panic("implement me")
}

func (m *mockMCPClient) Ping(ctx context.Context) error {
	panic("implement me")
}

func (m *mockMCPClient) ListResources(ctx context.Context, request mcp.ListResourcesRequest) (*mcp.ListResourcesResult, error) {
	panic("implement me")
}

func (m *mockMCPClient) ListResourceTemplates(ctx context.Context, request mcp.ListResourceTemplatesRequest) (*mcp.ListResourceTemplatesResult, error) {
	panic("implement me")
}

func (m *mockMCPClient) ReadResource(ctx context.Context, request mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	panic("implement me")
}

func (m *mockMCPClient) Subscribe(ctx context.Context, request mcp.SubscribeRequest) error {
	panic("implement me")
}

func (m *mockMCPClient) Unsubscribe(ctx context.Context, request mcp.UnsubscribeRequest) error {
	panic("implement me")
}

func (m *mockMCPClient) ListPrompts(ctx context.Context, request mcp.ListPromptsRequest) (*mcp.ListPromptsResult, error) {
	panic("implement me")
}

func (m *mockMCPClient) GetPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	panic("implement me")
}

func (m *mockMCPClient) ListTools(ctx context.Context, request mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
	return &mcp.ListToolsResult{
		Tools: []mcp.Tool{
			{
				Name:        "name",
				Description: "description",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"input": map[string]interface{}{"type": "string"},
					},
					Required: []string{"input"},
				},
			},
			{
				Name:        "name2",
				Description: "description",
			},
		},
	}, nil
}

func (m *mockMCPClient) CallTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "hello",
			},
		},
		IsError: false,
	}, nil
}

func (m *mockMCPClient) SetLevel(ctx context.Context, request mcp.SetLevelRequest) error {
	panic("implement me")
}

func (m *mockMCPClient) Complete(ctx context.Context, request mcp.CompleteRequest) (*mcp.CompleteResult, error) {
	panic("implement me")
}

func (m *mockMCPClient) Close() error {
	panic("implement me")
}

func (m *mockMCPClient) OnNotification(handler func(notification mcp.JSONRPCNotification)) {
	panic("implement me")
}
