package taskdata

import(
	"encoding/json"
	"net/http"
	outp "be_progate_task/connection"//koneksi

	"github.com/gorilla/mux"//routing mux
)

type Task struct {
    Id string `json:"id"`
    Task *string `json:"task"`
    Assignee *string `json:"assignee"`
    Deadline *string `json:"deadline"`
    Status *string `json:"status"`
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
        var err = row.Scan(&each.Id, &each.Task, &each.Assignee, &each.Status,&each.Deadline)

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

//update task status selesai atau belum
func ChangeStatusTask(w http.ResponseWriter, r *http.Request){

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

    CStatus := "1"

	result, err := db.Exec("UPDATE task SET status = $2 WHERE id = $1" , userID, CStatus)

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
