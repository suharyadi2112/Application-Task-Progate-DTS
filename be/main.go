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
		r.Post("/task_post", b.PostTask)//post task
		r.Delete("/task_del/{userID}", b.DelTask_id)//post task
		r.Put("/task_up/{userID}", b.UpTask_id)//post task

	})

	http.ListenAndServe(":9999", r)//web service running
}
