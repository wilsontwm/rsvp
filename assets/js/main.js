var loadingOverlay;
$(document).ready(function(){
    loadingOverlay = document.querySelector('.loading');

    // Reset the loading when press on back button
    window.onpageshow = function(event) {
        if (event.persisted) {
            resetLoading();
        }
    };

    // Display the flash message
    window.Flash.create('.flash-message');

    // Whenever submit a form, toggle the loading     
    $('form').submit(function(){
        toggleLoading();
    });
});

// Show/hide the loading screen
function toggleLoading(){    
    if(loadingOverlay) {        
        document.activeElement.blur();
        if (loadingOverlay.classList.contains('hidden')){
            loadingOverlay.classList.remove('hidden');
        } else {
            loadingOverlay.classList.add('hidden');
        }
    }
    
}

// Reset the loading screen when page load
function resetLoading(){
    if (!loadingOverlay.classList.contains('hidden')){
        loadingOverlay.classList.add('hidden');
    }
}