/* -------------------- Check Browser --------------------- */

function browser() {
	
	//var isOpera = !!(window.opera && window.opera.version);  // Opera 8.0+
	//var isFirefox = testCSS('MozBoxSizing');                 // FF 0.8+
	var isSafari = Object.prototype.toString.call(window.HTMLElement).indexOf('Constructor') > 0;
	    // At least Safari 3+: "[object HTMLElementConstructor]"
	var isChrome = !isSafari && testCSS('WebkitTransform');  // Chrome 1+
	//var isIE = /*@cc_on!@*/false || testCSS('msTransform');  // At least IE6

	function testCSS(prop) {
	    return prop in document.documentElement.style;
	}
	
	if (isSafari || isChrome) {
		
		return true;
		
	} else {
		
		return false;
		
	}
	
}

/* -------------------- Placeholder for IE --------------------- */

jQuery(document).ready(function($) {

	// Invoke the plugin
    $('input, textarea').placeholder();
    // That’s it, really.
    // Now display a message if the browser supports placeholder natively
    var html;
    
});

/* -------------------- Charts --------------------- */

jQuery(document).ready(function($){
		
	function randNum(){
		return (Math.floor( Math.random()* (1+40-20) ) ) + 20;
	}
	
	if($("#stockPrice").length)
	{
		var price = [[1, randNum()-10], [2, randNum()-10], [3, randNum()-10], [4, randNum()],[5, randNum()],[6, 4+randNum()],[7, 5+randNum()],[8, 6+randNum()],[9, 6+randNum()],[10, 8+randNum()],[11, 9+randNum()],[12, 10+randNum()],[13,11+randNum()],[14, 12+randNum()],[15, 13+randNum()],[16, 14+randNum()],[17, 15+randNum()],[18, 15+randNum()],[19, 16+randNum()],[20, 17+randNum()],[21, 18+randNum()],[22, 19+randNum()],[23, 20+randNum()],[24, 21+randNum()],[25, 14+randNum()],[26, 24+randNum()],[27,25+randNum()],[28, 26+randNum()],[29, 27+randNum()], [30, 31+randNum()]];

		var plot = $.plot($("#stockPrice"),
			   [ { data: price, label: "price" } ], {
				   series: {
					   lines: { show: true,
								lineWidth: 2,
								fill: true, fillColor: { colors: [ { opacity: 0.5 }, { opacity: 0.2 } ] }
							 },
					   points: { show: true },
					   shadowSize: 1
				   },
				   grid: { hoverable: true, 
						   clickable: true, 
						   tickColor: "#eee",
						   borderWidth: 0,
						 },
				   colors: ["#414141"],
					xaxis: {ticks:11, tickDecimals: 0},
					yaxis: {tickFormatter: function (v) { return v + "USD"; }},
				 });

		function showTooltip(x, y, contents) {
			$('<div id="tooltip">' + contents + '</div>').css( {
				position: 'absolute',
				display: 'none',
				top: y + 5,
				left: x + 5,
				border: '1px solid #fdd',
				padding: '2px',
				'background-color': '#dfeffc',
				opacity: 0.80
			}).appendTo("body").fadeIn(200);
		}

		var previousPoint = null;
		$("#stockPrice").bind("plothover", function (event, pos, item) {
			$("#x").text(pos.x.toFixed(2));
			$("#y").text(pos.y.toFixed(2));

				if (item) {
					if (previousPoint != item.dataIndex) {
						previousPoint = item.dataIndex;

						$("#tooltip").remove();
						var x = item.datapoint[0].toFixed(2),
							y = item.datapoint[1].toFixed(2);

						showTooltip(item.pageX, item.pageY,
									item.series.label + " of " + x + " = " + y);
					}
				}
				else {
					$("#tooltip").remove();
					previousPoint = null;
				}
		
		});
		
	}
	
	function randNumFB(){
		return ((Math.floor( Math.random()* (1+40-20) ) ) + 20);
	}
	
	/* ---------- Chart with points ---------- */
	if($("#facebookChart").length)
	{	
		var likes = [[1, 5+randNumFB()], [2, 10+randNumFB()], [3, 15+randNumFB()], [4, 20+randNumFB()],[5, 25+randNumFB()],[6, 30+randNumFB()],[7, 35+randNumFB()],[8, 40+randNumFB()],[9, 45+randNumFB()],[10, 50+randNumFB()],[11, 55+randNumFB()],[12, 60+randNumFB()],[13, 65+randNumFB()],[14, 70+randNumFB()],[15, 75+randNumFB()],[16, 80+randNumFB()],[17, 85+randNumFB()],[18, 90+randNumFB()],[19, 85+randNumFB()],[20, 80+randNumFB()],[21, 75+randNumFB()],[22, 80+randNumFB()],[23, 75+randNumFB()],[24, 70+randNumFB()],[25, 65+randNumFB()],[26, 75+randNumFB()],[27,80+randNumFB()],[28, 85+randNumFB()],[29, 90+randNumFB()], [30, 95+randNumFB()]];

		var plot = $.plot($("#facebookChart"),
			   [ { data: likes, label: "Fans"} ], {
				   series: {
					   lines: { show: true,
								lineWidth: 2,
								fill: true, fillColor: { colors: [ { opacity: 0.5 }, { opacity: 0.2 } ] }
							 },
					   points: { show: true, 
								 lineWidth: 2 
							 },
					   shadowSize: 0
				   },
				   grid: { hoverable: true, 
						   clickable: true, 
						   tickColor: "#ddd",
						   borderWidth: 0
						 },
				   colors: ["#3B5998"],
					xaxis: {ticks:6, tickDecimals: 0},
					yaxis: {ticks:3, tickDecimals: 0},
				 });

		function showTooltip(x, y, contents) {
			$('<div id="tooltip">' + contents + '</div>').css( {
				position: 'absolute',
				display: 'none',
				top: y + 5,
				left: x + 5,
				border: '1px solid #fdd',
				padding: '2px',
				'background-color': '#dfeffc',
				opacity: 0.80
			}).appendTo("body").fadeIn(200);
		}

		var previousPoint = null;
		$("#facebookChart").bind("plothover", function (event, pos, item) {
			$("#x").text(pos.x.toFixed(2));
			$("#y").text(pos.y.toFixed(2));

				if (item) {
					if (previousPoint != item.dataIndex) {
						previousPoint = item.dataIndex;

						$("#tooltip").remove();
						var x = item.datapoint[0].toFixed(2),
							y = item.datapoint[1].toFixed(2);

						showTooltip(item.pageX, item.pageY,
									item.series.label + " of " + x + " = " + y);
					}
				}
				else {
					$("#tooltip").remove();
					previousPoint = null;
				}
		});
	
	}
	
	function randNumTW(){
		return ((Math.floor( Math.random()* (1+40-20) ) ) + 20);
	}
	
	/* ---------- Chart with points ---------- */
	if($("#twitterChart").length)
	{	
		var followers = [[1, 5+randNumTW()], [2, 10+randNumTW()], [3, 15+randNumTW()], [4, 20+randNumTW()],[5, 25+randNumTW()],[6, 30+randNumTW()],[7, 35+randNumTW()],[8, 40+randNumTW()],[9, 45+randNumTW()],[10, 50+randNumTW()],[11, 55+randNumTW()],[12, 60+randNumTW()],[13, 65+randNumTW()],[14, 70+randNumTW()],[15, 75+randNumTW()],[16, 80+randNumTW()],[17, 85+randNumTW()],[18, 90+randNumTW()],[19, 85+randNumTW()],[20, 80+randNumTW()],[21, 75+randNumTW()],[22, 80+randNumTW()],[23, 75+randNumTW()],[24, 70+randNumTW()],[25, 65+randNumTW()],[26, 75+randNumTW()],[27,80+randNumTW()],[28, 85+randNumTW()],[29, 90+randNumTW()], [30, 95+randNumTW()]];

		var plot = $.plot($("#twitterChart"),
			   [ { data: followers, label: "Followers"} ], {
				   series: {
					   lines: { show: true,
								lineWidth: 2,
								fill: true, fillColor: { colors: [ { opacity: 0.5 }, { opacity: 0.2 } ] }
							 },
					   points: { show: true, 
								 lineWidth: 2 
							 },
					   shadowSize: 0
				   },
				   grid: { hoverable: true, 
						   clickable: true, 
						   tickColor: "#ddd",
						   borderWidth: 0
						 },
				   colors: ["#1BB2E9"],
					xaxis: {ticks:6, tickDecimals: 0},
					yaxis: {ticks:3, tickDecimals: 0},
				 });

		function showTooltip(x, y, contents) {
			$('<div id="tooltip">' + contents + '</div>').css( {
				position: 'absolute',
				display: 'none',
				top: y + 5,
				left: x + 5,
				border: '1px solid #fdd',
				padding: '2px',
				'background-color': '#dfeffc',
				opacity: 0.80
			}).appendTo("body").fadeIn(200);
		}

		var previousPoint = null;
		$("#twitterChart").bind("plothover", function (event, pos, item) {
			$("#x").text(pos.x.toFixed(2));
			$("#y").text(pos.y.toFixed(2));

				if (item) {
					if (previousPoint != item.dataIndex) {
						previousPoint = item.dataIndex;

						$("#tooltip").remove();
						var x = item.datapoint[0].toFixed(2),
							y = item.datapoint[1].toFixed(2);

						showTooltip(item.pageX, item.pageY,
									item.series.label + " of " + x + " = " + y);
					}
				}
				else {
					$("#tooltip").remove();
					previousPoint = null;
				}
		});
	
	}

});

