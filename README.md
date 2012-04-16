# Embeddable Hacker News vote / counter button

![HN Button](http://img.skitch.com/20120415-bp8igiq74w53f91swt6tcy9cx8.jpg)

Async, embeddable submit + vote counter button for Hacker News.

- If the story has not been posted to HN, "Submit" button is shown, otherwise latest point count is displayed. 
- Auto-detects Google Analytics and registers clicks events (see reports under `Traffic Sources > Social > Social Plugins`).

### Embedding the button

**Step 1**, place the HN link where you want the button appear on the page:

```html
<!-- Auto-detect URL of current page and title if necessary -->
<a href="http://news.ycombinator.com/submit" class="hn-share-button">Vote on HN</a>

<!-- Override the URL and Title for the button -->
<a href="http://news.ycombinator.com/submit" class="hn-share-button" data-title="Some Title" data-url="http://www.igvita.com/">Vote on HN</a>
```

**Step 2**, add the following loader snippet right before the `</body>` tag:

```html
<script>
	(function(d, t) {
		var g = d.createElement(t),
		    s = d.getElementsByTagName(t)[0];
		g.src = '//hnbutton.appspot.com/static/hn.js';
		s.parentNode.insertBefore(g, s);
	}(document, 'script'));
</script>
```

_Note: you can safely embed multiple buttons on the same page._

### Misc

* Kudos to @sbashyal and @stbullard for the button styling (hnlike.com)
* (MIT License) - Copyright (c) 2012 Ilya Grigorik