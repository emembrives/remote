(function(){

	// the minimum version of jQuery we want
	var vmin = "1.3";
	var v = "1.11.1";

	// check prior inclusion and version
	if (window.jQuery === undefined || window.jQuery.fn.jquery < vmin) {
		var done = false;
		var script = document.createElement("script");
		script.src = "http://ajax.googleapis.com/ajax/libs/jquery/" + v + "/jquery.min.js";
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
	
	function initMyBookmarklet() {
		(window.myBookmarklet = function() {
			var socket = new WebSocket("ws://etienne.membrives.fr/remote/conn");
			socket.onmessage = function(message) {
				if (message.data == "a") {
					// left
					jQuery.event.trigger({ type : 'keypress', which : 37 });
				} else if (message.data == "d") {
					// right
					jQuery.event.trigger({ type : 'keypress', which : 39 });
				}
			}
			socket.onerror = function(error) {
				console.log(error);
			}
		})();
	}

})();

