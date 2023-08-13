/**
 * Table of contents
 *
 * 1. Components
 * 2. Header
 * 3. Core
 * 4. Elements
 * 5. Other
 * 6. Plugins
 */

(function($){ "use strict";
$(document).ready( function() {

/* -----------------------------------------------------------------------------

	1. COMPONENTS

----------------------------------------------------------------------------- */

	/* -------------------------------------------------------------------------
		QUANTITY FIELD
	------------------------------------------------------------------------- */

	$( '.quantity-field' ).each(function() {

		var $this = $(this),
			$input = $this.find( '.quantity-field__input' ),
			$addBtn = $this.find( '.quantity-field__btn--add' ),
			$removeBtn = $this.find( '.quantity-field__btn--remove' );

		// Get numeric value
		var getValue = function() {
			return ! isNaN( parseInt( $input.val() ) ) && parseInt( $input.val() ) > 0 ? parseInt( $input.val() ) : 0;
		};

		// Add
		$addBtn.on( 'click', function() {
			$input.val( getValue() + 1 );
		});

		// Remove
		$removeBtn.on( 'click', function() {

			if ( getValue() > 0 ) {
				$input.val( getValue() - 1 );
			}

		});

	});

/* -----------------------------------------------------------------------------

	2. HEADER

----------------------------------------------------------------------------- */

	/* -------------------------------------------------------------------------
		HEADER LANGUAGES
	------------------------------------------------------------------------- */

	$( '.header-languages' ).each(function() {

		var $this = $(this),
			$toggle = $this.find( '.header-languages__toggle' ),
			$list = $this.find( '.header-languages__list' );

		// Toggle list
		$toggle.on( 'click', function() {

			// Toggle on desktop
			if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() > 1199 ) {

				$list.fadeToggle( 200, function() {

					// Bind closeList event
					$(window).on( 'click', closeList );
					$( '.header-search__toggle' ).on( 'click', closeList );

				});

			}

			// Toggle on mobile
			else {

				$list.slideToggle( 200 );
				$( '#header' ).toggleClass( 'header--languages-active' );

			}

			$toggle.toggleClass( 'header-languages__toggle--active' );

		});

		// Hide on click outside
		var closeList = function() {

       		$toggle.removeClass( 'header-languages__toggle--active' );
			$list.fadeOut( 200 );

			// Unbind closeList event
			$(window).unbind( 'click', closeList );
			$( '.header-search__toggle' ).unbind( 'click', closeList );

		};
		$this.on( 'click', function( e ) {
    		e.stopPropagation();
		});

		// Reset on screen transition
		$(document).on( 'lsvrBeautyspotScreenTransition', function() {

			$toggle.removeClass( 'header-languages__toggle--active' );
			$( '#header' ).removeClass( 'header--languages-active' );
			$list.removeAttr( 'style' );
			$(window).unbind( 'click', closeList );
			$( '.header-search__toggle' ).unbind( 'click', closeList );

		});

	});

	/* -------------------------------------------------------------------------
		HEADER MENU
	------------------------------------------------------------------------- */

	$( '.header-menu' ).each(function() {

		var $this = $(this);

		// Reset menu function
		function resetMenu() {
			$this.find( '.header-menu__item' ).removeClass( 'header-menu__item--hover header-menu__item--active' );
			$this.find( '.header-menu__submenu' ).removeAttr( 'style' );
			$this.find( '.header-menu__submenu-toggle' ).removeClass( 'header-menu__submenu-toggle--active' );
		}

		// Reset menu when click on link without submenu
		$this.find( '.header-menu__item-link' ).each(function() {
			$(this).on( 'click', function() {
				if ( $(this).parent().parent().find( '> .header-menu__submenu' ).length < 1 ) {
					resetMenu();
				}
			});
		});

		// Parse submenus
		$this.find( '.header-menu__submenu' ).each(function() {

			var $submenu = $(this),
				$parent = $submenu.parent(),
				$toggle = $parent.find( '> .header-menu__submenu-toggle' ),
				$linkWrapper = $parent.find( '> .header-menu__item-link-wrapper' ),
				$link = $linkWrapper.find( '> .header-menu__item-link' );

			// Show desktop submenu function
			function desktopShowSubmenu() {
				$submenu.show();
				$parent.addClass( 'header-menu__item--hover' );
			}

			// Hide desktop submenu function
			function desktopHideSubmenu() {
				$submenu.hide();
				$parent.removeClass( 'header-menu__item--hover' );
			}

			// Show mobile submenu function
			function mobileShowSubmenu() {
				$submenu.slideDown( 150 );
				$parent.addClass( 'header-menu__item--active' );
			}

			// Hide mobile submenu function
			function mobileHideSubmenu() {
				$submenu.slideUp( 150 );
				$parent.removeClass( 'header-menu__item--active' );
			}

			// Desktop interaction
			if ( $submenu.parent().parent().hasClass( 'header-menu__list' ) ) {

				// Desktop mouseover and focus action
				$parent.on( 'mouseover focus', function() {
					if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() > 1199 ) {
						desktopShowSubmenu();
					}
				});

				// Desktop mouseleave and blur action
				$parent.on( 'mouseleave blur', function() {
					if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() > 1199 ) {
						desktopHideSubmenu();
					}
				});

				// Desktop click or key enter
				$link.on( 'click', function() {

					if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() > 1199 && ! $parent.hasClass( 'header-menu__item--hover' ) ) {

						// Hide opened submenus
						$parent.siblings( '.header-menu__item.header-menu__item--hover' ).each(function() {

							$(this).removeClass( 'header-menu__item--hover' );
							$(this).find( '> .header-menu__submenu' ).hide();

							if ( $(this).hasClass( 'header-menu__item--dropdown' ) ) {

								$(this).find( '.header-menu__item--hover' ).removeClass( 'header-menu__item--hover' );
								$(this).find( '.header-menu__submenu' ).hide();

							}

						});

						// Show submenu
						desktopShowSubmenu();

						// Hide on click outside
						$( 'html' ).on( 'click.lsvrBeautyspotHeaderMenuCloseSubmenuOnClickOutside touchstart.lsvrBeautyspotHeaderMenuCloseSubmenuOnClickOutside', function(e) {

							desktopHideSubmenu();
							$( 'html' ).unbind( 'click.lsvrBeautyspotHeaderMenuCloseSubmenuOnClickOutside touchstart.lsvrBeautyspotHeaderMenuCloseSubmenuOnClickOutside' );

						});

						// Disable link
						$parent.on( 'click touchstart', function(e) {
							e.stopPropagation();
						});
						return false;

					} else {
						resetMenu();
					}

				});

			}

			// Mobile interactions
			$toggle.on( 'click', function() {

				$toggle.toggleClass( 'header-menu__submenu-toggle--active' );
				if ( $toggle.hasClass( 'header-menu__submenu-toggle--active' ) ) {
					mobileShowSubmenu();
				} else {
					mobileHideSubmenu();
				}

			});

		});

		// Reset on screen transition
		$(document).on( 'lsvrBeautyspotScreenTransition', function() {
			resetMenu();
		});

	});

	/* -------------------------------------------------------------------------
		HEADER SEARCH
	------------------------------------------------------------------------- */

	$( '.header-search' ).each(function() {

		var $this = $(this),
			$toggle = $this.find( '.header-search__toggle' ),
			$form = $this.find( '.header-search__form' );

		// Toggle form
		$toggle.on( 'click', function() {

			$form.fadeToggle( 200, function() {

				// Input focus
				if ( $(this).is( ':visible' ) ) {
					$(this).find( 'input' ).focus();
				}

				// Bind closeForm event
				$(window).on( 'click', closeForm );
				$( '.header-languages__toggle' ).on( 'click', closeForm );

			});
			$toggle.toggleClass( 'header-search__toggle--active' );

		});

		// Hide on click outside
		var closeForm = function() {

       		$toggle.removeClass( 'header-search__toggle--active' );
			$form.fadeOut( 200 );

			// Unbind closeForm event
			$(window).unbind( 'click', closeForm );
			$( '.header-languages__toggle' ).unbind( 'click', closeForm );

		};
		$this.on( 'click', function( e ) {
    		e.stopPropagation();
		});

		// Reset on screen transition
		$(document).on( 'lsvrBeautyspotScreenTransition', function() {

			$toggle.removeClass( 'header-search__toggle--active' );
			$form.removeAttr( 'style' );
			$(window).unbind( 'click', closeForm );
			$( '.header-languages__toggle' ).unbind( 'click', closeForm );

		});

	});

	/* -------------------------------------------------------------------------
		HEADER PANEL
	------------------------------------------------------------------------- */

	// Check for panel collision
	var panelOverlap = function() {

		var $headerInner = $( '#header .header__inner' ),
			$branding = $( '#header .header-branding' ),
			$menu = $( '#header .header-menu' ),
			$panel = $( '#header .header-panel' ),
			panelGap = 60;

		// Check for overlap
		if ( ( $branding.length > 0 || $menu.length > 0 ) && $panel.length > 0 ) {

			// Has both branding and menu
			if ( ( $branding.length > 0 && $menu.length > 0 ) && ( $headerInner.height() - ( $branding.height() + $menu.height() ) < $panel.height() + panelGap ) ) {
				return true;
			}

			// Has only branding
			else if ( $branding.length > 0 && ( $headerInner.height() - $branding.height() < $panel.height() + panelGap ) ) {
				return true;
			}

			// Has only menu
			else if ( $menu.length > 0 && ( $headerInner.height() - $menu.height() < $panel.height() + panelGap ) ) {
				return true;
			}

		}

		return false;

	};

	// Toggle
	$( '#header.header--has-contact.header--has-collision-detection' ).each(function() {

		var $header = $(this),
			$panelInner = $(this).find( '.header-panel__inner' ),
			$toggle = $(this).find( '.header-panel__toggle' ),
			$cta = $(this).find( '.header-cta' ),
			$contact = $(this).find( '.header-contact' ),
			headerDefaultClasses;

		// Save default header classes
		if ( $header.hasClass( 'header--has-collapsed-panel' ) ) {
			headerDefaultClasses = 'header--has-collapsed-panel';
		} else {
			headerDefaultClasses = 'header--has-expanded-panel';
		}
		if ( $header.hasClass( 'header--has-collapsible-panel' ) ) {
			headerDefaultClasses += ' header--has-collapsible-panel';
		}

		// Toggle
		$toggle.on( 'click', function() {

			// Desktop toggle
			if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() > 1199 ) {

				$cta.slideToggle( 200 );
				$contact.slideToggle( 200, function() {

					$header.toggleClass( 'header--has-collapsed-panel header--has-expanded-panel' );

					// Check for collision on open
					if ( $contact.is( ':visible' ) && true === panelOverlap() ) {
						$header.addClass( 'header--has-panel-collision' );
					}

				});

			}

			// Mobile toggle
			else {

				$panelInner.slideToggle( 300 );

			}

		});

		// Function for checking for the collision
		var checkForCollision = function() {

			if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() > 1199 && true === panelOverlap() ) {
				$header.removeClass( 'header--has-expanded-panel' ).addClass( 'header--has-panel-collision header--has-collapsible-panel header--has-collapsed-panel' );
				$cta.hide();
				$contact.hide();
				$toggle.show();
			}

		};

		// Check for collision after initial load
		checkForCollision();

		// Check for collision on screen transition
		$(document).on( 'lsvrBeautyspotScreenTransition', function() {

			$header.removeClass( 'header--has-panel-collision header--has-collapsible-panel header--has-expanded-panel header--has-collapsed-panel' );
			$header.addClass( headerDefaultClasses );
			$panelInner.removeAttr( 'style' );
			$cta.removeAttr( 'style' );
			$contact.removeAttr( 'style' );
			$toggle.removeAttr( 'style' );
			checkForCollision();

		});

	});

	/* -------------------------------------------------------------------------
		MOBILE TOGGLE
	------------------------------------------------------------------------- */

	$( '.header-mobile-toggle' ).each(function() {

		var $toggle = $(this),
			$panel = $( '.header-panel' ),
			$menu = $( '.header-menu' ),
			$search = $( '.header-search' );

		// Toggle
		$toggle.on( 'click', function() {

			$search.slideToggle( 300 );
			$panel.slideToggle( 300 );
			$menu.slideToggle( 300 );
			$(this).toggleClass( 'header-mobile-toggle--active' );

		});

		// Reset on screen transition
		$(document).on( 'lsvrBeautyspotScreenTransition', function() {

			$search.removeAttr( 'style' );
			$panel.removeAttr( 'style' );
			$menu.removeAttr( 'style' );
			$toggle.removeClass( 'header-mobile-toggle--active' );

		});

	});


