package kafka

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ReadConfig() kafka.ConfigMap {
    // reads the client configuration from client.properties
    // and returns it as a key-value map
    m := make(map[string]kafka.ConfigValue)

    file, err := os.Open("client.properties")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if !strings.HasPrefix(line, "#") && len(line) != 0 {
            kv := strings.Split(line, "=")
            parameter := strings.TrimSpace(kv[0])
            value := strings.TrimSpace(kv[1])
            m[parameter] = value
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Failed to read file: %s", err)
        os.Exit(1)
    }

    return m
}

func Produce() {
	// creates a new producer instance
	conf := ReadConfig()
	p, _ := kafka.NewProducer(&conf)
	topic := "orders-topic"

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
			for e := range p.Events() {
					switch ev := e.(type) {
					case *kafka.Message:
							if ev.TopicPartition.Error != nil {
									fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
							} else {
									fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
											*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
							}
					}
			}
	}()

	// produces a sample message to the user-created topic
	p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte("key"),
			Value:          []byte("value"),
	}, nil)

	// send any outstanding or buffered messages to the Kafka broker and close the connection
	p.Flush(15 * 1000)
	p.Close()
}