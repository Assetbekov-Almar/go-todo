package todohandlers

import (
	"database/sql"
	"net/http"
	"package/models"
	"package/utils"
)

type TodoHandler struct {
    DB *sql.DB
}

func (h *TodoHandler) ReadHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT * FROM todo;")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Deadline); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, todos)
}