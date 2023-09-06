package main

import (
	"log"
	outp "be_progate_task/connection"//koneksiPostgres
    helper "be_progate_task/helper"//helper
    "github.com/tealeg/xlsx"
    "sync"
)

type Task struct {
    Email string `json:"email"`
    NamaTask string `json:"task"`
    Assignee string `json:"assignee"`
    Deadline string `json:"deadline"`
}

func main() {

	amqpmRbt, err := outp.Rabbitcon()
    defer amqpmRbt.Close()

	// koneksi db
	db, err := outp.Dbcon()
    if err != nil {
    	log.Fatalf(err.Error())
    }
    defer db.Close()

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

		var wg sync.WaitGroup

		for _, sheet := range xlFile.Sheets {
			for iRow, vRow := range sheet.Rows {
				if !helper.IsRowEmpty(vRow) {
					wg.Add(1)
					go func(sheet *xlsx.Sheet, iRow int, wg *sync.WaitGroup) {
						defer wg.Done()
						if iRow > 0 {
							for _, cell := range sheet.Rows[iRow].Cells {
								if cell.String() != "" {
									log.Printf("%s\t", cell.String())
								}
							}
						}
					}(sheet, iRow, &wg)
				}
			}
		}

		wg.Wait()
		log.Printf("Received: %s", msg.Body)
	}

}

