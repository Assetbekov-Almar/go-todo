package todohandlers

import (
	"fmt"
	"net/http"
	"package/utils"

	"github.com/gorilla/mux"
)

func (h *TodoHandler) DeleteHandler(w http.ResponseWriter, r * http.Request) {
	id := mux.Vars(r)["id"]
	result, err := h.DB.Exec("DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		errMsg := fmt.Sprintf("Error deleting todo with id %s: %v", id, err)
		utils.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errMsg := fmt.Sprintf("Error getting rows affected: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	if rowsAffected == 0 {
		errMsg := fmt.Sprintf("No todo found with id %s", id)
		utils.RespondWithError(w, http.StatusNotFound, errMsg)
		return
	}
}