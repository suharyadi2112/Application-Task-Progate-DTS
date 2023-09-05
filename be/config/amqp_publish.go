package rabbitmq

import (
    "github.com/streadway/amqp"
)

func DeclareAndPublishMessage(ch *amqp.Channel, queueName, body string) error {
    q, err := ch.QueueDeclare(
        queueName, // Nama antrian
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return err
    }

    err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
        })
    return err
}