/* -------------------- Search --------------------- */

jQuery(document).ready(function($){
	
	$('.search > :input').on('keyup',function(){
		
		$(this).attr('class', 'activeSearch');

		var count;
		var timeToEnd = 1000;

		$(':input').keydown(function(){

			clearTimeout(count);

			count = setTimeout(endCount, timeToEnd);

		});

	});

	function endCount(){

		$('.search > :input').attr('class','search-form');

	}

});

/* -------------------- Buttons 3D Style --------------------- */

jQuery(document).ready(function($){

	$(':button').each(function(){

		if($(this).hasClass('btn')) {

			$(this).wrap('<div class="btn-overlay" />');

			var inner = $(this).html();
			
			if(browser()) {

				$(this).html('<span>' + inner + '</span>');
				
			}	

		}

	});
	
		
	$('a').each(function(){

		if($(this).hasClass('btnOverlay')) {
			
			$(this).wrap('<div class="btn-overlay" />');

			var inner = $(this).html();
			
			if(browser()) {

				$(this).html('<span>' + inner + '</span>');
				
			}	

		}

	});

});

/* -------------------- Twitter --------------------- */

jQuery(document).ready(function($){
	
	$.getJSON('http://api.twitter.com/1/statuses/user_timeline/lukaszholeczek.json?count=3&callback=?', function(tweets){
		$("#twitter").html(tz_format_twitter(tweets));
	}); 

});