/* -----------------------------------------------------------------------------

	3. CORE

----------------------------------------------------------------------------- */

	/* -------------------------------------------------------------------------
		CONTACT FORM
	------------------------------------------------------------------------- */

	$( '.lsvr-form--contact' ).each(function() {

		var $form = $(this),
			$validationErrorMessage = $form.find( '.lsvr-form__message--validation-error' ),
			$connectionErrorMessage = $form.find( '.lsvr-form__message--connection-error' ),
			$successMessage = $form.find( '.lsvr-form__message--success' ),
			$submitBtn = $form.find( '.lsvr-form__submit' );

		// Submit
		$form.on( 'submit', function( e ) {

			// Check if not loading
			if ( ! $form.hasClass( 'lsvr-form--loading' ) ) {

				// Form is valid - make ajax request
				if ( true === $form.lsvrBeautyspotValidateForm() && $form.hasClass( 'lsvr-form--ajax' ) ) {

					e.preventDefault();

					$connectionErrorMessage.slideUp( 200 );
					$validationErrorMessage.slideUp( 200 );
					$successMessage.slideUp( 200 );
					$form.addClass( 'lsvr-form--loading' );
					$submitBtn.attr( 'data-label', $submitBtn.text() ).text( $submitBtn.data( 'loading-label' ) );

					// Make ajax request
					$.ajax({
						type: 'post',
	            		dataType: 'json',
						url: $form.attr( 'action' ),
						data: $form.serialize(),
						success: function( response ) {

							$form.removeClass( 'lsvr-form--loading' );
							$submitBtn.text( $submitBtn.attr( 'data-label' ) );

							// Connection error
							if ( response.hasOwnProperty( 'type' ) && 'connection-error' === response.type && response.hasOwnProperty( 'message' ) ) {

								$connectionErrorMessage.html( '<p>' + response.message + '</p>' );
								$connectionErrorMessage.slideDown( 200 );

							}

							// Validation error
							else if ( response.hasOwnProperty( 'type' ) && 'validation-error' === response.type && response.hasOwnProperty( 'message' ) ) {

								$validationErrorMessage.html( '<p>' + response.message + '</p>' );
								$validationErrorMessage.slideDown( 200 );

							}

							// Success
							else if ( response.hasOwnProperty( 'type' ) && 'success' === response.type && response.hasOwnProperty( 'message' ) ) {

								$successMessage.html( '<p>' + response.message + '</p>' );
								$successMessage.slideDown( 200 );

								// Reset all fields
								$form.find( 'input, textarea' ).val( '' );

							}

							// Unable to parse
							else {

								$form.removeClass( 'lsvr-form--loading' );
								$submitBtn.text( $submitBtn.attr( 'data-label' ) );
								$connectionErrorMessage.slideDown( 200 );

							}

						},
						error: function() {

							$form.removeClass( 'lsvr-form--loading' );
							$submitBtn.text( $submitBtn.attr( 'data-label' ) );
							$connectionErrorMessage.slideDown( 200 );

						}
					});

				}

				// Form is valid
				else if ( true === $form.lsvrBeautyspotValidateForm() ) {
					$successMessage.slideDown( 200 );
					$validationErrorMessage.slideUp( 200 );
				}

				// Form is invalid
				else {
					$successMessage.slideUp( 200 );
					$validationErrorMessage.slideDown( 200 );
					return false;
				}

			}

			// Form is loading
			else {
				return false;
			}

		});

	});


