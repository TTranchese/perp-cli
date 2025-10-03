# Testing the Perplexity CLI

## Prerequisites
- Install Go from https://golang.org/dl/
- Ensure Go is in your PATH

## Build the CLI
```bash
cd c:/_WATest/perp-cli
go build -o perp.exe main.go
```

## Test Scenarios

### 1. Test with command-line auth flag
```bash
./perp.exe --auth=your_test_token_here
```
Expected: Should start with "Authentication token set via command line."

### 2. Test interactive mode without initial auth
```bash
./perp.exe
```
Then try these commands:
- `help` - Should show available commands
- `status` - Should show "Authentication: Not configured"
- `auth` - Should prompt for token input
- `auth test_token` - Should set token directly
- `status` - Should show masked token
- `exit` - Should quit

### 3. Test invalid commands
In interactive mode, try:
- `invalid_command` - Should show "Unknown command" message

## Expected Output Examples

### Starting with auth flag:
```
Authentication token set via command line.
Welcome to Perplexity CLI!
Type 'help' for available commands or 'exit' to quit.
perp>
```

### Interactive auth:
```
perp> auth
Enter your Perplexity API token: [you type token here]
Authentication token set successfully.
perp> status
Authentication: Configured (test...)
```