package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Chains struct {
	Base string `json`
}

func main() {
	if os.Args[1] == "base" {

	}
	const url = "wss://rpc.ankr.com/base/ws/5b3e494193a39bb51d06eb3cb4735c305870850adcd0987f03bd6f88ea308737" // url string

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
