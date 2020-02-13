window.addEventListener("load", function(evt) {
    var ws = new WebSocket("ws://" + document.location.host + "/ws");

    ws.onopen = function(evt) {
        console.log("OPEN");
	var sendmsg = {
	    message: "hello",
	};
	ws.send(JSON.stringify(sendmsg));
    }
    
    ws.onclose = function(evt) {
        console.log("CLOSE"); 
        ws = null;
    }
    
    // We assume the incoming message is a JSON string containing a single
    // field 'message' with a string as a value.
    ws.onmessage = function(evt) {

	//var obj = JSON.parse(evt.data);
	var data = evt.data;
	for (int i = 0; i < evt.data.length ) {
	    var t = data[0];
	    var l = data[1];
	    var val = data[2:];
	}

	/*
	for (id in obj) {
	    var ele = document.getElementById(id);
	    if (!ele) {
		console.log("Unknown element: " + id);
		continue;
	    }

	    switch (id) {
	    default:
		if (ele) {
		    ele.innerHTML = obj[id];
		}
	    }
	}
	*/
    }
    
    ws.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    }
    return false;
});
