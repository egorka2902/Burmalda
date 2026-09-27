document.addEventListener('DOMContentLoaded', async function () {
    document.getElementById('back').addEventListener('click', function() {
        window.location.href = 'main.html'
    })


    let f = document.getElementById('pa')
    let s = document.getElementById('ke')
    let pkc = localStorage.getItem('token')
    let a = f.placeholder
    let b = s.placeholder
    let k
    document.getElementById('izm').addEventListener('click', function() {
        if (f.placeholder == a) {
            f.placeholder = b
            s.placeholder = a
            let temp = f.value
            f.value = s.value
            s.value = temp
            f.type = 'password'
            s.type = 'text'
        } else {
            f.placeholder = a
            s.placeholder = b
            let temp = f.value
            f.value = s.value
            s.value = temp
            s.type = 'password'
            f.type = 'text'
        }

    })
    document.getElementById('sendChange').addEventListener('click', async function() {
        if (f.placeholder == a) {
            k = s.value
        }
        if (f.placeholder == b) {
            k = f.value
        }
        let login
        let tmp
        if (document.getElementById('shifr').checked) {
            tmp = btoa(btoa(k))
        } else {
            tmp = k
        }
        let logg = await fetch('http://localhost:7777/returnLogin', {
            method: "POST",
            'Content-Type': 'application/json',
            body: JSON.stringify({Login: null, Password: null, SecKey: localStorage.getItem('token')})
        })
        if (logg.status == 200) {
            login = await logg.text()
        }
        let res = await fetch('http://localhost:7777/changeData', {
            method: 'POST',
            'Content-Type': 'application/json',
            body: JSON.stringify({Login: login, Password: f.value, SecKey: tmp})
        })
        if (res.status == 200) {
            if (tmp != pkc) {
                alert("Key copied to clipboard")
                await navigator.clipboard.writeText(tmp)
            }
            else {
                alert("Password was changed")
            }
            localStorage.removeItem('token')
            window.location.href = "logpage.html"

        }
        
    })
    document.getElementById('copyKey').addEventListener('click', async function() {
        await navigator.clipboard.writeText(pkc)
    })
})