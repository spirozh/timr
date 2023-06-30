function start() {
    let timersElt = document.getElementById("timers");

    let es = new EventSource("/api/sse/");

    let error = false;

    es.onopen = (ev) => {
        console.log('onopen(ev): ', ev)
        if (error) {
            window.location = window.location;
        }
    }

    es.onmessage = (ev) => {
        console.log('onmessage(ev): ', ev);

        msg = JSON.parse(ev.data)
        id = `timer-${msg.name}`
        e = document.getElementById(id)
        if (!e) {
            e = document.createElement('div')
            e.id = id;
            timersElt.appendChild(e);
        }
        if (!msg.state) {
            timersElt.removeChild(e);
        }
        e.innerText = ev.data;
    }
    es.onerror = (ev) => {
        console.log('onerror(ev): ', ev);

        // children = timersElt.children;
        // for (let i = 0; i < children.length; i++) {
        //     timersElt.removeChild(children[i]);
        // }
        error = true;
    }
}
