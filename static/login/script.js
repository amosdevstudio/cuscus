function deleteAllCookies() {
    const cookies = document.cookie.split(";");

    for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i];
        const eqPos = cookie.indexOf("=");
        const name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
    }
}

function setCookie(cName, cValue, expDays) {
    if (cName === undefined){return}
    let date = new Date()
    date.setTime(date.getTime() + (expDays * 24 * 60 * 60 * 1000))
    const expires = "expires=" + date.toUTCString()
    document.cookie = cName + "=" + cValue + "; " + expires + "; path=/"
}

function setSessionId(evt){
    setCookie("username", document.getElementById("username").value, 30)
    setCookie("sessionid", evt.detail.value, 30)
}

document.addEventListener("setSessionId", setSessionId)
//deleteAllCookies()
