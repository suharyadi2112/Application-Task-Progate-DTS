package main 

import(
	"net/http"

	b "be_progate_task/process"//mengimport file lain untuk digunakan functionnya

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main(){
	
	r := chi.NewRouter()//Chi Route
	r.Route("/progatedts", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/task", b.Gettask)//all task
		r.Get("/task_id/{userID}", b.Gettask_byid)//task by id
	})
	
	http.ListenAndServe(":8080", r)//web service running
}
