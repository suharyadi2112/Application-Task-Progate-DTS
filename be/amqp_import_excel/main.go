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
)

type Task struct {
    Email string `json:"email"`
    NamaTask string `json:"task"`
    Assignee string `json:"assignee"`
    Deadline string `json:"deadline"`
}

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
					func(sheet *xlsx.Sheet, iRow int) {
						if iRow > 0 {
							pusherClient.Trigger("my-channel", "my-event", map[string]interface{}{
					            "percentage": iRow,
					            "totalRows": totalRows,
					        })
							for _, cell := range sheet.Rows[iRow].Cells {
								if cell.String() != "" {
									// log.Printf("%s\t", cell.String())
								}
							}
						}
						bar.Add(1)
					}(sheet, iRow)
				}
			}
		}
		bar.Finish()
		log.Printf("Received: %s", msg.Body)
	}

}

