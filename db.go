package main

import (
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

type User struct {
    Name  string `db:"username"`
    PasswordHash string `db:"password_hash"`
    UserId int `db:"userid"`
}

type Message struct {
    Message string `db:"message"`
    Username string `db:"username"`
    MessageId string `db:"messageid"`
}

var db *sqlx.DB

func initDB() {
    var err error
    db, err = sqlx.Connect("postgres", "user=postgres dbname=cuscus sslmode=disable password=fottutapassword host=localhost")
    if err != nil {
        log.Fatalln(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    } else {
        log.Println("Successfuly connected")
    }
}

func countMsgs() int {
    count := 0
    db.Get(&count, "SELECT currval('messageid') FROM messages;")
    return count
}

func userExists(username string) bool {
    exists := 0
    err := db.Get(&exists, "SELECT 1 FROM users WHERE username = $1", username)
    if err != nil {
        log.Println(err)
    }
    return (exists == 1)
}

func addUser(username string, pwd string){
    _, err := db.Queryx("INSERT INTO users (username, password_hash) VALUES ($1, crypt($2, gen_salt('md5')));", username, pwd)
    if err != nil{
        log.Println(err)
    }
}

func addMessage(message string, username string) int {
    count := 0
    err := db.Get(&count, "INSERT INTO messages (message, username) VALUES ($1, $2) RETURNING messageid;", message, username)
    if err != nil{
        log.Println(err)
    }
    return count
}

func getLastMessages(lastN int) *[]Message {
    messages := []Message{}
    err := db.Select(&messages, "SELECT * FROM MESSAGES ORDER BY messageid DESC LIMIT $1;", lastN)
    if err != nil{
        log.Println(err)
    }
    return &messages
}

func authUser(username string, pwd string) bool{
    result := false
    err := db.Get(&result, `SELECT (password_hash = crypt($1, password_hash)) AS password_match FROM users WHERE username = $2;`, pwd, username)
    if err != nil{
        log.Println(err)
    }
    return result
}

func closeDB (){
    db.Close()
}