/* ------------------- Fancybox --------------------- */
jQuery(document).ready(function($){
	
	(function() {

		$('[rel=image]').fancybox({
			type        : 'image',
			openEffect  : 'fade',
			closeEffect	: 'fade',
			nextEffect  : 'fade',
			prevEffect  : 'fade',
			helpers     : {
				title   : {
					type : 'inside'
				}
			}
		});

		$('[rel=image-gallery]').fancybox({
			nextEffect  : 'fade',
			prevEffect  : 'fade',
			helpers     : {
				title   : {
					type : 'inside'
				},
				buttons  : {},
				media    : {}
			}
		});


	})();
	
});
	
/* ------------------- Client Carousel --------------------- */

jQuery(document).ready(function($){

	$('.clients-carousel').flexslider({
	    animation: "slide",
		easing: "swing",
	    animationLoop: true,
	    itemWidth: 200,
	    itemMargin: 1,
	    minItems: 1,
	    maxItems: 8,
		controlNav: false,
		directionNav: false,
		move: 2
     });

});

/* ------------------ Back To Top ------------------- */
jQuery(document).ready(function($){

	jQuery('#under-footer-back-to-top a').click(function(){
		jQuery('html, body').animate({scrollTop:0}, 300); 
		return false; 
	});

});	

/* ------------------ Tooltips ----------------- */

jQuery(document).ready(function($) {

    $('.tooltips').tooltip({
      selector: "a[rel=tooltip]"
    })

});

/* ------------------ Progress Bar ------------------- */	

jQuery(document).ready(function($){
	
	$('.meter > span').each(function() {
		
		var totalWidth = $(this).parent().width();
		
		var barWidth = $(this).width();
		
		var percent = barWidth/totalWidth * 100;

		$(this).data('origWidth', $(this).width()).width(0).animate({
			width: $(this).data('origWidth')
		}, 1200, function(){
			$(this).css('width', percent + '%');
		});

		
	});
});

/* ------------------- Parallax --------------------- */

jQuery(document).ready(function($){
	
	$('#da-slider').cslider({
		autoplay	: true,
		bgincrement	: 0
	});

});

/* ------------------- Layered slider --------------------- */

