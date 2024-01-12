package authhandlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"package/models"
	"package/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
    DB *sql.DB
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    return string(bytes), nil
}

func createToken(expMinutes int, mySigningKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	expDuration := time.Duration(expMinutes) * time.Minute
	claims["exp"] = time.Now().Add(expDuration).Unix()

	tokenString, err := token.SignedString([]byte(mySigningKey))

	if err != nil {
		log.Fatalf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func generateJWT() (string, string, string) {
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	accessToken, accessErr := createToken(15, accessSecret)

	if accessErr != nil {
		return "", "", "Error creating access token"
	}

	refreshToken, refreshErr := createToken(60, refreshSecret)

	if refreshErr != nil {
		return "", "", "Error creating access token"
	}
	
	return accessToken, refreshToken, ""
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	row := h.DB.QueryRow("SELECT password FROM users WHERE username = $1", user.Username)

	var hashedPassword string
    err = row.Scan(&hashedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
        } else {
            utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
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