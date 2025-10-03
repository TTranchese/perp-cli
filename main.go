package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sgaunet/perplexity-go/v2"
)

type CLI struct {
	authToken string
	client    *perplexity.Client
}

func main() {
	// Parse command line flags
	authFlag := flag.String("auth", "", "Perplexity API authentication token")
	flag.Parse()

	cli := &CLI{}

	// Check if auth token was provided via flag
	if *authFlag != "" {
		cli.authToken = *authFlag
		cli.client = perplexity.NewClient(*authFlag)
		fmt.Println("Authentication token set via command line.")
	}

	// Start interactive mode
	fmt.Println("Welcome to Perplexity CLI!")
	fmt.Println("Type 'help' for available commands or 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("perp> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return
		case "help":
			cli.showHelp()
		case "auth":
			cli.handleAuth(parts)
		case "status":
			cli.showStatus()
		case "ask":
			cli.handleAsk(parts)
		default:
			// If not a known command, treat it as a question
			cli.handleQuestion(input)
		}
	}
}

func (c *CLI) showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  auth <token>  - Set the Perplexity API authentication token")
	fmt.Println("  status        - Show current authentication status")
	fmt.Println("  ask <query>   - Ask a question to Perplexity AI")
	fmt.Println("  help          - Show this help message")
	fmt.Println("  exit/quit     - Exit the CLI")
	fmt.Println()
	fmt.Println("You can also type any question directly without using 'ask'")
}

func (c *CLI) handleAuth(parts []string) {
	var token string

	if len(parts) < 2 {
		fmt.Print("Enter your Perplexity API token: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			token = strings.TrimSpace(scanner.Text())
			if token == "" {
				fmt.Println("No token provided.")
				return
			}
		}
	} else {
		token = parts[1]
	}

	// Initialize the Perplexity client
	c.authToken = token
	c.client = perplexity.NewClient(token)
	fmt.Println("Authentication token set successfully.")
}

func (c *CLI) showStatus() {
	if c.authToken != "" {
		// Show only first few characters for security
		masked := c.authToken
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		fmt.Printf("Authentication: Configured (%s)\n", masked)
	} else {
		fmt.Println("Authentication: Not configured")
	}
}

func (c *CLI) handleAsk(parts []string) {
	if len(parts) < 2 {
		fmt.Println("Please provide a question. Usage: ask <your question>")
		return
	}

	question := strings.Join(parts[1:], " ")
	c.askPerplexity(question)
}

func (c *CLI) handleQuestion(question string) {
	c.askPerplexity(question)
}

func (c *CLI) askPerplexity(question string) {
	if c.client == nil {
		fmt.Println("Please set your authentication token first using the 'auth' command.")
		return
	}

	fmt.Println("Asking Perplexity AI...")

	// Create the request using the perplexity-go package
	messages := []perplexity.Message{
		{
			Role:    "user",
			Content: question,
		},
	}

	// Make the API call
	ctx := context.Background()
	resp, err := c.client.CreateChatCompletion(ctx, "llama-3.1-sonar-small-128k-online", messages)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}

	// Display the response
	if resp != nil && len(resp.Choices) > 0 {
		fmt.Println("\n--- Response ---")
		fmt.Println(resp.Choices[0].Message.Content)
		if resp.Usage != nil {
			fmt.Printf("\n--- Usage: %d tokens ---\n", resp.Usage.TotalTokens)
		}
	} else {
		fmt.Println("No response received from Perplexity AI")
	}
}
