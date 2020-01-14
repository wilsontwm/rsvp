var currentRunningNumber = 1;
$(document).ready(function(){
    var $addAttendeeLink = document.getElementById("btn-add-attendee");

    // Click via link in statement
    $addAttendeeLink.onclick = function() {
        addAttendee();
        
        // Hide the statement
        $('#text-add-attendee').hide();
    };

    $("#btn-submit-rsvp").click(function(e){
        e.preventDefault();
        var queries = $('#form-rsvp').serialize();
        toggleLoading();
        
        axios.post('/rsvp', queries)
        .then(function (response) {
            console.log(response);
            toggleLoading();
            if(response['data']['success']) {                
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
            toggleLoading();
            Swal.fire({
                title: 'Error!',
                text: "Please make sure that you have filled up all the required fields.",
                icon: 'error',
                confirmButtonText: 'OK'
            });
        });
    });
});

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