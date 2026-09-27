document.addEventListener("DOMContentLoaded", async function() {


    let chat = document.getElementById('chats')
    let snd_v = document.getElementById('snd_n')
    let searchUser = document.getElementById('searchUs')
    let result = document.getElementById('findedUs')
    let chatPlace = document.getElementById('chatPlace_n')
    let tf = false
    if (searchUser.value.length <= 2) {
        result.innerHTML = ''
    }
    if (localStorage.getItem('token') != '') {
        let goy = document.getElementById('log')
        let sec = localStorage.getItem('token')
        res = await fetch("http://localhost:7777/authkey", {
            method: 'POST',
            'Content-Type': "application/json",
            body: JSON.stringify({Login: null, Password: null, SecKey: sec})
        })
        if (res.status == 200) {
            goy.id = "prof"
            let t = await res.text()
            console.log(t)
            goy.textContent = t
            
            document.getElementById('prof').addEventListener('click', function() {
                window.location.href = 'settings.html'
            })
        }
        else {
            document.getElementById('log').addEventListener("click", function() {
                window.location.href = "logpage.html"
            })
        }
        
    }
    let inp = document.getElementById('sms_n')
    let snd = document.getElementById('snd_n')
    let forf = document.getElementById('to_')
    let usTable
    let zovv = document.getElementById('prof').textContent
    let urs = await fetch("http://localhost:7777/returnChats?lg=" + encodeURIComponent(zovv), {
        method: 'POST',
        'Content-Type': 'application/json',
        body: JSON.stringify({From: zovv})
    })
    if (urs.status == 200) {
        let ar = await urs.text()
        ar = JSON.parse(ar)
        usTable = ar
    }
    
    usTable.forEach(el => {
        let li
        if (el.split('/')[0] == document.getElementById('prof').textContent) {
            li = document.createElement('li')
            li.textContent = el.split('/')[1]
            li.classList = 'tableUser'
            li.id = 'ppsk'
            li.addEventListener('click', async function() {
                forf.id = 'to_us'
                forf.textContent = el.split('/')[1]
                inp.id = 'sms_v'
                snd.id = 'snd_v'
                chatPlace.id = 'chatPlace_v'
                
            })
            chat.appendChild(li)
        }
        else {
            let li = document.createElement('li')
            li.textContent = el.split('/')[0]
            li.classList = 'tableUser'
            li.id = 'ppsk'
            li.addEventListener('click', async function() {
                forf.id = 'to_us'
                forf.textContent = el.split('/')[0]
                inp.id = 'sms_v'
                snd.id = 'snd_v'
                chatPlace.id = 'chatPlace_v'
                
            })
            chat.appendChild(li)
        }
        
    })
    if (!localStorage.getItem('token')) {
        window.location.href = 'logpage.html'
    }
    searchUser.addEventListener('input', async function() {
        let val = this.value.trim()
        if (val.length > 2) {
            let res = await fetch('http://localhost:7777/searchUser', {
                method: 'POST',
                'Content-Type': 'application/json',
                body: JSON.stringify({Login: val})
            })
            if (res.status == 200) {
                let kpc = await res.text()
                zov = JSON.parse(kpc)
                console.log(zov)
                results(zov)
            }
        }
        else {
            result.innerHTML = ''
        }
    })
    async function results(arr) {
        result.innerHTML = ''
        arr.forEach(us => {
            let li = document.createElement('li')
            li.textContent = us
            li.classList = 'usli'
            li.addEventListener('click', async function() {
                document.getElementById('chatPlace_n').id = 'chatPlace_v'
                chatPlace.innerHTML = ''
                tf = true
                searchUser.value = us
                result.innerHTML = ''
                forf.textContent = us
                forf.id = 'to_us'
                inp.id = 'sms_v'
                snd.id = 'snd_v'
                let re = await fetch('http://localhost:7777/getMessages', {
                    method: 'POST',
                    'Content-Type': 'application/json',
                    body: JSON.stringify({From: document.getElementById('prof').textContent, To: document.getElementById('to_us').textContent})

                })
                if (re.status == 200) {
                    let aw  = await re.text()
                    aw = JSON.parse(aw) //массив с сообщ
                    aw.forEach(cont => {
                        let mms =document.createElement('li')
                        if (cont.Content != document.getElementById('prof').value) {
                            mms.className = 'msg'
                            console.log('ppsk')
                        }
                        else {
                            mms.className = 'msg_my'
                        }
                        mms.textContent = cont.Content
                        document.getElementById('chatPlace_v').appendChild(mms)
                        //document.getElementById('chatPlace_v').appendChild(document.createElement('br'))
                    })
                    
                }
                
            })
            result.appendChild(li)
        });
    }
    document.getElementById('exit').addEventListener('click', function() {
        localStorage.removeItem('token')
        window.location.href = 'logpage.html'
    })
    snd_v.addEventListener('click', async function() {
        ws.onopen = function() {
            console.log('connected')
        };
        ws.send(JSON.stringify({To: document.getElementById('to_us').textContent, Content: document.getElementById('sms_v').value}))
        let msm  = document.createElement('li')
        msm.className = 'msg_my'
        msm.textContent = document.getElementById('sms_v').value
        document.getElementById('chatPlace_v').appendChild(msm)
        document.getElementById('chatPlace_v').appendChild(document.createElement('br'))
    })
    const ws = new WebSocket('ws://localhost:7777/ws?key=' + encodeURIComponent(localStorage.getItem('token')))
    
    ws.onmessage = function(event) {
        const mes = JSON.parse(event.data)
        let mm = document.createElement('li')
        mm.className = 'msg'
        mm.textContent = mes.Content
        chatPlace.appendChild(mm)
        console.log(mes, mes.Content)
    }
    if (tf) {
        console.log('ZZZ')
        
    }
    const ppskEl = document.querySelectorAll('.tableUser')
    if (chat.childElementCount > 0) {
        ppskEl.forEach(z => {

        
            z.addEventListener('click' , async function() {
                chatPlace.innerHTML = ''
                chatPlace.id = 'chatPlace_v'
                let re = await fetch('http://localhost:7777/getMessages', {
                    method: 'POST',
                    'Content-Type': 'application/json',
                    body: JSON.stringify({From: document.getElementById('prof').textContent, To: document.getElementById('to_us').textContent})

                })
                if (re.status == 200) {
                    let aw  = await re.text()
                    aw = JSON.parse(aw) //массив с сообщ
                    aw.forEach(cont => {
                        let mms =document.createElement('li')
                        if (cont.User != document.getElementById('prof').textContent) {
                            mms.className = 'msg'
                        }
                        else {
                            mms.className = 'msg_my'
                        }
                        mms.textContent = cont.Content

                        document.getElementById('chatPlace_v').appendChild(mms)
                        //document.getElementById('chatPlace_v').appendChild(document.createElement('br'))
                    })

                }
            })
        })
    }
    inp.addEventListener('keydown', async (event) => {
        if (event.key == 'Enter') {
            ws.send(JSON.stringify({To: document.getElementById('to_us').textContent, Content: document.getElementById('sms_v').value}))
            let mm = document.createElement('li')
            mm.className = 'msg_my'
            mm.textContent = mes.Content
            chatPlace.appendChild(mm)
        }
    })
})