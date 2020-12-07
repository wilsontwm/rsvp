const secondsDuration = 1000;
const minutesDuration = secondsDuration * 60;
const hoursDuration = minutesDuration * 60;
const daysDuration = hoursDuration * 24;

const weddingDate = new Date('6/19/2021 19:00');
const today = new Date();

var currentRunningNumber = 1;
$(document).ready(function(){
    var $addAttendeeLink = document.getElementById("btn-add-attendee");

    // Click via link in statement
    if (typeof($addAttendeeLink) != 'undefined' && $addAttendeeLink != null) {
        $addAttendeeLink.onclick = function() {
            addAttendee();
            
            // Hide the statement
            $('#text-add-attendee').hide();
        };
    }
    

    $("#btn-submit-rsvp").click(function(e){
        e.preventDefault();
        var queries = $('#form-rsvp').serialize();
        toggleLoading();
        $('#btn-submit-rsvp').html("Loading...");
        $('#btn-submit-rsvp').attr("disabled", true);

        axios.post('/rsvp', queries)
        .then(function (response) {
            $('#btn-submit-rsvp').html("Submit");
            $('#btn-submit-rsvp').attr("disabled", false);
            grecaptcha.reset();

            if(response['data']['success']) {      
                $('#form-rsvp')[0].reset();          
                $('#rsvp-modal').modal('toggle');
                Swal.fire({
                    title: 'Wilson & Shu Zhen:',
                    text: response['data']['message'],
                    icon: 'success',
                    confirmButtonText: 'OK'
                })
            } else {
                Swal.fire({
                    title: 'Ops!',
                    text: response['data']['message'],
                    icon: 'error',
                    confirmButtonText: 'OK'
                });
            }
        
        })
        .catch(function (error) {      
            $('#btn-submit-rsvp').html("Submit");
            $('#btn-submit-rsvp').attr("disabled", false);
            grecaptcha.reset();
            Swal.fire({
                title: 'Error!',
                text: "Please make sure that you have filled up all the required fields.",
                icon: 'error',
                confirmButtonText: 'OK'
            });
        });
    });

    // Smooth scrolling using jQuery easing
    $('a.js-scroll-trigger[href*="#"]:not([href="#"])').click(function() {
        if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
            var target = $(this.hash);
            target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
            if (target.length) {
            $('html, body').animate({
                scrollTop: (target.offset().top - 54)
            }, 1000, "easeInOutExpo");
            return false;
            }
        }
    });

    // Closes responsive menu when a scroll trigger link is clicked
    $('.js-scroll-trigger').click(function() {
        $('.navbar-collapse').collapse('hide');
    });

    // Activate scrollspy to add active class to navbar items on scroll
    $('body').scrollspy({
        target: '#mainNav',
        offset: 56
    });

    // Collapse Navbar
    var navbarCollapse = function() {
        if ($("#mainNav").offset().top > 100) {
            $("#mainNav").addClass("navbar-shrink");
        } else {
            $("#mainNav").removeClass("navbar-shrink");
        }
    };

    // Collapse now if page is not at top
    navbarCollapse();
    // Collapse the navbar when page is scrolled
    $(window).scroll(navbarCollapse);
    
    // Setup the countdown timer
    Countdown.init();

    // Image gallery
    // Gallery image hover
    $( ".img-wrapper" ).hover(
        function() {
            $(this).find(".img-overlay").animate({opacity: 1}, 600);
            //test.animate({opacity: 1}, 600);
            //console.log(test);
        }, function() {
            $(this).find(".img-overlay").animate({opacity: 0}, 600);
        }
    );

    // Lightbox
    var $overlay = $('<div id="overlay"></div>');
    var $image = $("<img>");
    var $prevButton = $('<div id="prevButton"><i class="fa fa-angle-left"></i></div>');
    var $nextButton = $('<div id="nextButton"><i class="fa fa-angle-right"></i></div>');
    var $exitButton = $('<div id="exitButton"><i class="fa fa-times"></i></div>');

    // Add overlay
    $overlay.append($image).prepend($prevButton).append($nextButton).append($exitButton);
    $("#gallery").append($overlay);

    // Hide overlay on default
    $overlay.hide();
  
    // When an image is clicked
    $(".img-overlay").click(function(event) {
        // Prevents default behavior
        event.preventDefault();
        // Adds href attribute to variable
        var imageLocation = $(this).prev().attr("href") != undefined && $(this).prev().attr("href") != '' ? $(this).prev().attr("href") : $(this).prev().attr("src");
        // Add the image src to $image
        $image.attr("src", imageLocation);
        $image.css("max-height", $(window).height());
        // Fade in the overlay
        $overlay.fadeIn("slow");
    });

    // When the overlay is clicked
    $overlay.click(function() {
        // Fade out the overlay
        $(this).fadeOut("slow");
    });

    // On press of keyboard key
    $("body").keyup(function(e) {
        if(e.keyCode == 37) { // left
            $prevButton.click();
        }
        else if(e.keyCode == 39) { // right
            $nextButton.click();
        }
        else if(e.keyCode == 27) { // escape
            $exitButton.click();
        }
    });
    
  
    // When next button is clicked
    $nextButton.click(function(event) {
        // Hide the current image
        $("#overlay img").hide();
        // Overlay image location
        var $currentImgSrc = $("#overlay img").attr("src");
        // Image with matching location of the overlay image
        var $currentImg = $('#image-gallery img[src="' + $currentImgSrc + '"]');
        // Finds the next image
        var $nextImg = $($currentImg.closest(".image").next().find("img"));
        // All of the images in the gallery
        var $images = $("#image-gallery img");
        // If there is a next image
        if ($nextImg.length > 0) { 
            // Fade in the next image
            $("#overlay img").attr("src", $nextImg.attr("src")).fadeIn(800);
        } else {
            // Otherwise fade in the first image
            $("#overlay img").attr("src", $($images[0]).attr("src")).fadeIn(800);
        }
        // Prevents overlay from being hidden
        event.stopPropagation();
    });

    // When previous button is clicked
    $prevButton.click(function(event) {
        // Hide the current image
        $("#overlay img").hide();
        // Overlay image location
        var $currentImgSrc = $("#overlay img").attr("src");
        // Image with matching location of the overlay image
        var $currentImg = $('#image-gallery img[src="' + $currentImgSrc + '"]');
        // Finds the next image
        var $nextImg = $($currentImg.closest(".image").prev().find("img"));
        // Fade in the next image
        $("#overlay img").attr("src", $nextImg.attr("src")).fadeIn(800);
        // Prevents overlay from being hidden
        event.stopPropagation();
    });

    // When the exit button is clicked
    $exitButton.click(function() {
        // Fade out the overlay
        $("#overlay").fadeOut("slow");
    });
});

