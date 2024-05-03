package main

import (
	kafka "github.com/mnm458/zeke-relayer/pkg/services/kafka"
	"context"
	"fmt"
	"log"
	"os"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
)

type Transaction struct {

}

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

	//Ramp contract Address
	contractAddress := common.HexToAddress("0x3782CB11C629eFA26Fa9B3DAB9C2d4f4eCFCBAc4")
    query := ethereum.FilterQuery{
        Addresses: []common.Address{contractAddress},
    }

    logs := make(chan types.Log)
    sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
    if err != nil {
        log.Fatal(err)
    }

    for {
        select {
        case err := <-sub.Err():
            log.Fatal(err)
        case vLog := <-logs:
            fmt.Println(vLog.TxHash) 
			kafka.Produce(vLog.TxHash.Hex(), args)
        }
    }
}

