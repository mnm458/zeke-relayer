package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("RESERVOIR_API_KEY")
	if apiKey == "" {
		log.Fatal("RESERVOIR_API_KEY environment variable is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Base chain
	url := fmt.Sprintf("wss://ws-base-sepolia.reservoir.tools?api_key=%s", apiKey)
	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", apiKey)},
		},
	})
	if err != nil {
		log.Fatal("Error connecting to Reservoir: ", err)
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	fmt.Println("Connected to Reservoir")

	var message map[string]interface{}
	for {
		err := wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}

		fmt.Println("Message received: ", message)

		// When the connection is ready, subscribe to the top-bids event
		if message["status"] == "ready" {
			fmt.Println("Subscribing")
			err = wsjson.Write(ctx, conn, map[string]interface{}{
				"type":  "subscribe",
				"event": "top-bid.changed",
			})
			if err != nil {
				log.Println("Error subscribing: ", err)
				break
			}

			// To unsubscribe, send the following message
			// err = wsjson.Write(ctx, conn, map[string]interface{}{
			//     "type":    "unsubscribe",
			//     "event": "top-bid.changed",
			// })
			// if err != nil {
			//     log.Println("Error unsubscribing: ", err)
			//     break
			// }
		}
	}
}
