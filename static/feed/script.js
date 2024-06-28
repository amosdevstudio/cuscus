let msgCount = 0

function setMsgCount(evt){
    msgCount = evt.detail.value
}

let sessionid = document.cookie.split("sessionid=")[1]
let username = document.cookie.split("username=")[1].split("; sessionid=")[0]

document.body.addEventListener("setMsgCount", setMsgCount)
