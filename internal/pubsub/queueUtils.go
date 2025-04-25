package pubsub

import (
    `context`
    `encoding/json`
    `fmt`
    amqp `github.com/rabbitmq/amqp091-go`
    `log`
)

type QueueType int

const (
    QueueTypeTransient QueueType = iota
    QueueTypePersistent
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
    data, err := json.Marshal(val)
    if err != nil {
        fmt.Errorf("Error marshalling JSON: %v", err)
        return err
    }
    if err := ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
        ContentType: "application/json",
        Body:        data,
    }); err != nil {
        fmt.Errorf("Error publishing JSON: %v", err)
        return err
    }
    return nil
}

func DeclareAndBind(
        conn *amqp.Connection,
        exchange,
        queueName,
        key string,
        simpleQueueType QueueType,
) (*amqp.Channel, amqp.Queue, error) {
    
    channel, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
        return nil, amqp.Queue{}, err
    }
    
    exclusive, autoDelete, durable := false, false, true
    if simpleQueueType == QueueTypeTransient {
        exclusive, autoDelete, durable = true, true, false
    }
    queue, err := channel.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
    if err != nil {
        log.Fatal(err)
        return nil, amqp.Queue{}, err
    }
    
    if err := channel.QueueBind(queue.Name, key, exchange, false, nil); err != nil {
        log.Fatal(err)
        return nil, amqp.Queue{}, err
    }
    return channel, queue, nil
}
