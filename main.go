package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func natsParameters(address string, channel string, name string) {
	// Check if the NATS channel is empty
	if strings.TrimSpace(channel) == "" {
		log.Fatal("Channel name cannot be empty")
	}

	// Connect to the NATS server with authentication
	nc, err := nats.Connect(address, nats.UserInfo("testuser", "testpassword"))
	if err != nil {
		log.Fatal(err)
	}
	// Close the connection when the function natsParameters() returns
	defer nc.Close()
	log.Printf("Connected to NATS server at %s", nc.ConnectedUrl())

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the stream exists
	streamName := "CHAT"
	_, err = js.StreamInfo(streamName)
	if err != nil {
		// Create the stream, and show messages from the last hour
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  []string{channel},
			Retention: nats.WorkQueuePolicy,
			MaxAge:    time.Hour,
		})
		if err != nil {
			log.Fatalf("Error creating stream: %v", err)
		}
	}

	// Log the channel name being used
	log.Printf("Subscribing to channel: %s", channel)

	// Subscribe to the channel, and show all messages
	sub, err := js.SubscribeSync(channel, nats.DeliverAll())
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}
	defer sub.Unsubscribe()

	// Create a context to handle the graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	//Goroutine to recieve messages
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := sub.NextMsg(30 * time.Second) // Increased timeout
				if err != nil {
					if err == nats.ErrTimeout {
						log.Println("No messages received within the timeout")
						continue
					}
					log.Printf("Error receiving message: %v", err)
					continue
				}
				log.Printf("Received message: %s", string(msg.Data))
			}
		}
	}(ctx)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down...")
		cancel() // Used to stop the goroutine that receives messages
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("Error unsubscribing: %v", err)
		}
		nc.Close()
		os.Exit(0)
	}()

	// Read the input from the user and publish messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Read user input
		text := scanner.Text()
		if strings.TrimSpace(text) == "" {
			continue
		}
		message := fmt.Sprintf("%s: %s", name, text)
		log.Printf("Publishing message: %s", message)
		if err := nc.Publish(channel, []byte(message)); err != nil {
			log.Printf("Error publishing message: %v", err)
		}
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <NATS_ADDRESS>, <NATS_CHANNEL>, <NATS_NAME>")
		return
	}

	address := os.Args[1]
	channel := os.Args[2]
	name := os.Args[3]

	natsParameters(address, channel, name)
}
