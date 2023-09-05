package connect

import(
    "github.com/streadway/amqp"
)

func Rabbitcon() (*amqp.Channel, error) {
    connMrabbit, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return nil, err
    }

    chMrabbit, err := connMrabbit.Channel()
    if err != nil {
        return nil, err
    }

    return chMrabbit, nil
}
