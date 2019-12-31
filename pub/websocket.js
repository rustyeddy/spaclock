window.addEventListener("load", function(evt) {
    var clock = document.getElementById("clock");
    var ws = new WebSocket("ws://localhost:8000/ws");

    ws.onopen = function(evt) {
        console.log("OPEN");
	ws.send("hello");
    }
    
    ws.onclose = function(evt) {
        console.log("CLOSE");
        ws = null;
    }
    
    ws.onmessage = function(evt) {
        console.log("RESPONSE: " + evt.data);
	clock.textContent = evt.data;
    }
    
    ws.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    }
    return false;
});
