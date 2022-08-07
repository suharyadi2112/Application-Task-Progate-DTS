package main 

import(
	"net/http"

	b "be_progate_task/process"//mengimport file lain untuk digunakan functionnya

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main(){
	
	r := chi.NewRouter()//Chi Route
	r.Use(middleware.Logger)

	r.Get("/task", b.Gettask)// get all task
	
	http.ListenAndServe(":8080", r)//web service running
}
