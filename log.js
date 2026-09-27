document.addEventListener("DOMContentLoaded", async function() {
    if (localStorage.getItem('token') != '') {
        let s = localStorage.getItem('token')
        let re = await fetch('http://localhost:7777/authkey', {
            method: 'POST',
            'Content-Type': 'application/json',
            body: JSON.stringify({Login: null, Password: null, SecKey: s})
        })
        if (re.status == 200) {
            window.location.href = 'main.html'
        }
    }
    let cht = '1'
})

async function kpc(event) {
    let btn = event.currentTarget
    let btnName = btn.id
    if (btnName == "auth") {
        let log = document.getElementById('login').value
        let pass = document.getElementById('password').value
        let tt = log + "/" + pass + ":" + String(log.length + pass.length) + "//" + String(log.length + pass.length)
        let res = await fetch("http://localhost:7777/auth", {
            method: 'POST',
            'Content-Type': 'application/json',
            body: JSON.stringify({Login: log, Password: pass, SecKey: btoa(btoa(tt))})
        })
        if (res.status == 200) {
            key = (await res.text()).toString()
            let sec = atob(atob(key))
            localStorage.setItem('token', key)
            window.location.href = "main.html"

        } 
        if (res.status == 502) {
            alert('password is incorrect')
        }
    }
    else if (btnName == "authkey") {
        let key = document.getElementById('skey').value
        let res = await fetch("http://localhost:7777/authkey", {
            method: 'POST',
            'Content-Type': 'application/json',
            body: JSON.stringify({Login: null, Password: null,SecKey: key})
        })
        if (res.status == 200) {
            window.location.href = 'main.html'
            localStorage.setItem('token', key)

        }
        if (res.status == 502) {
            alert('key is incorrect')
        }
        if (res.status == 409) {
            alert('db error')
        }
        if (res.status == 424) {
            alert('user not found')
        }
    }

}