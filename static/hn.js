(function(d){
	var eventMethod = window.addEventListener ? "addEventListener" : "attachEvent";
	var eventer = window[eventMethod];
	var messageEvent = eventMethod == "attachEvent" ? "onmessage" : "message";

	eventer(messageEvent,function(e) {
		if (typeof(_gaq) != "undefined") {
			_gaq.push(['_trackSocial', 'hacker news', e.data])
		}
	},false);

	var hnAnchorElements = document.getElementsByClassName("hn-share-button");
	for(var e = hnAnchorElements.length - 1; e >= 0; e--) {
		title = hnAnchorElements[e].getAttribute("data-title") || window.document.title;
		url   = hnAnchorElements[e].getAttribute("data-url")   || window.location.href;

		var i = document.createElement("iframe");
		i.src = "http://hnbutton.appspot.com/button?title="+encodeURIComponent(title)+"&url="+encodeURIComponent(url);
		i.scrolling = "auto"; i.frameBorder = "0"; i.width = "75px"; i.height = "20px";

		hnAnchorElements[e].parentNode.insertBefore(i, hnAnchorElements[e]);
		hnAnchorElements[e].parentNode.removeChild(hnAnchorElements[e]);
	}
})(document);