jQuery(document).ready(function($){
	
	$('#layerslider').layerSlider({
		skinsPath : 'css/skins/',
		skin : 'tabula',
		thumbnailNavigation : 'always',				
	});
	
});		


/* -------------------- Isotope --------------------- */

jQuery(document).ready(function () {
	
	$('#wall').imagesLoaded(function() {
		
		var $container = $('#wall');
			$select = $('#filters select');

		// initialize Isotope
		$container.isotope({
		// options...
		resizable: false, // disable normal resizing
		// set columnWidth to a percentage of container width
	  	masonry: { columnWidth: $container.width() / 12 }
		});

		// update columnWidth on window resize
		$(window).smartresize(function(){
		
			$container.isotope({
			// update columnWidth to a percentage of container width
				masonry: { columnWidth: $container.width() / 12 }
			});
		});


		$container.isotope({
			itemSelector : '.item'
		});

		$select.change(function() {
			
			var filters = $(this).val();

				$container.isotope({
					filter: filters
				});
			
			});

			var $optionSets = $('#filters .option-set'),
		  	$optionLinks = $optionSets.find('a');

		  	$optionLinks.click(function(){
			
				var $this = $(this);
				// don't proceed if already selected
				if ( $this.hasClass('selected') ) {
			  		return false;
				}
			var $optionSet = $this.parents('.option-set');
			$optionSet.find('.selected').removeClass('selected');
			$this.addClass('selected');

			// make option object dynamically, i.e. { filter: '.my-filter-class' }
			var options = {},
				key = $optionSet.attr('data-option-key'),
				value = $this.attr('data-option-value');
			// parse 'false' as false boolean
			value = value === 'false' ? false : value;
			options[ key ] = value;
			if ( key === 'layoutMode' && typeof changeLayoutMode === 'function' ) {
			  // changes in layout modes need extra logic
			  changeLayoutMode( $this, options )
			} else {
			  // otherwise, apply new options
			  $container.isotope( options );
			}

			return false;
			
		  });
		
	});
	
});

/* ------------------ Tabs ----------------- */

function AutomaticTabsWidhtReset() {
	
	$('.nav-tabs').find('li').each(function(){

		$(this).find('a').css('width', 'auto');

	});
		
}

function AutomaticTabsWidht() {
		
	var containerWidth = $('.tab-content').width();
	
	var tabsNavWidth = 0;
	
	$('.nav-tabs').find('li').each(function(){
	
		var liWidth = $(this).outerWidth();
		
		tabsNavWidth = tabsNavWidth + liWidth;
		
	});
		
		
	if(tabsNavWidth > containerWidth) {
		
		var elements = $(".nav-tabs li").size();
		
		var newWidth = containerWidth / elements - 15;
		
		$('.nav-tabs').find('li').each(function(){

			$(this).find('a').css('width', newWidth);

		});	
		
	}
	
}


jQuery(document).ready(function($) {
	
	$('#myTab a').click(function (e) {
	  e.preventDefault();
	  $(this).tab('show');
	})
	
});

/* -------------------- Width Functions --------------------- */

jQuery(document).ready(function($){
	
	widthFunctions();

});


$(window).bind("resize", widthFunctions);

function widthFunctions(e) {
	
	AutomaticTabsWidhtReset();
	AutomaticTabsWidht();
	
    var winHeight = $(window).height();
    var winWidth = $(window).width();
    
	if (winWidth < 980 && winWidth > 767) {
		
		if($("#wall").width()) {
			
			if($(".item").hasClass("span3")) {

				$(".item").removeClass("span3");
				$(".item").addClass("span4");

			}
			
		}
		
		if($(".lr-page").hasClass("span4 offset4")) {
			
			$(".lr-page").removeClass("span4 offset4");
			$(".lr-page").addClass("span6 offset3");
			
			$("#page-title").removeClass("span4 offset4");
			$("#page-title").addClass("span6 offset3");
		}
						
	} else {
		
		if($("#wall").width()) {
			
			if($(".item").hasClass("span4")) {

				$(".item").removeClass("span4");
				$(".item").addClass("span3");

			}
			
		}
		
		if($(".lr-page").hasClass("span6 offset3")) {
			
			$(".lr-page").removeClass("span6 offset3");
			$(".lr-page").addClass("span4 offset4");
			
			$("#page-title").removeClass("span6 offset3");
			$("#page-title").addClass("span4 offset4");
		}
			
	}
	
}