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
			// Add indicator to the group that the type of event is present in it
			currentGroup[currentEvent.tag] = true;

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

function getDateSnapGuides(upperTimestamp, intervalInSeconds, numberOfSnaps) {
	snapGuides = [new Date(upperTimestamp * 1000)];
	currentSnapGuide = upperTimestamp;
	for (var i = 0; i < numberOfSnaps; i++) {
		currentSnapGuide = currentSnapGuide - intervalInSeconds;
		snapGuides.push(new Date(currentSnapGuide * 1000)); 
    }

    return snapGuides;
}

function getDayRangeFromUpperBound(upperTimestampInSeconds) {	
    lowerTimestamp = upperTimestampInSeconds - 86400;
    return [new Date(lowerTimestamp * 1000), new Date(upperTimestampInSeconds * 1000)];
}

function radiansFromDegree(degree) {
	return degree * 3.1416 / 180;
}

function generateEventMarkers(userEventGroup) {
	var parts = []
    if (userEventGroup.hasOwnProperty("Insulin") && userEventGroup.hasOwnProperty("Carbs")) {
      left = new Object();
      left.type = "left";
      left.tag = "Carbs";
      left.date = userEventGroup[0].date;
      left.x = userEventGroup[0].x;
      left.y = userEventGroup[0].y;
      
      parts.push(left);

      right = new Object();
      right.type = "right";
      right.tag = "Insulin";
      right.date = userEventGroup[0].date;
      right.x = userEventGroup[0].x;
      right.y = userEventGroup[0].y;

      parts.push(right);
    } else {
      whole = new Object();
      whole.type = "full";
      whole.tag = userEventGroup[0].tag;
      whole.date = userEventGroup[0].date;
      whole.x = userEventGroup[0].x;
      whole.y = userEventGroup[0].y;

      parts.push(whole);
    }
    
    return parts;
}