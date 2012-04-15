# Embeddable Hacker News vote / counter button

![HN Button](http://img.skitch.com/20120415-bp8igiq74w53f91swt6tcy9cx8.jpg)

Async, embeddable submit + vote counter button for Hacker News. If the story has not be posted to HN, it will automatically show the "Submit" button, otherwise it will report the latest point count for the story.

### Embedding the button

**Step 1**, place the HN button where you want to appear:

```html
<!-- Auto-detect URL of current page and title if necessary -->
<a href="http://news.ycombinator.com/submit" class="hn-share-button">Vote on HN</a>

<!-- Override the URL and Title for the button -->
<a href="http://news.ycombinator.com/submit" class="hn-share-button" data-title="Some Title" data-url="http://www.igvita.com/">Vote on HN</a>
```

**Step 2**, add the following loader snippet right before the closing body tag:

```html
<script>
	(function() {
	   var hn = document.createElement('script'); hn.type = 'text/javascript';
	   hn.async = true; hn.src = 'http://hnbutton.appspot.com/static/hn.js';
	   var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(hn, s);
	})();
</script>
```

_Note: you can safely embed multiple buttons on the same page._

### License

(MIT License) - Copyright (c) 2012 Ilya Grigorik