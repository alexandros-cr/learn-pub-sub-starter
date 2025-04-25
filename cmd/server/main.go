package main

import (
    "fmt"
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
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
    <-sigs
    fmt.Println("Signal received, shutting down...")
}
