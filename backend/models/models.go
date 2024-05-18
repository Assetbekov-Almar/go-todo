package models

type Todo struct {
    ID          int    `json:"id" schema:"id"`
    Title       string `json:"title" schema:"title"`
    Description string `json:"description" schema:"description"`
    Deadline    string `json:"deadline" schema:"deadline"`
}

type User struct {
    ID          int    `json:"id" schema:"id"`
    Username    string `json:"username" schema:"username"`
    Password    string `json:"password" schema:"password"`
    LastLogout  string `json:"lastlogout" schema:"lastlogout"`
}

type MovieResult struct {
    Adult            bool     `json:"adult"`
    BackdropPath     string   `json:"backdrop_path"`
    GenreIds         []int    `json:"genre_ids"`
    ID               int      `json:"id"`
    OriginalLanguage string   `json:"original_language"`
    OriginalTitle    string   `json:"original_title"`
    Overview         string   `json:"overview"`
    Popularity       float64  `json:"popularity"`
    PosterPath       string   `json:"poster_path"`
    ReleaseDate      string   `json:"release_date"`
    Title            string   `json:"title"`
    Video            bool     `json:"video"`
    VoteAverage      float64  `json:"vote_average"`
    VoteCount        int      `json:"vote_count"`
}

type Movie struct {
    Page       int  `json:"page" schema:"page"`
    Results []MovieResult `json:"results"`    
}