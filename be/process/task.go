package taskdata

import(
	"encoding/json"
	"net/http"
	outp "be_progate_task/connection"//koneksiPostgres
    amqpFunc "be_progate_task/config"//koneksi
    helper "be_progate_task/helper"//helper
	"github.com/gorilla/mux"//routing mux
    "log"
    "io/ioutil"
    "time"
)

type Task struct {
    Id string `json:"id"`
    Task *string `json:"task"`
    Assignee *string `json:"assignee"`
    Deadline *string `json:"deadline"`
    Status *string `json:"status"`
    Email *string `json:"email"`
}
type ResponseArr struct {
    Status string `json:"status"`
    Data []Task `json:"data"`
}
type ResponseSingle struct{
    Status string `json:"status"`
    Data Task `json:"data"`
}

//dapatkan semua data dari task
func Gettask(w http.ResponseWriter, r *http.Request){

    // semua origin mendapat ijin akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // semua method diperbolehkan masuk
    w.Header().Set("Access-Control-Allow-Methods", "*")
    // semua header diperbolehkan untuk disisipkan
    w.Header().Set("Access-Control-Allow-Headers", "*")

    db, err := outp.Dbcon()//koneksi dari file lain(package lain)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    defer db.Close()

    row, err := db.Query("Select * from task")//query dari sql
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    defer db.Close()

    var result []Task //variable menampung data Task dari DB

    for row.Next(){//looping data
        var each = Task{}//Task dari struct
        var err = row.Scan(&each.Id, &each.Task, &each.Assignee, &each.Deadline,&each.Status, &each.Email)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        result = append(result, each)
    }

    defer db.Close()

    var res = ResponseArr{Status : "success", Data : result}

    json.NewEncoder(w).Encode(res)

}

//dapatkan data berdasarkan id
func Gettask_byid(w http.ResponseWriter, r *http.Request){

    // semua origin mendapat ijin akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // semua method diperbolehkan masuk
    w.Header().Set("Access-Control-Allow-Methods", "*")
    // semua header diperbolehkan untuk disisipkan
    w.Header().Set("Access-Control-Allow-Headers", "*")

    vars := mux.Vars(r)
    userID := vars["userID"]

    db, err := outp.Dbcon()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    SqlQuery := `SELECT * FROM task WHERE id=$1;`

    var task Task

    row := db.QueryRow(SqlQuery, userID)

    err = row.Scan(&task.Id, &task.Task, &task.Assignee, &task.Status,&task.Deadline)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var res = ResponseSingle{Status : "success", Data : task}

    json.NewEncoder(w).Encode(res)
}

//tambah task
func PostTask(w http.ResponseWriter, r *http.Request){

    // semua origin mendapat ijin akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // semua method diperbolehkan masuk
    w.Header().Set("Access-Control-Allow-Methods", "*")
    // semua header diperbolehkan untuk disisipkan
    w.Header().Set("Access-Control-Allow-Headers", "*")

	
    db, err := outp.Dbcon()//koneksi
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    task := r.FormValue("task")
    assignee := r.FormValue("assignee")
    deadline := r.FormValue("deadline")
    status := "0"
  	
    sqlStatement := `INSERT INTO task (task, assignee, deadline, status) VALUES ($1, $2, $3, $4)`

    _, err = db.Exec(sqlStatement, task,assignee,deadline,status)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	defer db.Close()
	
	json.NewEncoder(w).Encode("success")

}

