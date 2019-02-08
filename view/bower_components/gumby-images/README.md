Gumby Responsive Images
=======================

This module will preload and insert images based on media queries and/or feature detection. For more detailed documentation please check out the [Gumby docs](http://gumbyframework.com).

Installation
------------

A bower package is available to install this module into your project. We recommend using this method to install Gumby and any extra UI modules, however you can alternatively move the individuals files into your project.

	$ bower install gumby-images

Include gumby.images.js in the same fashion as your other UI modules, after gumby.js and before gumby.init.js. In production you should minify JavaScript files into a single optimized gumby.min.js file, ensuring the order (gumby.js, UI modules, gumby.init.js) is retained.

	<!--
	Include gumby.js followed by UI modules.
	Or concatenate and minify into a single file-->
	<script src="/bower_components/gumby/js/libs/gumby.js"></script>
	<script src="/bower_components/gumby/js/libs/ui/gumby.skiplink.js"></script>
	<script src="/bower_components/gumby/js/libs/ui/gumby.toggleswitch.js"></script>
	<script src="/bower_components/gumby-images/gumby.images.js"></script>
	<script src="/bower_components/gumby/js/libs/gumby.init.js"></script>

	<!-- In production minifiy and combine the above files into gumby.min.js -->
	<script src="js/gumby.min.js"></script>
	<script src="js/plugins.js"></script>
	<script src="js/main.js"></script>

Usage
-----

Using the responsive images module is simple. Add a `gumby-media` attribute to any element containing any number of comma separated media query / image path pairs. Media queries and their associated images should be separated with a pipe. You can also add a `gumby-supports` attribute to any element, containing any number of comma seaparted features / image path pairs, in the same format as `gumby-media`. Features will be tested with Modernizr so ensure your Modernizr build contains all the tests you require. Both `gumby-media` and `gumby-supports` will fallback to the image supplied in `gumby-default`. If applied to an `<img>` the `src` will be updated otherwise the `background-image` will be used.

	<img gumby-media="only screen and (max-width: 768px) and (min-width: 501px)|img/img_silence_demo-768.jpg,
					  only screen and (max-width: 500px)|img/img_silence_demo-500.jpg"
		 gumby-default="img/img_silence_demo.jpg" />

	<img gumby-supports="webp|img/img_silence_demo.webp"
			 gumby-default="img/img_silence_demo.jpg" />

### Shorthand

The characters `<` and `>` can be used as shorthand for `max-width` and `min-width` to make it easier and less verbose.

For exampe, the following...

	<div gumby-media="< 768px | 2-1-0"></div>

Will be converted to...

	<div gumby-media="(max-width: 768px) | 2-1-0"></div>

*The media queries are passed directly to [JavaScript's matchMedia function](https://developer.mozilla.org/en-US/docs/Web/API/window.matchMedia) which is not supported in <= IE9, but fear not, you can include [Paul Irish's polyfil](https://github.com/paulirish/matchMedia.js/) and all will be well*

**MIT Open Source License**

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.



