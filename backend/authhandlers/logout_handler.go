package authhandlers

import (
	"fmt"
	"net/http"
	"os"
	"package/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Access-Token")

	if accessToken == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "No access token")
		return
	}

	accessSecret := []byte(os.Getenv("ACCESS_SECRET"))

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return accessSecret, nil
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !token.Valid {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	var userID string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if id, ok := claims["userID"].(string); ok {
			userID = id
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "the id claim is not present")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusInternalServerError, "the token's claims cannot be cast to jwt.MapClaims")
		return
	}

	if userID == "" {
		utils.RespondWithError(w, http.StatusInternalServerError, "user ID is empty")
		return
	}

	_, err = h.DB.Exec("UPDATE users SET lastlogout = $1 WHERE id = $2", time.Now(), userID)
	if err != nil {
		errMsg := fmt.Sprintf("Error adding a new user: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	
	w.WriteHeader(http.StatusOK)
    w.Write([]byte("Logged out successfully"))
}