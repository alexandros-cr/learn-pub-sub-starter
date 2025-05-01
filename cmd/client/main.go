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
  gameState := gamelogic.NewGameState(userName)
  repl(gameState)
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
  <-sigs
  fmt.Println("Shutting down...")
}

func contains(slice []string, val string) bool {
  for _, item := range slice {
    if item == val {
      return true
    }
  }
  return false
}

func repl(gameState *gamelogic.GameState) {
  locations := []string{"americas", "europe", "africa", "asia", "antarctica", "australia"}
  untitTypes := []string{"infantry", "cavalry", "artillery"}

replLoop:
  for {
    input := gamelogic.GetInput()
    if len(input) == 0 {
      continue
    }
    switch input[0] {
    case "spawn":
      // validate location
      if len(input) == 1 {
        fmt.Println("Please specify a location to spawn")
        break
      }
      location := input[1]
      if !contains(locations, location) {
        fmt.Println("please enter a valid location")
        break
      }
      // validate unit type
      if len(input) == 2 {
        fmt.Println("Please specify a unit type to spawn")
        break
      }
      unitType := input[2]
      if !contains(untitTypes, unitType) {
        fmt.Println("please enter a valid unit type")
        break
      }
      // spawn unit
      if err := gameState.CommandSpawn(input); err != nil {
        fmt.Printf("Failed to spawn unit: %v\n", err)
      }
      break
    case "move":
      // validate location
      if len(input) == 1 {
        fmt.Println("Please specify a location to move")
        break
      }
      location := input[1]
      if !contains(locations, location) {
        fmt.Println("please enter a valid location")
        break
      }
      // validate id
      if len(input) == 2 {
        fmt.Println("Please specify a unit to move")
        break
      }
      // move unit
      if _, err := gameState.CommandMove(input); err != nil {
        fmt.Println("Failed to move unit: %v", err)
      } else {
        fmt.Println("Move unit successful")
      }
      break
    case "status":
      gameState.CommandStatus()
      break
    case "help":
      gamelogic.PrintClientHelp()
      break
    case "spam":
      fmt.Println("Spamming not allowed yet")
      break
    case "quit":
      gamelogic.PrintQuit()
      break replLoop
    default:
      fmt.Println("Please enter a valid command")
      break
    }
  }
}
