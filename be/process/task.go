package taskdata 

import(
	"encoding/json"
	"net/http"
	outp "be_progate_task/connection"

    "github.com/go-chi/chi/v5"//routing
    "io/ioutil"//post

		"fmt"
)

type Task struct {
    Id string `json:"id"`
    Task string `json:"task", validate:"required"`
    Assignee string `json:"assignee"`
    Deadline string `json:"deadline"`
    Status string `json:"status"`
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

    row, err := db.Query("Select * from task")//query dari mysql
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    defer db.Close()

    var result []Task //variable menampung data Task dari DB

    for row.Next(){//looping data
        var each = Task{}//Task dari struct
        var err = row.Scan(&each.Id, &each.Task, &each.Assignee, &each.Deadline, &each.Status)

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

    userID := chi.URLParam(r, "userID")

    db, err := outp.Dbcon()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    SqlQuery := `SELECT * FROM task WHERE id=$1;`

    var task Task

    row := db.QueryRow(SqlQuery, userID)

    err = row.Scan(&task.Id, &task.Task, &task.Assignee, &task.Deadline, &task.Status)

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

    w.Header().Set("Content-Type", "application/json")
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

    taskStruct := Task{}//gunakan jika masukan bukan array

    body, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(body, &taskStruct)

    task := taskStruct.Task
    assignee := taskStruct.Assignee
    deadline := taskStruct.Deadline
    status := taskStruct.Status

		fmt.Println(assignee, deadline,task)

    sqlStatement := `INSERT INTO task (task, assignee, deadline, status) VALUES ($1, $2, $3, $4)`

    _, err = db.Exec(sqlStatement, task,assignee,deadline,status)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
		defer db.Close()

    var res = ResponseSingle{Status : "success", Data : taskStruct}

    json.NewEncoder(w).Encode(res)

}

//delete data berdasarkan id
func DelTask_id(w http.ResponseWriter, r *http.Request){

    // semua origin mendapat ijin akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // semua method diperbolehkan masuk
    w.Header().Set("Access-Control-Allow-Methods", "*")
    // semua header diperbolehkan untuk disisipkan
    w.Header().Set("Access-Control-Allow-Headers", "*")

    userID := chi.URLParam(r, "userID")

		db, err := outp.Dbcon()//koneksi
		if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}

		defer db.Close()
		_, err = db.Exec("DELETE FROM task WHERE id=$1", userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()

		json.NewEncoder(w).Encode("success")

}
