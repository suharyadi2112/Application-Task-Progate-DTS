package main

import(
    "fmt"
    "net/http"
    "html/template"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
 )

//home
func home(w http.ResponseWriter, r *http.Request){
    var tmpl = template.Must(template.New("index").ParseFiles("index.html"))
    if err := tmpl.Execute(w, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
//edit
func edit(w http.ResponseWriter, r *http.Request){

    Taskid := chi.URLParam(r, "task_id")
    var tmpl = template.Must(template.New("edit").ParseFiles("edit.html"))
    var data_id = map[string]interface{}{
        "id_task": Taskid,
    }
    if err := tmpl.Execute(w, data_id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {

    //chi routing
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    //start routing
    r.Get("/", home)//home
    r.Get("/home", home)//home
    r.Get("/edit/{task_id}", edit)//home

    fmt.Println("server started at localhost:9000")
    http.ListenAndServe(":9000", r)
}

