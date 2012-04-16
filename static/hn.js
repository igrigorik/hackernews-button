(function(w){
  var eventMethod = w.addEventListener ? "addEventListener" : "attachEvent";
  var eventer = w[eventMethod];
  var messageEvent = eventMethod == "attachEvent" ? "onmessage" : "message";

  eventer(messageEvent,function(e) {
    if (e.origin !== 'http://hnbutton.appspot.com/') { return; }
    if (typeof(_gaq) != "undefined") {
      if (e.data == 'vote' || e.data == 'submit') {
        _gaq.push(['_trackSocial', 'Hacker News', e.data]);
      }
    }
  }, false);

  var d = window.document, e,
      hnAnchorElements = d.getElementsByClassName("hn-share-button");

  for(e = hnAnchorElements.length - 1; e >= 0; e--) {
    var anchor = hnAnchorElements[e],
        title = anchor.getAttribute("data-title") || d.title,
        url = anchor.getAttribute("data-url") || w.location.href,
        i = d.createElement("iframe");

    i.src = "http://hnbutton.appspot.com/button?title="+encodeURIComponent(title)+"&url="+encodeURIComponent(url);
    i.scrolling = "auto"; i.frameBorder = "0"; i.width = "75px"; i.height = "20px";

    anchor.parentNode.insertBefore(i, anchor);
    anchor.parentNode.removeChild(anchor);
  }
})(window);