/* -----------------------------------------------------------------------------

	4. ELEMENTS

----------------------------------------------------------------------------- */

	/* -------------------------------------------------------------------------
		LSVR ACCORDION
	------------------------------------------------------------------------- */

	$( '.lsvr-accordion' ).each(function() {

		var $accordion = $(this),
			isClassic = $accordion.hasClass( 'lsvr-accordion--toggle' ) ? false : true;

		// Parse accordion items
		$accordion.find( '.lsvr-accordion__item' ).each(function() {

			var $this = $(this),
				$header = $this.find( '.lsvr-accordion__item-header' ),
				$content = $this.find( '.lsvr-accordion__item-content-wrapper' );

			// Toggle
			$header.on( 'click', function() {

				// Classic mode
				if ( true === isClassic ) {

					$accordion.find( '.lsvr-accordion__item--expanded' ).not( $this ).each(function() {
						$(this).find( '.lsvr-accordion__item-content-wrapper' ).slideUp( 200 );
						$(this).removeClass( 'lsvr-accordion__item--expanded' );
					});

				}

				$content.slideToggle( 200 );
				$this.toggleClass( 'lsvr-accordion__item--expanded' );

			});

		});

	});

	/* -------------------------------------------------------------------------
		LSVR FORM
	------------------------------------------------------------------------- */

	// Validation form
	$( '.lsvr-form:not( .lsvr-form--contact )' ).each(function() {

		var $form = $(this),
			$validationErrorMessage = $form.find( '.lsvr-form__message--validation-error' );

		// Submit
		$form.on( 'submit', function( e ) {

			// Form is valid
			if ( true === $form.lsvrBeautyspotValidateForm() ) {
				$validationErrorMessage.slideUp( 200 );
			}

			// Form is invalid
			else {
				$validationErrorMessage.slideDown( 200 );
				return false;
			}

		});

	});

	// Checkbox toggle
	$( '.lsvr-form__field-input--checkbox[data-toggle-element]' ).each(function() {

		var $this = $(this),
			elementID = $this.data( 'toggle-element' );

		// Toggle element
		$this.on( 'change', function() {
			$( '#' + elementID ).slideToggle( 200, function() {

				// Toggle disabled sate for all inputs
				if ( $(this).is( ':visible' ) ) {
					$(this).find( 'input, select, textarea' ).prop( 'disabled', false );
				} else {
					$(this).find( 'input, select, textarea' ).prop( 'disabled', true );
				}

				// Remove field validation error statuses
				$(this).find( '.lsvr-form__field-input--error' ).removeClass( 'lsvr-form__field-input--error' );

			});
		});

	});

	/* -------------------------------------------------------------------------
		LSVR GRID
	------------------------------------------------------------------------- */

	// Masonry
	if ( $.fn.masonry && $.fn.imagesLoaded ) {
		$( '.lsvr-grid--masonry' ).each(function() {

			var $this = $(this),
				originLeft = $( 'html' ).attr( 'dir' ) && 'rtl' === $( 'html' ).attr( 'dir' ) ? false : true;

			// Wait for images to load
			$this.imagesLoaded(function() {
				$this.masonry({
					originLeft : originLeft
				});
			});

		});
	}

	/* -------------------------------------------------------------------------
		LSVR SLIDE LIST
	------------------------------------------------------------------------- */

	$( '.lsvr-slide-list' ).each(function() {

		var $this = $(this),
			$slides = $this.find( '.lsvr-slide-list__item' ),
			$active = $slides.first(),
			autoplay = $this.attr( 'data-autoplay' ) ? parseInt( $this.attr( 'data-autoplay' ) ) : 0,
			$prevBtn = $this.find( '.lsvr-slide-list__nav-button--prev' ),
			$nextBtn = $this.find( '.lsvr-slide-list__nav-button--next' ),
			$next;

		// Set default active
		$active.addClass( 'lsvr-slide-list__item--active' );

		// Slide
		function slide( direction ) {

			var slideDelay = $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() > 768 ? 1000 : 0;

			$this.addClass( 'lsvr-slide-list--animation-in-progress' );

			// Previous
			if ( 'prev' === direction ) {

				if ( $active.prev().length > 0 ) {
					$next = $active.prev();
				} else {
					$next = $slides.last();
				}

			}

			// Next
			else {

				if ( $active.next().length > 0 ) {
					$next = $active.next();
				} else {
					$next = $slides.first();
				}

			}

			// Start animation
			$next.addClass( 'lsvr-slide-list__item--next lsvr-slide-list__item--animate-in' );
			$active.addClass( 'lsvr-slide-list__item--animate-out' );

			// Reset after animation ends
			setTimeout( function() {

				$next.removeClass( 'lsvr-slide-list__item--next lsvr-slide-list__item--animate-in' );
				$next.addClass( 'lsvr-slide-list__item--active' );
				$active.removeClass( 'lsvr-slide-list__item--active lsvr-slide-list__item--animate-out' );
				$active = $next;
				$this.removeClass( 'lsvr-slide-list--animation-in-progress' );

			}, slideDelay );

		}

		// Click on next
		$nextBtn.on( 'click', function() {
			if ( ! $this.hasClass( 'lsvr-slide-list--animation-in-progress' ) ) {
				slide( 'next' );
			}
		});

		// Click on prev
		$prevBtn.on( 'click', function() {
			if ( ! $this.hasClass( 'lsvr-slide-list--animation-in-progress' ) ) {
				slide( 'prev' );
			}
		});

		// Autoplay function
		var autoplayTimeout;
		function startAutoplay( delay ) {
			autoplayTimeout = setTimeout( function() {
				slide( 'next' );
				startAutoplay( delay );
			}, delay * 1000 );
		}

		// Set autoplay
		if ( autoplay > 0 ) {

			// Initial start
			startAutoplay( autoplay );

			// Pause on hover
			$this.on( 'mouseenter', function() {
				clearTimeout( autoplayTimeout );
			});

			// Resume oon leave
			$this.on( 'mouseleave', function() {
				startAutoplay( autoplay );
			});

		}

	});

	/* -------------------------------------------------------------------------
		LSVR TABS
	------------------------------------------------------------------------- */

	$( '.lsvr-tabs' ).each(function() {

		var $this = $(this),
			$headerItems = $this.find( '.lsvr-tabs__header-item' ),
			$contentItems = $this.find( '.lsvr-tabs__content-item' );

		$headerItems.each(function() {

			var $parent = $(this),
				$link = $(this).find( '.lsvr-tabs__header-item-link' ),
				index = $(this).index();

			// Switch tabs on click
			$link.on( 'click', function( e ) {

				if ( ! $parent.hasClass( '.lsvr-tabs__header-item--active' ) ) {

					e.preventDefault();

					// Set header item as active
					$headerItems.filter( '.lsvr-tabs__header-item--active' ).removeClass( 'lsvr-tabs__header-item--active' );
					$parent.addClass( 'lsvr-tabs__header-item--active' );

					// Show content item
					$contentItems.filter( '.lsvr-tabs__content-item--active' ).removeClass( 'lsvr-tabs__content-item--active' ).hide();
					$contentItems.filter( ':nth-child(' + ( index + 1 ) + ')' ).addClass( 'lsvr-tabs__content-item--active' ).show();

				}

			});

		});

	});


