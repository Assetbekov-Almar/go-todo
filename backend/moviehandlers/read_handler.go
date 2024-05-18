package moviehandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"package/models"
	"package/utils"
)

// 3ecd9c3e778aeccc6a7d58d8bb572daf
var MOVIE_DB_SECRET = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzZWNkOWMzZTc3OGFlY2NjNmE3ZDU4ZDhiYjU3MmRhZiIsInN1YiI6IjY2NDkwYWI2OTU5OTEyYjNkOTdhNDNiOCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.aiPUSVlhqzvHnAlDO3LLSOWP2DtmCIWIiBh6FJZ2WS8"

type MovieHandler struct {
	DB *sql.DB
}

func authInTheMovieDB(w http.ResponseWriter) {
	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", MOVIE_DB_SECRET))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var movies models.Movie	
    err := json.Unmarshal(body, &movies)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, movies)
}

func (h *MovieHandler) ReadHandler(w http.ResponseWriter, r *http.Request) {
	authInTheMovieDB(w)
}