// Create Countdown
var Countdown = {
  
    // Backbone-like structure
    $el: $('.countdown'),
    
    // Params
    countdown_interval: null,
    total_seconds     : 0,
    
    // Initialize the countdown  
    init: function() {
        const diffTime = (weddingDate - today) >= 0 ? (weddingDate - today) : 0;
        const diffDays = Math.floor(diffTime / daysDuration); 
        const diffHours = Math.floor((diffTime % daysDuration) / hoursDuration); 
        const diffMinutes = Math.floor(((diffTime % daysDuration) % hoursDuration) / minutesDuration);
        const diffSeconds = Math.round((((diffTime % daysDuration) % hoursDuration) % minutesDuration) / secondsDuration);
        
        // DOM
        this.$ = {
            days   : this.$el.find('.bloc-time.days .figure'),
            hours  : this.$el.find('.bloc-time.hours .figure'),
            minutes: this.$el.find('.bloc-time.min .figure'),
            seconds   : this.$el.find('.bloc-time.sec .figure'),
        };
  
        // Init countdown values
        this.values = {            
            days    : diffDays,
            hours   : diffHours,
            minutes : diffMinutes,
            seconds : diffSeconds
        };
        
        // Initialize total seconds
        this.total_seconds = this.values.days * 24 * 60 * 60 + (this.values.hours * 60 * 60) + this.values.minutes * 60 + this.values.seconds;

        // Animate countdown to the end 
        this.count();  
    },
    
    count: function() {
      
      var that    = this,
          $day_1  = this.$.days.eq(0),
          $day_2  = this.$.days.eq(1);
          $day_3  = this.$.days.eq(2),
          $hour_1 = this.$.hours.eq(0),
          $hour_2 = this.$.hours.eq(1),
          $min_1  = this.$.minutes.eq(0),
          $min_2  = this.$.minutes.eq(1),
          $sec_1  = this.$.seconds.eq(0),
          $sec_2  = this.$.seconds.eq(1),
      
          this.countdown_interval = setInterval(function() {
          if(that.total_seconds > 0) {
  
            --that.values.seconds;              

            if(that.values.minutes >= 0 && that.values.seconds < 0) {

                that.values.seconds = 59;
                --that.values.minutes;
            }            
  
            if(that.values.hours >= 0 && that.values.minutes < 0) {

                that.values.minutes = 59;
                --that.values.hours;
            }

            if(that.values.days >= 0 && that.values.hours < 0) {

                that.values.hours = 23;
                --that.values.days;
            }
  
            // Update DOM values  
            // Days
            that.checkDay(that.values.days, $day_1, $day_2, $day_3);

            // Hours
            that.checkHour(that.values.hours, $hour_1, $hour_2);

            // Minutes
            that.checkHour(that.values.minutes, $min_1, $min_2);
            
            // Secibds
            that.checkHour(that.values.seconds, $sec_1, $sec_2);

            --that.total_seconds;
          }
          else {
            clearInterval(that.countdown_interval);
          }
      }, 1000);    
    },
    
    animateFigure: function($el, value) {
      
       var that         = this,
           $top         = $el.find('.top'),
           $bottom      = $el.find('.bottom'),
           $back_top    = $el.find('.top-back'),
           $back_bottom = $el.find('.bottom-back');
  
      // Before we begin, change the back value
      $back_top.find('span').html(value);
  
      // Also change the back bottom value
      $back_bottom.find('span').html(value);
  
      // Then animate
      TweenMax.to($top, 0.8, {
          rotationX           : '-180deg',
          transformPerspective: 300,
            ease                : Quart.easeOut,
          onComplete          : function() {
  
              $top.html(value);
  
              $bottom.html(value);
  
              TweenMax.set($top, { rotationX: 0 });
          }
      });
  
      TweenMax.to($back_top, 0.8, { 
          rotationX           : 0,
          transformPerspective: 300,
            ease                : Quart.easeOut, 
          clearProps          : 'all' 
      });    
    },
    
    setNumber: function(value, $el) {
        var fig_value   = $el.find('.top').html(),
            val         = value.toString().charAt(0);

        if(fig_value !== val) this.animateFigure($el, val);
    },

    checkDay: function(value, $el_1, $el_2, $el_3) {
              
        var hundredth   = Math.floor(value / 100) % 10,
            tenth       = Math.floor(value / 10) % 10,
            oneth       = value % 10;
            
        this.setNumber(hundredth, $el_1);
        this.setNumber(tenth, $el_2);
        this.setNumber(oneth, $el_3);
        
    },

    checkHour: function(value, $el_1, $el_2) {
      
        var tenth       = Math.floor(value / 10) % 10,
            oneth       = value % 10;
            
        this.setNumber(tenth, $el_1);
        this.setNumber(oneth, $el_2);
    }
};

