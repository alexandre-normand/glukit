/**
* Gumby Images
*/
!function() {

	'use strict';

	function Images($el) {

		Gumby.debug('Initializing Responsive Images module', $el);

		this.$el = $el;

		this.type = '';
		this.supports = '';
		this.media = '';
		this.def = '';
		this.current = '';

		// set up module based on attributes
		this.setup();

		var scope = this;
		$(window).on('load gumby.trigger '+(!this.media || 'resize'), function(e) {
			scope.fire();
		});

		this.$el.on('gumby.initialize', function() {
			Gumby.debug('Re-initializing Responsive Images module', scope.$el);
			scope.setup();
			scope.fire();
		});

		scope.fire();
	}

	// set up module based on attributes
	Images.prototype.setup = function() {
		// is this an <img> or background-image?
		this.type = this.$el.is('img') ? 'img' : 'bg';
		// supports attribute in format test:image
		this.supports = Gumby.selectAttr.apply(this.$el, ['supports']) || false;
		// media attribute in format mediaQuery:image
		this.media = Gumby.selectAttr.apply(this.$el, ['media']) || false;
		// default image to load
		this.def = Gumby.selectAttr.apply(this.$el, ['default']) || false;

		// parse support/media objects
		if(this.supports) { this.supports = this.parseAttr(this.supports); }
		if(this.media) { this.media = this.parseAttr(this.media); }

		// check functions
		this.checks = {
			supports : function(val) {
				return Modernizr[val];
			},
			media: function(val) {
				return window.matchMedia(val).matches;
			}
		};
	};

	// fire required checks and load resulting image
	Images.prototype.fire = function() {
		// feature supported or media query matched
		var success = false;

		// if support attribute supplied and Modernizr is present
		if(this.supports && Modernizr) {
			success = this.handleTests('supports', this.supports);
		}

		// if media attribute supplied and matchMedia is supported
		// and success is still false, meaning no supporting feature was found
		if(this.media && window.matchMedia && !success) {
			success = this.handleTests('media', this.media);
		}

		// no feature supported or media query matched so load default if supplied
		if(!success && this.def) {
			success = this.def;
		}

		// no image to load
		if(!success) {
			return false;
		}

		// preload image and insert or set background-image property if not already set
		if(this.current !== success) {
			this.current = success;
			this.insertImage(this.type, success);
		}
	};

	// handle media object checking each prop for matching media query
	Images.prototype.handleTests = function(type, array) {
		var scope = this,
			supported = false;

		$(array).each(function(key, val) {
			// media query matched
			// supplied in order of preference so halt each loop
			if(scope.check(type, val.test)) {
				supported = val.img;
				return false;
			}
		});

		return supported;
	};

	// return the result of test function
	Images.prototype.check = function(type, val) {
		return this.checks[type](val);
	};

	// preload image and insert or set background-image property
	Images.prototype.insertImage = function(type, img) {
		var scope = this,
			image = $(new Image());

		image.load(function() {
			type === 'img' ? scope.$el.attr('src', img) : scope.$el.css('background-image', 'url('+img+')');

			// trigger custom loaded event
			Gumby.debug('Triggering onChange event', img, scope.$el);
			scope.$el.trigger('gumby.onChange', [img]);
		}).attr('src', img);
	};

	// parse attribute strings, media/support
	Images.prototype.parseAttr = function(support) {
		var scope = this,
			supp = support.split(','),
			res = [], splt = [];

		// multiple can be supplied so loop round and create object
		$(supp).each(function(key, val) {
			splt = val.split('|');
			if(splt.length !== 2) {
				return true;
			}

			// object containing Modernizr test or media query and image url
			res.push({
				'test' : scope.shorthand(splt[0]),
				'img' : splt[1]
			});
		});

		return res;
	};

	// replace < and > with min/max width media queries
	Images.prototype.shorthand = function(str) {
		// replace < and >
		if(str.indexOf('>') > -1 || str.indexOf('<') > -1) {
			str = str.replace('>', 'min-width: ').replace('<', 'max-width: ');
		}

		// check if media query (prevent wrapping feature detection tests) AND if media query, wrapped in ()
 		if(str.indexOf('width') > -1 
 				&& str.charAt(0) !== '(' && str.charAt(str.length - 1) !== ')') {
			str = '('+str+')';
		}

		return str;
	};

	// add initialisation
	Gumby.addInitalisation('images', function(all) {
		$('[gumby-supports],[data-supports],[supports],[gumby-media],[data-media],[media]').each(function() {
			var $this = $(this);

			// this element has already been initialized
			// and we're only initializing new modules
			if($this.data('isImage') && !all) {
				return true;

			// this element has already been initialized
			// and we need to reinitialize it
			} else if($this.data('isImage') && all) {
				$this.trigger('gumby.initialize');
				return true;
			}

			// mark element as initialized
			$this.data('isImage', true);
			new Images($this);
		});
	});

	// register UI module
	Gumby.UIModule({
		module: 'images',
		events: ['onChange', 'trigger'],
		init: function() {
			Gumby.initialize('images');
		}
	});
}();
