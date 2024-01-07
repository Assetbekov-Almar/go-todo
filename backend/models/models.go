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
}