package todohandlers

import (
	"fmt"
	"net/http"
	"package/models"
	"package/utils"

	"github.com/gorilla/schema"
)

func (h *TodoHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := schema.NewDecoder()
	var todo models.Todo
	err = decoder.Decode(&todo, r.PostForm)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = h.DB.Exec("INSERT INTO todo (title, description, deadline) VALUES ($1, $2, $3)", todo.Title, todo.Description, todo.Deadline)
	if err != nil {
		errMsg := fmt.Sprintf("Error inserting new todo: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, errMsg)
		fmt.Fprintf(w, "Error inserting new todo: %v", err)
		return
	}

}