/* -----------------------------------------------------------------------------

	5. OTHER

----------------------------------------------------------------------------- */

	/* -------------------------------------------------------------------------
		MAGNIFIC POPUP
	------------------------------------------------------------------------- */

	if ( $.fn.magnificPopup ) {

		var js_strings = {
			'mp_tClose' : 'Close (Esc)',
			'mp_tLoading' : 'Loading...',
			'mp_tPrev' : 'Previous (Left arrow key)',
			'mp_tNext' : 'Next (Right arrow key)',
			'mp_image_tError' : 'The image could not be loaded.',
			'mp_ajax_tError' : 'The content could not be loaded.'
		};

		$.extend( true, $.magnificPopup.defaults, {
			tClose: js_strings.mp_tClose,
			tLoading: js_strings.mp_tLoading,
			gallery: {
				tPrev: js_strings.mp_tPrev,
				tNext: js_strings.mp_tNext,
				tCounter: '%curr% / %total%'
			},
			image: {
				tError: js_strings.mp_image_tError,
			},
			ajax: {
				tError: js_strings.mp_ajax_tError,
			}
		});

		// Init lightbox
		$( '.open-in-lightbox' ).magnificPopup({
			type: 'image',
			removalDelay: 300,
			mainClass: 'mfp-fade',
			gallery: {
				enabled: true
			}
		});

	}

});
})(jQuery);

