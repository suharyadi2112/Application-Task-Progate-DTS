package main

import(
	"net/http"
	"fmt"
    "github.com/joho/godotenv"
	b "be_progate_task/process"//mengimport file lain untuk digunakan functionnya
	"github.com/gorilla/mux"//mux route
    "log"
)

func main(){

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Fail Load .env")
    }
	k := mux.NewRouter()

    // IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
    k.HandleFunc("/task_changestatus/{userID}", b.ChangeStatusTask).Methods("GET")
    k.HandleFunc("/task", b.Gettask).Methods("GET")
    k.HandleFunc("/task_id/{userID}", b.Gettask_byid).Methods("GET")
    k.HandleFunc("/task_post", b.PostTask).Methods("POST")
    k.HandleFunc("/task_del/{userID}", b.DelTask_id).Methods("GET")
    k.HandleFunc("/task_up/{userID}", b.UpTask_id).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
    k.HandleFunc("/task_excel", b.PostExcel).Methods("POST")
    	
    k.Use(loggingMiddleware)

	fmt.Println("server started at localhost:9009")
	http.ListenAndServe("127.0.0.1:9009", k)//web service running
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        fmt.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}
