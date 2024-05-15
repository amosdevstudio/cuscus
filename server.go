package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strconv"
)

const (
    PORT = ":8080"
    FEED_PATH = "/feed/"
    LOGIN_PATH = "/login/"

    MSG_TEMPLATE = "<div class=\"message\"><p class=\"username\">%s</p><p class=\"msg-text\">%s</p></div>"
    TMP_TEMPLATE = "<div id=\"tmp-lastMsg\">%d</div>"
)

var (
    msgCount = 0
)

func getMsgs(lastMsgStr string) []byte{
    var buffer bytes.Buffer
    var lastmsgid int
    if lastMsgStr == "" {lastmsgid = 0} else {
        var err error
        lastmsgid, err = strconv.Atoi(lastMsgStr)
        if err != nil {fmt.Println(err); return []byte("")}
    }

    msgs := getLastMessages(msgCount - lastmsgid)
    for _, msg := range *msgs{
        buffer.WriteString(fmt.Sprintf(MSG_TEMPLATE, html.EscapeString(msg.Username), html.EscapeString(msg.Message)))
    }
    buffer.WriteString(fmt.Sprintf(TMP_TEMPLATE, msgCount))

    return buffer.Bytes()
}

func handleChat (w http.ResponseWriter, r *http.Request){
    w.Write(getMsgs(r.FormValue("lastmsgid")))
}

func handlePost(w http.ResponseWriter, r *http.Request){
    username := r.FormValue("username")
    if authUser(username, r.FormValue("pwd")){
        msgCount = addMessage(r.FormValue("text"), username)
    } else {
        w.WriteHeader(http.StatusUnauthorized)
    }
}

func handleLogin (w http.ResponseWriter, r *http.Request){
    username := r.FormValue("username")
    pwd := r.FormValue("pwd")
    exists := userExists(username)
    if username == "" || pwd == "" {
        w.Write([]byte("NUH - UH!! Insert username and password."))
    } else if exists && authUser(username, pwd){
        w.Write([]byte("YUH - UHH; password right. :) <a href=\"/feed\">Login</a>"))
    } else if exists {
        w.Write([]byte("NUH - UH!! Password is wrong. :("))
    } else {
        addUser(username, pwd)
        w.Write([]byte("YIPPIEE!!:3 User created. <a href=\"/feed\">Login</a>"))
    }

}


func main() {
    initDB()
    defer closeDB()


    feedFS := http.FileServer(http.Dir("./static/feed"))
    loginFS := http.FileServer(http.Dir("./static/login"))


    http.Handle(FEED_PATH, http.StripPrefix(FEED_PATH, feedFS))
    http.Handle(LOGIN_PATH, http.StripPrefix(LOGIN_PATH, loginFS))
    http.Handle("/", loginFS)

    http.HandleFunc("/chat/", handleChat)
    http.HandleFunc("/post/", handlePost)
    http.HandleFunc("/login-user/", handleLogin)

    fmt.Printf("Server running on port %s\n", PORT)
    err := http.ListenAndServe(PORT, nil)
    if err != nil{
        fmt.Println(err)
    }
}
