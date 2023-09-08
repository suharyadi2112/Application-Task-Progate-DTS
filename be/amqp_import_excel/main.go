package main

import (
	"log"
	outp "be_progate_task/connection"//koneksiPostgres
    helper "be_progate_task/helper"//helper
    "github.com/joho/godotenv"
    "github.com/tealeg/xlsx"
    "github.com/schollz/progressbar/v3"
    "github.com/pusher/pusher-http-go/v5"
    "os"
    "time"
    // "fmt"
)

func main() {

	//ENV
    err := godotenv.Load("../.env")
    if err != nil {
        log.Fatal("Fail Load .env")
    }

	amqpmRbt, err := outp.Rabbitcon()
    defer amqpmRbt.Close()

	// koneksi db
	db, err := outp.Dbcon()
    if err != nil {
    	log.Fatalf(err.Error())
    }
    defer db.Close()

    pusherClient := pusher.Client{
		AppID: os.Getenv("APP_ID_PUSHER"),
		Key: os.Getenv("KEY_PUSHER"),
		Secret: os.Getenv("SECRET_PUSHER"),
		Cluster: os.Getenv("CLUSTER_PUSHER"),
		Secure: true,
	}

	q, err := amqpmRbt.QueueDeclare(
		"ProDtsPostExcel", // Nama antrian
		true,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := amqpmRbt.Consume(
		q.Name, // Nama antrian
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}
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

