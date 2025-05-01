package main

import (
  "fmt"
  `github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic`
  `github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub`
  `github.com/bootdotdev/learn-pub-sub-starter/internal/routing`
  amqp "github.com/rabbitmq/amqp091-go"
  `log`
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
  defer channel.Close()
  
  if _, _, err := pubsub.DeclareAndBind(
    connection,
    routing.ExchangePerilTopic,
    routing.GameLogSlug,
    fmt.Sprintf("%s.*", routing.GameLogSlug),
    pubsub.QueueTypePersistent); err != nil {
    log.Fatal(err)
  }
  
  gamelogic.PrintServerHelp()
replLoop:
  for {
    input := gamelogic.GetInput()
    if len(input) == 0 {
      continue
    }
    switch input[0] {
    case "pause":
      fmt.Println("Pausing...")
      if err := pubsub.PublishJSON(
        channel,
        routing.ExchangePerilDirect,
        routing.PauseKey,
        routing.PlayingState{IsPaused: true},
      ); err != nil {
        log.Fatal(err)
      }
      break
    case "resume":
      fmt.Println("Resuming...")
      if err := pubsub.PublishJSON(
        channel,
        routing.ExchangePerilDirect,
        routing.PauseKey,
        routing.PlayingState{IsPaused: false},
      ); err != nil {
        log.Fatal(err)
      }
      break
    case "quit":
      fmt.Println("Quitting...")
      break replLoop
    default:
      fmt.Println("Unrecognized command:", input[0])
    }
  }
}
