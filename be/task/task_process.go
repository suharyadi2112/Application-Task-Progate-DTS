
package member 

import(
	"encoding/json"
	"net/http"
	outp "be/conn"
)

type Task struct {

    Id string `json:"id"`
    Task string `json:"task"`
    Assignee string `json:"assigne"`
    Date string `json:"date"`
    Status string `json:"status"`

}

func Gettask(w http.ResponseWriter, r *http.Request){

	
}
