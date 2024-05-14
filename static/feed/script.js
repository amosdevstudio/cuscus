let lastMsg = 0
let username = ""
let pwd = ""

function changeLastMsg(){
    let el = document.getElementById("tmp-lastMsg")
    if (el == null){return}
    lastMsg = el.innerText
    el.remove()
}

let cookie = document.cookie.split(";")
username = cookie[0].split('=')[1]
pwd = cookie[1].split('=')[1]

