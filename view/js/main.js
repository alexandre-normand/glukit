function EventGroup() {
    this.x;
    this.includesInsulin = false;
    this.includesFood = false;
    this.userEvents = [];
    this.date;
}

EventGroup.prototype.addEvent = function(userEvent) {
    if (this.userEvents.length < 1) {
        this.date = userEvent.date;
        this.x = userEvent.x;
    }
    switch (userEvent.tag) {
        case "Insulin":
            this.includesInsulin = true;
            break;
        case "Carbs":
            this.includesFood = true;
            break;
    }
    this.userEvents.push(userEvent);
}

EventGroup.prototype.valueOf = function() {
    return this.time;
};

function groupEvents(userEvents, resolutionInMinutes) {
    eventGroups = [];
    if (userEvents.length > 0) {
        firstEventOfGroup = userEvents[0];

        currentGroup = new EventGroup();
        eventGroups.push(currentGroup);
        for (var i = 0; i < userEvents.length; i++) {
            currentEvent = userEvents[i];

            if (eventWithinResolution(firstEventOfGroup, currentEvent, resolutionInMinutes)) {
                currentGroup.addEvent(currentEvent);
            } else {
                currentGroup = new EventGroup();
                currentGroup.addEvent(currentEvent);
                eventGroups.push(currentGroup);
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
    } else {
        return false;
    }
}

function getDateSnapGuides(lowerTimestamp, upperTimestamp) {
    // Approximate first snap boundary to the oldest full hour	
    lowerTimestamp.setMinutes(lowerTimestamp.getMinutes() + 55);
    lowerTimestamp.setMinutes(0);
    snapGuides = [lowerTimestamp];

    currentSnapGuide = lowerTimestamp;
    while (currentSnapGuide < upperTimestamp) {
        currentSnapGuide = moment(currentSnapGuide).add('day', 1);
        snapGuides.push(currentSnapGuide.toDate());
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
    if (userEventGroup.includesFood && userEventGroup.includesInsulin) {
        left = new Object();
        left.type = "left";
        left.tag = "Carbs";
        left.date = userEventGroup.date;
        left.x = userEventGroup.userEvents[0].x;
        left.y = userEventGroup.userEvents[0].y;

        parts.push(left);

        right = new Object();
        right.type = "right";
        right.tag = "Insulin";
        right.date = userEventGroup.date;
        right.x = userEventGroup.userEvents[0].x;
        right.y = userEventGroup.userEvents[0].y;

        parts.push(right);
    } else {
        whole = new Object();
        whole.type = "full";
        whole.tag = userEventGroup.userEvents[0].tag;
        whole.date = userEventGroup.date;
        whole.x = userEventGroup.userEvents[0].x;
        whole.y = userEventGroup.userEvents[0].y;

        parts.push(whole);
    }

    return parts;
}

function splitReadsInRangeSegments(glucoseReads, unit) {
    var segments = [];
    if (glucoseReads.length > 0) {
        previousRange = getRange(glucoseReads[0].y, unit);
        previousRead = glucoseReads[0];
        var reads = [];
        for (var i = 0; i < glucoseReads.length; i++) {
            reads.push(previousRead);
            currentRead = glucoseReads[i];
            range = getRange(currentRead.y, unit);

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

// Find the last segment that ends after the given time
function getFirstSegmentEndingAfter(segments, time, startIndex) {
    for (var i = startIndex; i < segments.length; i++) {
        segmentReads = segments[i].reads;
        if (segmentReads[segmentReads.length - 1].date >= time) {
            return i;
        }
    }

    return segments.length - 1;
}

function RangeAggregate() {
    this.normalTimeInMinutes = 0;
    this.lowTimeInMinutes = 0;
    this.highTimeInMinutes = 0;
}

RangeAggregate.prototype.getTotalTime = function() {
    return this.normalTimeInMinutes + this.lowTimeInMinutes + this.highTimeInMinutes;
};

// Sums up the time in minutes spent in in each range
function getRangeAggregate(segments, lowerBound, upperBound) {
    firstSegmentIndex = getFirstSegmentEndingAfter(segments, lowerBound, 0);
    lastSegmentIndex = getFirstSegmentEndingAfter(segments, upperBound, firstSegmentIndex);

    rangeAggregate = new RangeAggregate();

    // Calculate the slice of the first segment that we should be considering
    segmentReads = segments[firstSegmentIndex].reads;
    endOfFirstSegment = moment.unix(segmentReads[segmentReads.length - 1].x);
    durationInMinutes = endOfFirstSegment.diff(moment(lowerBound), 'minutes');

    rangeAggregate = addRangeToAggregate(segments[firstSegmentIndex].range, durationInMinutes, rangeAggregate);

    for (var i = firstSegmentIndex + 1; i < lastSegmentIndex; i++) {
        segmentReads = segments[i].reads;
        endOfSegment = moment.unix(segmentReads[segmentReads.length - 1].x);
        durationInMinutes = endOfSegment.diff(moment.unix(segmentReads[0].x), 'minutes');

        rangeAggregate = addRangeToAggregate(segments[i].range, durationInMinutes, rangeAggregate);
    }

    // Calculate the slice of the last segment
    segmentReads = segments[lastSegmentIndex].reads;
    durationInMinutes = moment(upperBound).diff(moment.unix(segmentReads[0].x), 'minutes');

    rangeAggregate = addRangeToAggregate(segments[lastSegmentIndex].range, durationInMinutes, rangeAggregate);

    return rangeAggregate;
}

function addRangeToAggregate(range, durationInMinutes, aggregate) {
    switch (range) {
        case RANGES.NORMAL:
            aggregate.normalTimeInMinutes = aggregate.normalTimeInMinutes + durationInMinutes;
            break;
        case RANGES.HIGH:
            aggregate.highTimeInMinutes = aggregate.highTimeInMinutes + durationInMinutes;
            break;
        case RANGES.LOW:
            aggregate.lowTimeInMinutes = aggregate.lowTimeInMinutes + durationInMinutes;
            break;
    }
    return aggregate;
}

function getRange(glucoseValue, unit) {
    var targetRangeUpperValue = getUpperRangeValue(unit);
    var targetRangeLowerValue = getLowerRangeValue(unit);
    if (glucoseValue > targetRangeUpperValue) {
        return RANGES.HIGH;
    } else if (glucoseValue <= targetRangeUpperValue && glucoseValue >= targetRangeLowerValue) {
        return RANGES.NORMAL;
    } else if (glucoseValue < targetRangeLowerValue) {
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
        } else if (currentElement.x > timestamp) {
            maxIndex = currentIndex - 1;
        } else {
            return currentIndex;
        }
    }

    return~ maxIndex;
}

function millisToDate(timestamp) {
    return new Date(timestamp);
}

function getHoverCoordinates(glucoseReads, time) {
    var glucoseIndex = Math.abs(binaryIndexOf.call(glucoseReads, time));
    var coordinates = new Object();

    if (glucoseIndex >= glucoseReads.length - 1) {
        read = glucoseReads[glucoseReads.length - 1];
        coordinates.x = glucoseReads[glucoseReads.length - 1].x;
        coordinates.y = glucoseReads[glucoseReads.length - 1].y;
    } else {
        coordinates.y = interpolateGlucoseRead(glucoseReads[glucoseIndex - 1], glucoseReads[glucoseIndex], time);
        coordinates.x = time;
    }

    return coordinates;
}

function DateRange() {
    this.start;
    this.end;
}


function getNightsDateRangesForTimeWindow(lowerBound, upperBound) {
    nights = [];
    // Start from the end of the range.
    nightRange = getNightEndingAt(upperBound);

    nights.push(nightRange);
    while (nightRange.end >= lowerBound) {
        currentNightEnd = moment(nightRange.end).subtract('days', 1).toDate();
        nightRange = getNightEndingAt(currentNightEnd);
        nights.push(nightRange);
    }

    return nights;
}

function getNightEndingAt(endOfNight) {
    nightRange = new DateRange();
    nightRange.end = endOfNight;
    nightRange.end.setMinutes(30);
    nightRange.end.setHours(6);
    nightRange.end.setMilliseconds(0);
    nightRange.end.setSeconds(0);
    window.nightRange = nightRange;

    nightRange.start = moment(nightRange.end).subtract('hours', 8).toDate();
    return nightRange;
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

function isHoveringEventGroup(userEventGroup, searchElement) {
    var timestamp = searchElement.getTime() / 1000;

    // We consider hovering if the cursor is withing 780 seconds (13 minutes)
    // from the event
    distanceInSeconds = Math.abs(timestamp - userEventGroup.x);
    if (distanceInSeconds <= 780) {
        return true;
    } else {
        return false;
    }
}

function appendUserEventsToHoverBox(hoverbox, userEventGroup) {
    for (var i = 0; i < userEventGroup.userEvents.length; i++) {
        userEvent = userEventGroup.userEvents[i];
        eventTime = new Date(userEvent.x * 1000);

        lineText = userEvent.value;
        if (userEvent.tag === "Insulin") {
            lineText = lineText + " units";
        } else {
            lineText = lineText + " grams";
        }
        lineText = lineText + " at " + moment(eventTime).format("HH:mm");

        hoverbox.append("p")
            .attr("class", userEvent.tag)
            .text(lineText);
    }
}

function getDayBoundaries(lowerTimestamp, upperTimestamp) {
    // Round the upper snap limit to the nearest day
    var upperTime = new Date(upperTimestamp * 1000);
    var lowerTime = new Date(lowerTimestamp * 1000);
    upperTime.setMinutes(0);
    upperTime.setHours(0);
    upperTime.setSeconds(0);
    upperTime.setMilliseconds(0);
    dayBoundaries = [upperTime];

    currentSnapGuide = upperTime;
    while (currentSnapGuide > lowerTime) {
        currentSnapGuide = moment(currentSnapGuide).subtract('days', 1);
        dayBoundaries.push(currentSnapGuide.toDate());
    }

    return dayBoundaries;
}