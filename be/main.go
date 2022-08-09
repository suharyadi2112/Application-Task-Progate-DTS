package main

import(
	"net/http"
	"fmt"

	b "be_progate_task/process"//mengimport file lain untuk digunakan functionnya
	"github.com/gorilla/mux"//mux route
)

func main(){

	k := mux.NewRouter()

    // IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
    k.HandleFunc("/task_changestatus/{userID}", b.ChangeStatusTask).Methods("GET")
    k.HandleFunc("/task", b.Gettask).Methods("GET")
    k.HandleFunc("/task_id/{userID}", b.Gettask_byid).Methods("GET")
    k.HandleFunc("/task_post", b.PostTask).Methods("POST")
    k.HandleFunc("/task_del/{userID}", b.DelTask_id).Methods("GET")
    k.HandleFunc("/task_up/{userID}", b.UpTask_id).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
    	
    k.Use(loggingMiddleware)

	fmt.Println("server started at localhost:9999")
	http.ListenAndServe(":9999", k)//web service running
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        fmt.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}
