package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	if len(os.Args) < 2 {
		log.Fatal("No chain argument provided")
	}

	args := os.Args[1]

	apiKey := os.Getenv("ANKR_API_KEY")

	var url string
	if args == "base" {
		url = fmt.Sprintf("wss://rpc.ankr.com/%s_sepolia/ws/%s", "base", apiKey)
	} else {
		url = fmt.Sprintf("wss://rpc.ankr.com/%s_sepolia/ws/%s", "mantle", apiKey)
	}
	fmt.Println(url)
	client, err := ethclient.Dial(url)

	if err != nil {
		panic(err)
	}

	ch := make(chan *types.Header, 1024)
	sub, err := client.SubscribeNewHead(context.Background(), ch)

	if err != nil {
		panic(err)
	}

	fmt.Println("---subscribe-----")

	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("---unsubscribe-----")
		sub.Unsubscribe()
	}()

	go func() {
		for c := range ch {
			fmt.Println(c.Number)
		}
	}()

	<-sub.Err()

}
