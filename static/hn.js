(function(window){
	var document = window.document,
		hnAnchorElements = document.getElementsByClassName("hn-share-button"),
		e;
	for (e = hnAnchorElements.length - 1; e >= 0; e--) {
		var anchor = hnAnchorElements[e],
			title = anchor.getAttribute("data-title") || document.title,
			url = anchor.getAttribute("data-url")   || window.location.href,
			i = document.createElement("iframe");
		i.src = "http://hnbutton.appspot.com/button?title="+encodeURIComponent(title)+"&url="+encodeURIComponent(url);
		i.scrolling = "auto"; i.frameBorder = "0"; i.width = "75px"; i.height = "20px";

		anchor.parentNode.insertBefore(i, anchor);
		anchor.parentNode.removeChild(anchor);
	}
})(window);