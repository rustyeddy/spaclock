window.addEventListener("load", function(evt) {
    var ws = new WebSocket("ws://" + document.location.host + "/ws");
    //ws.binaryType = 'arraybuffer';

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
	if (obj == null) {
	    console.log("WS bummer to message");
	    return;
	}

	console.log(obj);
	for (id in obj) {
	    var o = obj[id];
	    var key = o['k'];
	    var val = o['v'];

	    var ele = document.getElementById(key);
	    if (!ele) {
		console.log("Unknown element: " + id);
		continue;
	    }

	    ele.innerHTML = val;
	}
    }
    
    ws.onerror = function(evt) {
        console.log("ERROR: " + evt.data);
    }
    return false;
});
