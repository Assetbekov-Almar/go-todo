package authhandlers

import (
	"fmt"
	"net/http"
	"package/models"
	"package/utils"

	"github.com/gorilla/schema"
)

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := schema.NewDecoder()
	var user models.User
	err = decoder.Decode(&user, r.PostForm)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, hashedPassword)
	if err != nil {
		errMsg := fmt.Sprintf("Error adding a new user: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
}