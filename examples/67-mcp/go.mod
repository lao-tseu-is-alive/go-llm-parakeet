module 67-mcp

go 1.23.1

require (
	github.com/joho/godotenv v1.5.1
	github.com/parakeet-nest/parakeet v0.0.0-00010101000000-000000000000
)

require github.com/mark3labs/mcp-go v0.8.3 // indirect

replace github.com/parakeet-nest/parakeet => ../..
