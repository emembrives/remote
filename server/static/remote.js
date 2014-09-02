(function(){

	// the minimum version of jQuery we want
	var vmin = "1.3";
	var v = "1.11.1";

	// check prior inclusion and version
	if (window.jQuery === undefined || window.jQuery.fn.jquery < vmin) {
		var done = false;
		var script = document.createElement("script");
		script.src = "https://ajax.googleapis.com/ajax/libs/jquery/" + v + "/jquery.min.js";
		script.onload = script.onreadystatechange = function(){
			if (!done && (!this.readyState || this.readyState == "loaded" || this.readyState == "complete")) {
				done = true;
				initMyBookmarklet();
			}
		};
		document.getElementsByTagName("head")[0].appendChild(script);
	} else {
		initMyBookmarklet();
	}
	
	function dispatchEvent(keyCode) {
		var eventObj = document.createEventObject ?
        	document.createEventObject() : document.createEvent("Events");
		if(eventObj.initEvent){
			eventObj.initEvent("keydown", true, true);
		}

		eventObj.keyCode = keyCode;
		eventObj.which = keyCode;

		document.body.dispatchEvent ? document.body.dispatchEvent(eventObj) : document.body.fireEvent("onkeydown", eventObj);
	}

	function initMyBookmarklet() {
		(window.myBookmarklet = function() {
			var socket = new WebSocket("wss://etienne.membrives.fr/remote/");
			socket.onmessage = function(message) {
				var key = JSON.parse(message.data)
				if (key == "a") {
					// left
					dispatchEvent(37);
				} else if (key == "d") {
					// right
					dispatchEvent(39);
				}
			}
			socket.onerror = function(error) {
				console.log(error);
			}
		})();
	}

})();