(function($){ "use strict";

/* -----------------------------------------------------------------------------

	6. PLUGINS

----------------------------------------------------------------------------- */

	/* -------------------------------------------------------------------------
		FORM VALIDATION
	------------------------------------------------------------------------- */

	$.fn.lsvrBeautyspotValidateForm = function() {

		function emailValid( email ) {
			var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
			return re.test(email);
		}

		var $form = $(this),
			formValid = true;

		// Validate input
		var validateInput = function( input ) {

			var $input = input,
				value = $input.val(),
				placeholder = $input.data( 'placeholder' ) ? $input.data( 'placeholder' ) : false,
				inputValid = false;

			if ( value.trim() !== '' && ! ( placeholder && value === placeholder ) ) {

				// Email inputs
				if ( $input.hasClass( 'lsvr-form__field-input--email' ) ) {

					if ( ! emailValid( value ) ) {
						$input.addClass( 'lsvr-form__field-input--error' );
					}

					else {
						$input.removeClass( 'lsvr-form__field-input--error' );
						inputValid = true;
					}

				}

				// Select field
				else if ( $input.prop( 'tagName' ).toLowerCase() === 'select' ) {

					if ( value === null ) {
						$input.addClass( 'lsvr-form__field-input--error' );
					}

					else {
						$input.removeClass( 'lsvr-form__field-input--error' );
						inputValid = true;
					}

				}

				// Default field
				else {
					$input.removeClass( 'lsvr-form__field-input--error' );
					inputValid = true;
				}

			}
			else {
				$input.addClass( 'lsvr-form__field-input--error' );
			}

			return inputValid;

		};

		// Check required fields
		$form.find( '.lsvr-form__field-input--required:not([disabled])' ).each(function(){
			formValid = ! validateInput( $(this) ) ? false : formValid;
		});

		$form.find( '.lsvr-form__field-input--error' ).first().focus();

		return formValid;

	};

	/* -------------------------------------------------------------------------
		MEDIA QUERY BREAKPOINT
	------------------------------------------------------------------------- */

	if ( ! $.fn.lsvrBeautyspotGetMediaQueryBreakpoint ) {
		$.fn.lsvrBeautyspotGetMediaQueryBreakpoint = function() {

			if ( $( '#lsvr-media-query-breakpoint' ).length < 1 ) {
				$( 'body' ).append( '<span id="lsvr-media-query-breakpoint" style="display: none;"></span>' );
			}
			var value = $( '#lsvr-media-query-breakpoint' ).css( 'font-family' );
			if ( typeof value !== 'undefined' ) {
				value = value.replace( "\"", "" ).replace( "\"", "" ).replace( "\'", "" ).replace( "\'", "" );
			}
			if ( isNaN( value ) ) {
				return $( window ).width();
			}
			else {
				return parseInt( value );
			}

		};
	}

	var lsvrBeautyspotMediaQueryBreakpoint;
	if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint ) {
		lsvrBeautyspotMediaQueryBreakpoint = $.fn.lsvrBeautyspotGetMediaQueryBreakpoint();
		$(window).resize(function() {
			if ( $.fn.lsvrBeautyspotGetMediaQueryBreakpoint() !== lsvrBeautyspotMediaQueryBreakpoint ) {
				lsvrBeautyspotMediaQueryBreakpoint = $.fn.lsvrBeautyspotGetMediaQueryBreakpoint();
				$.event.trigger({
					type: 'lsvrBeautyspotScreenTransition',
					message: 'Screen transition completed.',
					time: new Date()
				});
			}
		});
	}
	else {
		lsvrBeautyspotMediaQueryBreakpoint = $(document).width();
	}

})(jQuery);