$(function() {
	var run = new Run();
});

var Run = function() {

	this.stopWatchInterval = null;
	this.stopWatchElem = null;

	if (window["WebSocket"]) {
		var conn = new WebSocket(this.getWsUrl());

		conn.onclose = function(evt) {
			alert("I has no connection");
		};

		conn.onmessage = this.onMessage.bind(this)

		this.conn = conn
	}
	// What to do?
};

Run.prototype.getWsUrl = function() {
	var loc = window.location,
		new_uri;
	if (loc.protocol === "https:") {
		new_uri = "wss:";
	} else {
		new_uri = "ws:";
	}
	new_uri += "//" + loc.host;
	new_uri += "/ws";
	return new_uri

};

Run.prototype.onMessage = function(msg) {
	// we should figure out what type of message this is ... its going to be ugly
	var msg = JSON.parse(msg.data);

	if (msg.durationNs) {
		this.runEnd(msg);
		return;
	}
	if (msg.started) {
		this.runStart(msg);
		return;
	}

	console.log("Dont know what this message means:", msg);
};

Run.prototype.runEnd = function(data) {
	clearInterval(this.stopWatchInterval);
	this.stopWatchInterval = null;
	$("#content").html("Tillykke!: <pre>" + (data.durationNs / 1000 / 1000).toHHMMSS() + "<pre>");
};

Run.prototype.updateStopwatch = function() {
	this.stopWatchElem.text((Date.now() - this.stopWatchStart).toHHMMSS())
};

Run.prototype.runStart = function(data) {
	
	if(this.stopWatchInterval) {
		clearInterval(this.stopWatchInterval);
	}

	this.stopWatchStart = Date.now();
	this.stopWatchElem = $(document.createElement("pre")).text("00:00:00.00");
	$("#content").html("Løøøb!").append(this.stopWatchElem);
	this.stopWatchInterval = setInterval(this.updateStopwatch.bind(this), 10);
};



Number.prototype.toHHMMSS = function() {
	var seconds = Math.floor(this) / 1000,
		hours = Math.floor(seconds / 3600);
	seconds -= hours * 3600;
	var minutes = Math.floor(seconds / 60);
	seconds -= minutes * 60;

	if (hours < 10) {
		hours = "0" + hours;
	}
	if (minutes < 10) {
		minutes = "0" + minutes;
	}
	if (seconds < 10) {
		seconds = "0" + seconds.toFixed(2);
	} else {
		seconds = seconds.toFixed(2);
	}
	return hours + ':' + minutes + ':' + seconds;
}