package main

import (
    `fmt`
    `github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic`
    `github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub`
    `github.com/bootdotdev/learn-pub-sub-starter/internal/routing`
    amqp `github.com/rabbitmq/amqp091-go`
    `log`
    `os`
    `os/signal`
    `syscall`
)

func main() {
    fmt.Println("Starting Peril client...")
    connectionString := "amqp://guest:guest@localhost:5672/"
    conn, err := amqp.Dial(connectionString)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    userName, err := gamelogic.ClientWelcome()
    if err != nil {
        log.Fatal(err)
    }
    _, _, err = pubsub.DeclareAndBind(
        conn,
        routing.ExchangePerilDirect,
        routing.PauseKey+"."+userName,
        routing.PauseKey,
        pubsub.QueueTypeTransient)
    if err != nil {
        log.Fatal(err)
    }
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
    <-sigs
    fmt.Println("Shutting down...")
}
