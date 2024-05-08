package main

import (
	"fmt"
	"net/http"
	"strings"
)
const PORT = ":8080"
const BUFFER_SIZE = 100000
const MSG_TEMPLATE = "<div class=\"message\">%s</div>"
var messages []interface{} = make([]interface{}, BUFFER_SIZE)
var head int = 0
var tail int = 0

func getMsgSlices(start int) [][]interface{}{
    result := make([][]interface{}, 2)
    if start > tail {
        result[0] = messages[start:]
        result[1] = messages[:tail+1]
        return result
    }
    result[0] = messages[start:tail+1]
    result[1] = make([]interface{}, 0)
    return result
}

func addMsg(msg string) {
    tail = (tail+1) % BUFFER_SIZE
    if tail == head{ head = (head + 1) % BUFFER_SIZE }
    messages[tail] = msg
}

func parseMessages(msgs [][]interface{}) [][]byte {
    parsed := make([][]byte, 2)
    parsed[0] = []byte(fmt.Sprintf(strings.Repeat(MSG_TEMPLATE, len(msgs[0])), msgs[0]...))
    parsed[1] = []byte(fmt.Sprintf(strings.Repeat(MSG_TEMPLATE, len(msgs[1])), msgs[1]...))
    return parsed
}

func handleChat(w http.ResponseWriter, r *http.Request) {
    switch r.Method{
    case "GET":
        msgs := parseMessages(getMsgSlices(head))
        w.Write(append(msgs[0], msgs[1]...))
        break
    case "POST":
         addMsg(r.FormValue("message"))
    }
}


func main() {
    messages[0] = "Welcome to CusCus!!"

    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/chat", handleChat)
    fmt.Printf("Server running on port %s\n", PORT)
    err := http.ListenAndServe(PORT, nil)

    if err != nil{
        fmt.Println(err)
    }
}
