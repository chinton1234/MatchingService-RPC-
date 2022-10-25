package main

import (
        "context"
        "log"
        "math/rand"
        "time"
        amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
        if err != nil {
                log.Panicf("%s: %s", msg, err)
        }
}

func randomString(l int) string {
        bytes := make([]byte, l)
        for i := 0; i < l; i++ {
                bytes[i] = byte(randInt(65, 90))
        }
        return string(bytes)
}

func randInt(min int, max int) int {
        return min + rand.Intn(max-min)
}

func sendRPC(n string) (res string, err error) {
        conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        q, err := ch.QueueDeclare(
                "",    // name
                false, // durable
                false, // delete when unused
                true,  // exclusive
                false, // noWait
                nil,   // arguments
        )
        failOnError(err, "Failed to declare a queue")

        msgs, err := ch.Consume(
                q.Name, // queue
                "",     // consumer
                true,   // auto-ack
                false,  // exclusive
                false,  // no-local
                false,  // no-wait
                nil,    // args
        )
        failOnError(err, "Failed to register a consumer")

        corrId := randomString(32)
        log.Printf(" [.] corrId %s", corrId)

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        err = ch.PublishWithContext(ctx,
                "",          // exchange
                "rpc_queue", // routing key
                false,       // mandatory
                false,       // immediate
                amqp.Publishing{
                        ContentType:   "text/plain",
                        CorrelationId: corrId,
                        ReplyTo:       q.Name,
                        Body:          []byte(n),
                })

        failOnError(err, "Failed to publish a message")

        for d := range msgs {
                log.Printf(" [.] Got %s", d.CorrelationId)
                if corrId == d.CorrelationId {
                        res,err = string(d.Body),nil
                        break
                }
        }

        return
}

func main() {
        rand.Seed(time.Now().UTC().UnixNano())

        n := "create 63513f36506c462416419f61"

        log.Printf("[x] Requesting matching(%s)", n)
        res, err := sendRPC(n)
        
        failOnError(err, "Failed to handle RPC request")
        log.Printf(res)
}