$(document).ready(function(){
    
    $(".btn-delete-attendee").click(function(){
        var url = "/dashboard/attendees/" + $(this).data("id") + "/delete"
        $("#form-delete-attendee").attr("action", url); //Will set it
        toggleLoading();
        // Submit the form
        $("#form-delete-attendee").submit();
    });

});
