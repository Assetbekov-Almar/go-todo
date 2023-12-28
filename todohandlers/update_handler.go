package todohandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"package/models"
	"package/utils"

	"github.com/gorilla/mux"
)

func (h *TodoHandler) UpdateHandler(w http.ResponseWriter, r * http.Request) {
	id := mux.Vars(r)["id"]

	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding todo: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	query := "UPDATE todo SET "
	params := []interface{}{}
	paramID := 1

	addParam := func(field string, value string) {
		if value != "" {
			query += fmt.Sprintf("%s = $%d, ", field, paramID)
			params = append(params, value)
			paramID++
		}
	}

	addParam("title", todo.Title)
	addParam("description", todo.Description)
	addParam("deadline", todo.Deadline)

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", paramID)
	params = append(params, id)

	_, err = h.DB.Exec(query, params...)
	if err != nil {
		errMsg := fmt.Sprintf( "Error updating todo with id %s: %v", id, err)
		utils.RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

}