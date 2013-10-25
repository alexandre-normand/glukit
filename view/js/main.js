// Close the dropdown menu if item in menu is clicked
$(document).ready(function closeMenu () {
	$('.top-bar').click(function(){
      $(this).removeClass('expanded');
	});
})

function groupEvents(userEvents, resolutionInMinutes) {
	eventGroups = [];
	if (userEvents.length > 0) {
		firstEventOfGroup = userEvents[0];
		
		currentGroup = [];
		for (var i = 0; i < userEvents.length; i++) {	

			currentEvent = userEvents[i];
			console.log("looking at event " + currentEvent);
			if (eventWithinResolution(firstEventOfGroup, currentEvent, resolutionInMinutes)) {
				currentGroup.push(currentEvent);
			} else {
				eventGroups.push(currentGroup);
				currentGroup = [currentEvent];
				firstEventOfGroup = currentEvent;
			}
			previousEvent = currentEvent;
		}
	}	

	return eventGroups;
}

function eventWithinResolution(first, other, resolutionInMinutes) {
	groupSizeInMillis = resolutionInMinutes * 60;	
	if ((other.x - first.x) <= groupSizeInMillis) {		
		return true;
	} 
	else {
		return false;
	}
}