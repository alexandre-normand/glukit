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
	// Round the upper snap limit to the nearest hour
	var upperTimestamp = new Date(upperTimestamp * 1000);
	upperTimestamp.setMinutes(upperTimestamp.getMinutes() + 30);
	upperTimestamp.setMinutes(0);
	snapGuides = [upperTimestamp];

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

function splitReadsInRangeSegments(glucoseReads) {
  var segments = [];
  if (glucoseReads.length > 0) {
    previousRange = getRange(glucoseReads[0].y);
    previousRead = glucoseReads[0];
    var reads = [];
    for (var i = 0; i < glucoseReads.length; i++) {
      reads.push(previousRead);
      currentRead = glucoseReads[i];
      range = getRange(currentRead.y);

      if (range != previousRange) {      	
      	// We could interpolate a read directly on the boundary 
      	// (we know the y value, we need to interpolate x) instead of duplicating
      	// a read to have a continuous line that crosses the boundaries
        reads.push(currentRead);
        segment = new Object();
        segment.reads = reads;
        segment.range = previousRange;
        segments.push(segment);
        reads = [];
      }

      previousRange = range;
      previousRead = currentRead;
    }
    if (reads.length > 0) {
      segment = new Object();
      segment.reads = reads;
      segment.range = previousRange;
      segments.push(segment);
    }
  }

  return segments;
}

function getRange(glucoseValue) {
  if (glucoseValue > TARGET_RANGE_UPPER_BOUND) {
    return RANGES.HIGH;
  } else if (glucoseValue <= TARGET_RANGE_UPPER_BOUND && glucoseValue >= TARGET_RANGE_LOWER_BOUND) {
    return RANGES.NORMAL;
  } else if (glucoseValue < TARGET_RANGE_LOWER_BOUND) {
    return RANGES.LOW;
  }
}    

/**
 * Performs a binary search on the host array. This method can either be
 * injected into Array.prototype or called with a specified scope like this:
 * binaryIndexOf.call(someArray, searchElement);
 *
 * @param {*} searchElement The item to search for within the array.
 * @return {Number} The index of the element which defaults to -1 when not found.
 */
function binaryIndexOf(searchElement) {
	'use strict';
	var timestamp = searchElement.getTime() / 1000;
	var minIndex = 0;
	var maxIndex = this.length - 1;
	var currentIndex;
	var currentElement;
	var resultIndex;

	while (minIndex <= maxIndex) {
		resultIndex = currentIndex = (minIndex + maxIndex) / 2 | 0;
		currentElement = this[currentIndex];		        
        
		if (currentElement.x < timestamp) {
			minIndex = currentIndex + 1;
		}
		else if (currentElement.x > timestamp) {
			maxIndex = currentIndex - 1;
		}
		else {
			return currentIndex;
		}
	}

	return ~maxIndex;
}

function millisToDate(timestamp) {
	return new Date(timestamp);
}

function getHoverCoordinates(glucoseReads, time) {
	var glucoseIndex = Math.abs(binaryIndexOf.call(glucoseReads, time));
    var coordinates = new Object();

	if (glucoseIndex == glucoseReads.length - 1) {
		read = glucoseReads[glucoseReads.length - 1];
		coordinates.x = reads.x;
		coordinates.y = reads.y;		
	} else {
		coordinates.y = interpolateGlucoseRead(glucoseReads[glucoseIndex], glucoseReads[glucoseIndex + 1], time);
		coordinates.x = time;
	}

	return coordinates;
}

function interpolateGlucoseRead(left, right, time) {
	timestamp = time.getTime() / 1000;

	// gap between the two reads
	gap = right.x - left.x;

    // Calculate interpolation weights for left/right values
	leftWeight = (timestamp - left.x) / gap;
	rightWeight = 1 - leftWeight;

	return left.y * leftWeight + right.y * rightWeight;
}