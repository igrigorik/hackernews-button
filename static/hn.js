/*jshint forin:true, noarg:true, noempty:true, eqeqeq:true, bitwise:true, strict:true, undef:true, curly:true, browser:true, indent:2, maxerr:50, expr:true */
(function (w) {
  "use strict";
  var j,
    d = w.document,
    hnAnchorElements = d.getElementsByClassName("hn-share-button"),
    eventMethod = w.addEventListener ? "addEventListener" : "attachEvent",
    eventer = w[eventMethod],
    messageEvent = eventMethod === "attachEvent" ? "onmessage" : "message",
    base = "http://hnbutton.appspot.com/";

  w._gaq || (w._gaq = []);
  eventer(messageEvent, function (e) {
    if (e.origin === base && (e.data === "vote" || e.data === "submit")) {
      w._gaq.push(["_trackSocial", "Hacker News", e.data]);
    }
  }, false);

  for (j = hnAnchorElements.length - 1; j >= 0; j--) {
    var anchor = hnAnchorElements[j],
      title = anchor.getAttribute("data-title") || d.title,
      url = anchor.getAttribute("data-url") || w.location.href,
      i = d.createElement("iframe");

    i.src = base + "button?title=" + encodeURIComponent(title) + "&url=" + encodeURIComponent(url);
    i.scrolling = "auto"; 
    i.frameBorder = "0"; 
    i.width = "75px"; 
    i.height = "20px";

    anchor.parentNode.insertBefore(i, anchor);
    anchor.parentNode.removeChild(anchor);
  }
})(window);