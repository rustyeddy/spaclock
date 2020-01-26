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
	var obj = JSON.parse(evt.data);
	for (id in obj) {
	    console.log(id + " - " + obj[id]);
	    switch (id) {
	    case "clock":
		console.log("  .. skipping clock");
		// Do not update clock
		break;

	    default:
		var ele = document.getElementById(id);
		if (ele) {
		    console.log("  .. " + id + " == " + obj[id]);
		    ele.innerHTML = obj[id];
		}
	    }
	}
    }
    
    ws.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    }
    return false;
});
