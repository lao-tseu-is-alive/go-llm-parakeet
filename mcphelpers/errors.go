package mcphelpers

import "fmt"

/*
type MCPClientError struct {
	Message string
}

func (e *MCPClientError) Error() string {
	return fmt.Sprintf("MCPClientError: %s", e.Message)
}
*/

type MCPClientCreationError struct {
	Message string
}

func (e *MCPClientCreationError) Error() string {
	return fmt.Sprintf("MCPClientCreationError: %s", e.Message)
}

type MCPClientInitializationError struct {
	Message string
}

func (e *MCPClientInitializationError) Error() string {
	return fmt.Sprintf("MCPClientInitializationError: %s", e.Message)
}

type MCPGetToolsError struct {
	Message string
}

func (e *MCPGetToolsError) Error() string {
	return fmt.Sprintf("MCPGetToolsError: %s", e.Message)
}

type MCPToolCallError struct {
	Message string
}

func (e *MCPToolCallError) Error() string {
	return fmt.Sprintf("MCPToolCallError: %s", e.Message)
}

type MCPResultExtractionError struct {
	Message string
}

func (e *MCPResultExtractionError) Error() string {
	return fmt.Sprintf("MCPResultExtractionError: %s", e.Message)
}
