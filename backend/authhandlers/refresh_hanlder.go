package authhandlers

import (
	"fmt"
	"net/http"
	"os"
	"package/utils"

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

	accessToken, refreshToken, errMsg := generateJWT()

	if errMsg != "" {
		utils.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	w.Header().Add("Access-Token", accessToken)
	w.Header().Add("Refresh-Token", refreshToken)
}