//delete data berdasarkan id
func DelTask_id(w http.ResponseWriter, r *http.Request){

    // semua origin mendapat ijin akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // semua method diperbolehkan masuk
    w.Header().Set("Access-Control-Allow-Methods", "*")
    // semua header diperbolehkan untuk disisipkan
    w.Header().Set("Access-Control-Allow-Headers", "*")

    vars := mux.Vars(r)
    userID := vars["userID"]

	db, err := outp.Dbcon()//koneksi
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	defer db.Close()
	result, err := db.Exec("DELETE FROM task WHERE id=$1", userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	rows, err := result.RowsAffected()

	if rows != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}else{
		json.NewEncoder(w).Encode("success")
	}

}
//update task diluar status
func UpTask_id(w http.ResponseWriter, r *http.Request){

    // semua origin mendapat ijin akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // semua method diperbolehkan masuk
    w.Header().Set("Access-Control-Allow-Methods", "*")
    // semua header diperbolehkan untuk disisipkan
    w.Header().Set("Access-Control-Allow-Headers", "*")

	db, err := outp.Dbcon()//koneksi
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	defer db.Close()

	vars := mux.Vars(r)
    userID := vars["userID"]//get param mux

	task := r.FormValue("task")
    assignee := r.FormValue("assignee")
    deadline := r.FormValue("deadline")

	result, err := db.Exec("UPDATE task SET task = $2, assignee = $3, deadline = $4 WHERE id = $1" , userID,task, assignee,deadline)

	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	rows, err := result.RowsAffected()
	defer db.Close()
	if rows != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}else{
		json.NewEncoder(w).Encode("success")
	}

}

func ChangeStatusTask(w http.ResponseWriter, r *http.Request) {
    // Setelah mengatasi CORS, pastikan semua origin, metode, dan header diperbolehkan
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "*")
    w.Header().Set("Access-Control-Allow-Headers", "*")

    // Mengambil parameter dari URL
    vars := mux.Vars(r)
    userID := vars["userID"]

    // Koneksi database
    db, err := outp.Dbcon()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    CStatus := "1"

    // Update status tugas di database
    result, err := db.Exec("UPDATE task SET status = $2 WHERE id = $1", userID, CStatus)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rows, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rows != 1 {
        http.Error(w, "Tidak dapat mengubah status tugas", http.StatusInternalServerError)
        return
    }

    // ---------------- AMQP -----------------//
    amqpmRbt, err := outp.Rabbitcon()
    if err != nil {
        log.Printf("Connection Failed mrabbit")
        return
    }
    defer amqpmRbt.Close()

    // Kirim pesan ke RabbitMQ untuk mengirim email
    body := userID
    queueName := "ProDtsSendMail"
    err = amqpFunc.DeclareAndPublishMessage(amqpmRbt, queueName, body)
    if err != nil {
        log.Println("Gagal mendeklarasikan dan mengirim pesan: %v", err)
    }
    // ---------------- AMQP -----------------//

    // Kirim respons JSON
    json.NewEncoder(w).Encode("success")
}


func PostExcel(w http.ResponseWriter, r *http.Request) {

    // Setelah mengatasi CORS, pastikan semua origin, metode, dan header diperbolehkan
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "*")
    w.Header().Set("Access-Control-Allow-Headers", "*")

    file, header, err := r.FormFile("files")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    if !helper.IsValidExcelFile(header.Filename){
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    TimeName := time.Now().UTC()
    uploadURL := "static_file/excel_import_task/" + helper.TrimDateFileName(TimeName) + "_" + helper.GenerateRandomString(3) + "_" + header.Filename

    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = ioutil.WriteFile(uploadURL, fileBytes, 0644)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }


    // ---------------- AMQP -----------------//
    amqpmRbt, err := outp.Rabbitcon()
    if err != nil {
        log.Printf("Connection Failed mrabbit")
        return
    }
    defer amqpmRbt.Close()

    // Kirim pesan ke RabbitMQ untuk memproses excel
    queueName := "ProDtsPostExcel"
    body := uploadURL
    err = amqpFunc.DeclareAndPublishMessage(amqpmRbt, queueName, body)
    if err != nil {
        log.Println("Gagal mendeklarasikan dan mengirim pesan: %v", err)
    }
    // ---------------- AMQP -----------------//


    res := map[string]interface{}{
        "message": "Berhasil",
        "data": "",
        "status": 200,
    }

    json.NewEncoder(w).Encode(res)
}


