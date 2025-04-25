package main

import (
    "fmt"
    `github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub`
    `github.com/bootdotdev/learn-pub-sub-starter/internal/routing`
    amqp "github.com/rabbitmq/amqp091-go"
    `log`
    `os`
    `os/signal`
    `syscall`
)

func main() {
    fmt.Println("Starting Peril server...")
    connectionString := "amqp://guest:guest@localhost:5672/"
    connection, err := amqp.Dial(connectionString)
    if err != nil {
        log.Fatal(err)
    }
    defer connection.Close()
    fmt.Println("Connected to AMQP server")
    
    channel, err := connection.Channel()
    if err != nil {
        log.Fatal(err)
    }
    
    if err := pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}); err != nil {
        log.Fatal(err)
    }
    
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
    <-sigs
    fmt.Println("Signal received, shutting down...")
}
