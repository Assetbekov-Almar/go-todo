package authhandlers

import (
	"fmt"
	"net/http"
	"os"
	"package/utils"
	"time"

	"github.com/golang-jwt/jwt"
)


func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Refresh-Token")

	if refreshToken == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "No refresh token")
		return
	}

	refreshSecret := []byte(os.Getenv("REFRESH_SECRET"))

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return refreshSecret, nil
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

	row := h.DB.QueryRow("SELECT lastlogout FROM users WHERE id = $1", userID)

	var lastLogout time.Time
    err = row.Scan(&lastLogout)
    if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

	accessToken, refreshToken, errMsg := generateJWT(userID, lastLogout)

	if errMsg != "" {
		utils.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	w.Header().Add("Access-Token", accessToken)
	w.Header().Add("Refresh-Token", refreshToken)
}