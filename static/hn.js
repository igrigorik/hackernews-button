/*jshint forin:true, noarg:true, noempty:true, eqeqeq:true, bitwise:true, strict:true, undef:true, curly:true, browser:true, indent:2, maxerr:50, expr:true */
(function (w) {
  "use strict";
  var j,
    d = w.document,
    getElementsByClassName = function(match, tag) {
      if (d.getElementsByClassName) {
        return d.getElementsByClassName(match);
      }
      var result = [],
        elements = d.getElementsByTagName(tag || '*'),
        i, elem; 
      match = " " + match + " ";
      for (i = 0; i < elements.length; i++) { 
        elem = elements[i];
        if ((" " + (elem.className || elem.getAttribute("class")) + " ").indexOf(match) > -1) {
          result.push(elem);
        }
      }
      return result; 
    },
    hnAnchorElements = getElementsByClassName("hn-share-button", "a"),
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
    i.className = "hn-share-iframe";

    anchor.parentNode.insertBefore(i, anchor);
    anchor.parentNode.removeChild(anchor);
  }
})(window);