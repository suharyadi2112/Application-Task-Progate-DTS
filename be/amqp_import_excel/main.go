package main

import (
	"log"
	outp "be_progate_task/connection"
    helper "be_progate_task/helper"
    "github.com/streadway/amqp"
    "github.com/joho/godotenv"
    "github.com/tealeg/xlsx"
    "github.com/schollz/progressbar/v3"
    "github.com/pusher/pusher-http-go/v5"
    "os"
    "time"
)

func main() {

	//ENV
    if err := godotenv.Load("../.env"); err != nil {
        log.Fatal("Fail Load .env")
    }

    amqpmRbt, err := outp.Rabbitcon()
    defer amqpmRbt.Close()
    handleErr(err, "Failed to connect to RabbitMQ")

    db, err := outp.Dbcon()
    defer db.Close()
    handleErr(err, "Failed to connect to the database")

    pusherClient := pusher.Client{
        AppID:   os.Getenv("APP_ID_PUSHER"),
        Key:     os.Getenv("KEY_PUSHER"),
        Secret:  os.Getenv("SECRET_PUSHER"),
        Cluster: os.Getenv("CLUSTER_PUSHER"),
        Secure:  true,
    }

    q, err := declareQueue(amqpmRbt, "ProDtsPostExcel")
    handleErr(err, "Failed to declare a queue")

    msgs, err := registerConsumer(amqpmRbt, q.Name)
    handleErr(err, "Failed to register a consumer")
   
	log.Println("Waiting for messages...")

	for msg := range msgs {
		uploadURL := msg.Body
		FixUrl := "../" + string(uploadURL)

		// Open Excel
		xlFile, err := xlsx.OpenFile(FixUrl)
		if err != nil {
			log.Fatalf("Fail to open file: %s", err)
		}
		
		Firstsheet := xlFile.Sheets[0]
		totalRows := helper.CountNonEmptyRows(Firstsheet)//hitung total row aktif
		bar := progressbar.Default(int64(totalRows))

		for _, sheet := range xlFile.Sheets {
			for iRow, vRow := range sheet.Rows {
				if !helper.IsRowEmpty(vRow) {
					if iRow > 0 {
						pusherClient.Trigger("my-channel", "my-event", map[string]interface{}{
				            "numberRow": iRow,
				            "totalRows": totalRows,
				        })
					
						dateValue, _ := sheet.Rows[iRow].Cells[2].GetTime(false)
						valDateValue := dateValue.Format("2006-01-02")

					    sqlStatement := `INSERT INTO task (task, assignee, deadline, status, email) VALUES ($1, $2, $3, $4, $5)`

					    _, err = db.Exec(sqlStatement, sheet.Rows[iRow].Cells[0].String(), sheet.Rows[iRow].Cells[1].String(), valDateValue, 0, sheet.Rows[iRow].Cells[3].String())

					    if err != nil {
					        log.Println(err.Error())
					        return
					    }
					}
					bar.Add(1)
				}
				time.Sleep(1 * time.Second)
			}
		}
		bar.Finish()
		log.Printf("Received: %s", msg.Body)
	}

}

func declareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
    return ch.QueueDeclare(
        queueName, // Nama antrian
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
}

func registerConsumer(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
    return ch.Consume(
        queueName, // Nama antrian
        "",        // consumer
        true,      // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
}

func handleErr(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %v", msg, err)
    }
}