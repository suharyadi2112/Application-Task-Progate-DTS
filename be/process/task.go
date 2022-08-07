package taskdata 

import(
	"encoding/json"
    "fmt"
	"net/http"
	outp "be_progate_task/connection"

    "github.com/go-chi/chi/v5"

)

type Task struct {
    Id string `json:"id"`
    Task string `json:"task"`
    Assignee string `json:"assigne"`
    Date string `json:"date"`
    Status string `json:"status"`
}
type Response struct {
    Status string `json:"status"`
    Data []Task `json:"data"`//get data berdasarkan struct Task, untuk response
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
        var err = row.Scan(&each.Id, &each.Task, &each.Assignee, &each.Date, &each.Status)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        result = append(result, each)
    }
    
    defer db.Close()

    var res = Response{Status : "success", Data : result}
    
    json.NewEncoder(w).Encode(res)
	
}

//dapatkan data berdasarkan id
func Gettask_byid(w http.ResponseWriter, r *http.Request){

    userID := chi.URLParam(r, "userID")

    // fetch `"key"` from the request context
    ctx := r.Context()
    key := ctx.Value("key").(string)

    // respond to the client
    w.Write([]byte(fmt.Sprintf("hi %v, %v", userID, key)))

}