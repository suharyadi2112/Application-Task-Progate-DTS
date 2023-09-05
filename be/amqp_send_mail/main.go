package main

import (
	"log"
	"io/ioutil"
	outp "be_progate_task/connection"//koneksiPostgres
    "github.com/joho/godotenv"
    "net/smtp"
    "strings"
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
    // Sender data.
    from := os.Getenv("MAIL_FROM")
    password := os.Getenv("PASSWORD")
    // smtp server configuration.
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")

    // template email
    templateFile := "../template_email/mark_as_done.html" 
    amqpmRbt, err := outp.Rabbitcon()

    defer amqpmRbt.Close()

	q, err := amqpmRbt.QueueDeclare(
		"ProDtsSendMail", // Nama antrian
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
		// koneksi db
		db, err := outp.Dbcon()
	    if err != nil {
	    	log.Println(err)
	        return
	    }
	    //dapatkan email
		SqlQuery := `SELECT email,task, assignee, deadline FROM task WHERE id=$1;`
	    var task Task
	    rowGet := db.QueryRow(SqlQuery, msg.Body)
	    err = rowGet.Scan(&task.Email, &task.NamaTask, &task.Assignee, &task.Deadline)

	    if err != nil {
	    	log.Println(err)
	        return
	    }

	    defer db.Close()

		log.Printf("Message process : %s", task.Email)

	    // Receiver email address.
	    to := task.Email
	    header := make(map[string]string)
		header["Subject"] = "Email Task Done"
		header["From"] = from
		header["To"] = to
		header["Content-Type"] = "text/html; charset=UTF-8"

		message := ""
		for key, value := range header {
		    message += key + ": " + value + "\r\n"
		}


		templateBytes, err := ioutil.ReadFile(templateFile)
		if err != nil {
		    log.Fatalf("Failed to read template file: %v", err)
		}

		message += "\r\n" + string(templateBytes)


		message = strings.Replace(message, "[nama-task]", task.NamaTask, 1)
		message = strings.Replace(message, "[assignee]", task.Assignee, 1)
		message = strings.Replace(message, "[deadline]", task.Deadline, 1)

		// log.Println(message)
	    // Authentication.
	    auth := smtp.PlainAuth("", from, password, smtpHost)

	    // Sending email.
	    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	    if err != nil {
	        log.Println(err)
	        return
	    }
	    log.Printf("Email Sent Successfully: %s", task.Email)
		log.Printf("Received: %s", msg.Body)
	}
}
