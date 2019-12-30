function startTime() {
    /*
    var today = new Date();
    var h = today.getHours();
    var m = today.getMinutes();
    var s = today.getSeconds();
    m = checkTime(m);
    s = checkTime(s);
    document.getElementById('clock').innerHTML =
	"<span id='clock' class='text-center'>" + h + ":" + m + ":" + s + "</span>";
    var t = setTimeout(startTime, 500);
    */
    var clock = document.getElementById("clock");

    if (ws) {
	return
    }
    //var ws = new WebSocket("ws://{{ .Host }}/ws");
    var ws = new WebSocket("ws://localhost:8000/ws")
    ws.onopen = function(evt) {
	console.log("websocket opening")
    }

    ws.onclose = function(evt) {
	data.textContent = 'Connection Closed';
    }

    ws.onmessage = function(evt) {
	console.log('update over websocket');
	data.textContent = evt.Data
    }

    ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }

}

function checkTime(i) {
    if (i < 10) {i = "0" + i};  // add zero in front of numbers < 10
    return i;
}