function addAttendee() {
    var html = getAddAttendeeTemplate(currentRunningNumber);
    var $addAttendeeContainer = document.getElementById("container-add-attendee");
    $addAttendeeContainer.insertAdjacentHTML('beforeend', html);
    currentRunningNumber++;
}

function removeAttendee(e) {
    // Get the parent element
    var $parent = e.closest(".attendee-row");
    if ($parent) $parent.parentNode.removeChild($parent);
    currentRunningNumber--;
    recalculateAttendeeNumber();
}

function recalculateAttendeeNumber() {
    var numbers = document.getElementsByClassName('attendee-number');

    for(var i = 1; i <= numbers.length; i++) {
        numbers[i-1].innerHTML = i;
    }

    if(numbers.length == 0) {
        $('#text-add-attendee').show();
    }
}

function getAddAttendeeTemplate(i) {
    var html =  '<div class="attendee-row">' +
                '<small class="text-muted col-xs-12">Attendee #<span class="attendee-number">' + i + '</span></small>' +
                '<div class="form-row">' +
                '<div class="form-group col-md-4 col-xs-6"><input name="names" type="text" class="form-control" id="inputName" placeholder="Name"></div>' + 
                '<div class="form-group col-md-4 col-xs-6"><input name="emails" type="email" class="form-control" id="inputEmail" placeholder="Email (Optional)" ></div>' + 
                '<div class="form-group col-md-3 col-xs-6"><input name="phones" type="text" class="form-control" id="inputPhone" placeholder="Phone (Optional)" ></div>' + 
                '<div class="col-md-1 col-xs-6">' +
                '<a href="#" class="btn btn-icon btn-add-attendee" onclick="addAttendee();"><i class="fa fa-plus text-success"></i></a>' + 
                '<a href="#" class="btn btn-icon btn-remove-attendee" onclick="removeAttendee(this);"><i class="fa fa-trash text-danger"></i></a>' +
                '</div>'
                '</div></div>';

    return html;
                    
}