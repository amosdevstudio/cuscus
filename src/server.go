package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

const (
    PORT = ":8080"
    MSG_TEMPLATE = "<div class=\"message\"><p class=\"username\">%s</p><p class=\"msg-text\">%s</p></div>"
    MAX_SENT_MSGS = 64
)

var(
    msgCount int = 0
)

func loginUser(w http.ResponseWriter, r *http.Request){
    username := r.FormValue("username")
    pwd := r.FormValue("pwd")

    //Set the response type to html
    w.Header().Add("Content-type", "text/html; charset=UTF-8")

    if !userExists(username) {
        w.Write([]byte("NUH - UH!! User doesn't exist... yet... <a href=\"/signup-page\">Sign Up</a>."))
        return
    }

    if authUser(username, pwd){
        //Generate a new session id for the user
        sessionid := changeSessionid(username)
        //Send it to the user
        w.Header().Add("HX-Trigger", fmt.Sprintf("{\"setSessionId\":\"%s\"}", sessionid))
        w.Write([]byte("YUH - UHH; password right. :) <a href=\"/feed\">Login</a>"))
        return
    }

    //If the user exists but typed in the wrong password
    w.Write([]byte("NUH - UH!! Password is wrong. :("))
}

func signupUser(w http.ResponseWriter, r *http.Request){
    username := r.FormValue("username")
    pwd := r.FormValue("pwd")

    //Set the response type to html
    w.Header().Add("Content-type", "text/html; charset=UTF-8")

    if username == "" || pwd == ""{
        w.Write([]byte("Insert username and password!!"))
        return
    }

    if userExists(username){
        w.Write([]byte("NUH - UHH; User already exists. Maybe you want to <a href=\"/login-page\">Login</a>?"))
        return
    }

    //Generate a new session id for the user
    sessionid := addUser(username, pwd)
    //Send it to the user
    w.Header().Add("HX-Trigger", fmt.Sprintf("{\"setSessionId\":\"%s\"}", sessionid))
    w.Write([]byte("YIUPPIEE!! User created!!! <a href=\"/feed\">Login</a>"))
}

func addPost(w http.ResponseWriter, r *http.Request){
    username := r.FormValue("username")
    sessionid := r.FormValue("sessionid")
    content := r.FormValue("text")

    if content == "" {
        w.WriteHeader(http.StatusNotImplemented)
        return
    }

    if !authSession(username, sessionid){
        // The session isn't valid / expired
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    msgCount = addMessage(content, username)
}

func serveChat(w http.ResponseWriter, r *http.Request){
    w.Header().Add("Content-type", "text/html; charset=UTF-8")

    clientMsgCount, err := strconv.Atoi(r.FormValue("msgCount"))
    if err != nil {
        w.Write([]byte("Hmmm... There's something wrong... Try reloading the page..."))
        return
    }

    lastMessages := msgCount - clientMsgCount
    if lastMessages > MAX_SENT_MSGS{
        lastMessages = MAX_SENT_MSGS
    }

    var buffer bytes.Buffer
    for _, msg := range(*getLastMessages(lastMessages)){
        buffer.WriteString(fmt.Sprintf(MSG_TEMPLATE, html.EscapeString(msg.Username), html.EscapeString(msg.Message)))
    }

    w.Header().Add("HX-Trigger",  fmt.Sprintf("{\"setMsgCount\":\"%d\"}", msgCount))
    w.Write(buffer.Bytes())
}

func cleanup(){
    closeDB()
    // Cleanup stuff...
}

func handleInterrupts(){
    //This function handles interrupts gracefully.
    //It sends the interrupts to a channel
    interruptChannel := make(chan os.Signal, 1)
    signal.Notify(interruptChannel, os.Interrupt)
    go func() {
        //When an interrupt is received
        <-interruptChannel
        fmt.Println("Program killed!")
        //It calls a cleanup function
        cleanup()
        //And then exits the program
        os.Exit(1)
    }()
}

func main() {
    //Initialize the database
    initDB()

    //Function to handle interrupts gracefully
    //(For some reason defer in golang is kinda broken)
    handleInterrupts()

    feedFS := http.FileServer(http.Dir("pages/feed"))
    loginFS := http.FileServer(http.Dir("pages/login"))
    signupFS := http.FileServer(http.Dir("pages/signup"))

    http.Handle("/feed/", http.StripPrefix("/feed/", feedFS))
    http.Handle("/login-page/", http.StripPrefix("/login-page/", loginFS))
    http.Handle("/signup-page/", http.StripPrefix("/signup-page/", signupFS))
    http.Handle("/", loginFS)

    http.HandleFunc("/chat/", serveChat)
    http.HandleFunc("/post/", addPost)
    http.HandleFunc("/login/", loginUser)
    http.HandleFunc("/signup/", signupUser)

    msgCount = countMsgs()

    fmt.Printf("Server running on port %s\n", PORT)
    err := http.ListenAndServeTLS(PORT, "certs/localhost.crt", "certs/localhost.key", nil)
    if err != nil{
        fmt.Println(err)
    }
}
