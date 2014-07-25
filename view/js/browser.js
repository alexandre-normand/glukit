var TARGET_RANGE_LOWER_BOUND = 63;
var TARGET_RANGE_UPPER_BOUND = 127;
// Default to 5 hours
var SNAP_DISTANCE_IN_MILLIS = 5 * 2700000;
var RANGES = {
    HIGH: "HIGH",
    NORMAL: "NORMAL",
    LOW: "LOW"
}

function showDataBrowser(pathPrefix) {
    var chartNode = document.getElementById("chart");
    chartNode.innerHTML = '';
    var margin = {
            top: 10,
            right: 10,
            bottom: 100,
            left: 40
        },
        margin2 = {
            top: 430,
            right: 10,
            bottom: 20,
            left: 40
        },
        width = 800 - margin.left - margin.right,
        height = 500 - margin.top - margin.bottom,
        viewfinderHeight = 500 - margin2.top - margin2.bottom;
    // Our dates are milliseconds since epoch so that parse function
    // just instantiates a new Date from it.
    var parseDate = function(d) {
        return new Date(d)
    };
    // Define the scales
    var x = d3.time.scale().range([0, width]),
        x2 = d3.time.scale().range([0, width]),
        y = d3.scale.linear().range([height, 0]),
        y2 = d3.scale.linear().range([viewfinderHeight, 0]);
    // Define the axes
    var xAxis = d3.svg.axis().scale(x).orient("bottom"),
        xAxis2 = d3.svg.axis().scale(x2).orient("bottom"),
        yAxis = d3.svg.axis().scale(y).orient("left");
    var glucoseLine = d3.svg.line()
        .x(function(d) {
            return x(d.date);
        })
        .y(function(d) {
            return y(d.y);
        });
    var viewfinderLine = d3.svg.line()
        .x(function(d) {
            return x2(d.date);
        })
        .y(function(d) {
            return y2(d.y);
        });
    // Find a way to change those as the viewfinder moves
    var userEventArc = d3.svg.arc()
        .innerRadius(0)
        .outerRadius(7)
        .startAngle(function(d) {
            switch (d.type) {
                case "full":
                    return 0 * Math.PI;
                case "left":
                    return 1 * Math.PI;
                case "right":
                    return 1 * Math.PI;
            }
        })
        .endAngle(function(d) {
            switch (d.type) {
                case "full":
                    return 2 * Math.PI;
                case "left":
                    return 2 * Math.PI;
                case "right":
                    return 0 * Math.PI;
            }
        });
    var svg = d3.select("#chart").append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .attr("align", "center");
    // That should be the lower bound of the times displayed on the graph
    // driven by the "self" data
    d3.json("/" + pathPrefix + "data", function(error, data) {
        showProfile(data, "self_dashboard");
        addTrendToProfile(pathPrefix, data, "self_dashboard");

        glucoseReads = data.data[0].data;
        userEvents = data.data[1].data;
        timeRangeLowerBound = glucoseReads[0].x;
        glucoseReads.forEach(function(d) {
            d.date = parseDate(d.x * 1000);
        });
        var highestOnChart = d3.max(glucoseReads, function(d) {
            return d.y;
        }) + 100;
        x.domain(d3.extent(glucoseReads.map(function(d) {
            return d.date;
        })));
        y.domain([0, highestOnChart]);
        x2.domain(x.domain());
        y2.domain(y.domain());
        var focus = svg.append("g")
            .attr("id", "mainContainer")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")");
        var context = svg.append("g")
            .attr("transform", "translate(" + margin2.left + "," + margin2.top + ")");
        focus.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate(0," + height + ")")
            .call(xAxis);
        focus.append("rect")
            .attr("class", "target_range")
            .attr("clip-path", "url(#rangeClip)")
            .attr("width", width)
            .attr("height", height);
        focus.append("g")
            .attr("class", "y axis")
            .call(yAxis);
        var focusLine = d3.svg.line()
            .x(function(d) {
                return x(d.x);
            })
            .y(function(d) {
                return d.y;
            })
            .interpolate('linear');
        var focusCoordinates = [];
        focusCoordinates[0] = new Object();
        focusCoordinates[1] = new Object();
        focusCoordinates[0].x = 0;
        focusCoordinates[0].y = 0;
        focusCoordinates[1].x = 0;
        focusCoordinates[1].y = height;
        var focusLineElement = focus.append("path")
            .attr("id", "focusLine")
            .attr("class", "focusLine")
            .attr("clip-path", "url(#clip)")
            .datum(focusCoordinates)
            .attr("d", focusLine);
        var hoverbox = d3.select("#hoverbox");
        var chartContainerElement = d3.select("#chart_container")[0][0];
        context.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate(0," + viewfinderHeight + ")")
            .call(xAxis2);
        addTargetRangeClip(svg, width, height, "rangeClip", TARGET_RANGE_LOWER_BOUND, TARGET_RANGE_UPPER_BOUND, highestOnChart, y);
        var dayBoundaries = getDayBoundaries(glucoseReads[0].x, glucoseReads[glucoseReads.length - 1].x);
        var dayBoundaryGroup = focus.append("g").attr("id", "dayBoundaries");
        dayBoundaryGroup.selectAll(".dayBoundary")
            .data(dayBoundaries)
            .enter()
            .append("text")
            .attr("class", "dayBoundary")
            .attr("clip-path", "url(#clip)")
            .text(function(d) {
                return moment(d).format("dddd, Do");
            })
            .attr("y", "20")
            .attr("x", function(d) {
                return x(d);
            });
        var nights = getNightsDateRangesForTimeWindow(moment.unix(glucoseReads[0].x).toDate(), moment.unix(glucoseReads[glucoseReads.length - 1].x).toDate());
        var timebar = focus.append("g")
            .attr("id", "timebar");
        timebar
            .selectAll(".night")
            .data(nights)
            .enter()
            .append("rect")
            .attr("class", "night")
            .attr("clip-path", "url(#clip)")
            .attr("width", function(d) {
                return x(d.end) - x(d.start);
            })
            .attr("height", 5)
            .attr("x", function(d) {
                return x(d.start);
            })
            .attr("y", height - 5);
        var segments = splitReadsInRangeSegments(glucoseReads);
        addToGraph(focus, "self", context, segments, glucoseReads, y, glucoseLine, true, viewfinderLine);
        // Trying out the grouping, it doesn't actually use any of this
        userEvents.forEach(function(d) {
            d.date = parseDate(d.x * 1000);
        });
        window.userEvents = userEvents;
        var userEventGroups = groupEvents(userEvents, 60);
        window.userEventGroups = userEventGroups;
        // The clip path needs to be set on the group of userEvents and NOT on the individual arcs because otherwise,
        // the arcs will get clipped before the transformation and we're never going to be able to get a full circle or
        // semi-circle
        userEventsSvg = focus.append("g").attr("class", "userEvents").attr("clip-path", "url(#clip)");
        for (var i = 0; i < userEventGroups.length; i++) {
            markers = generateEventMarkers(userEventGroups[i]);
            for (var j = 0; j < markers.length; j++) {
                userEventsSvg.append("path")
                    .attr("class", "event " + markers[j].tag)
                    .attr("transform", "translate(" + x(markers[j].date) + "," + y(markers[j].y) + ")")
                    .datum(markers[j])
                    .attr("d", userEventArc);
            }
        }

        function brushed() {
            x.domain(brush.empty() ? x2.domain() : brush.extent());
            focus.selectAll("path.self").attr("d", glucoseLine);
            focus.selectAll("path.steadySailor").attr("d", glucoseLine);
            focus.selectAll("#dayBoundaries")
            dayBoundaryGroup.selectAll(".dayBoundary")
                .attr("x", function(d) {
                    return x(d);
                });
            timebar
                .selectAll(".night")
                .attr("width", function(d) {
                    return x(d.end) - x(d.start);
                })
                .attr("x", function(d) {
                    return x(d.start);
                });
            focus.selectAll("path.event").attr("transform", function(d) {
                return "translate(" + x(d.date) + "," + y(d.y) + ")";
            });
            focus.select(".x.axis").call(xAxis);
            extent = brush.extent();
            viewfinderUpperLimit = extent[1];
            viewfinderLowerLimit = extent[0];
            //console.log("limits are " + viewfinderLowerLimit + ", " + viewfinderUpperLimit);
            rangeAggregate = getRangeAggregate(segments, viewfinderLowerLimit, viewfinderUpperLimit);
            var lowRangePercentage = Math.floor(rangeAggregate.lowTimeInMinutes / rangeAggregate.getTotalTime() * 100);
            var highRangePercentage = Math.floor(rangeAggregate.highTimeInMinutes / rangeAggregate.getTotalTime() * 100);
            var normalRangePercentage = Math.floor(rangeAggregate.normalTimeInMinutes / rangeAggregate.getTotalTime() * 100);
            $("#lowRangePercentage").text(lowRangePercentage + "% of the time");
            $("#normalRangePercentage").text(normalRangePercentage + "% of the time");
            $("#highRangePercentage").text(highRangePercentage + "% of the time");
        }
        // snap to days if close enough

        function snapBrush() {
            if (!d3.event.sourceEvent) return; // only transition after input
            extent = brush.extent();
            viewfinderUpperLimit = extent[1];
            // We iterate from most recent since we assume that most users will look at that data more often
            for (var i = 0; i < snapGuides.length; i++) {
                snap = snapGuides[i]
                diffInMillis = Math.abs(viewfinderUpperLimit - snap);
                if (diffInMillis <= SNAP_DISTANCE_IN_MILLIS) {
                    snapExtent = getDayRangeFromUpperBound(snap.getTime() / 1000);
                    d3.select(this).transition()
                        .call(brush.extent(snapExtent))
                        .call(brush.event)
                        .duration(500);
                    return;
                }
            }
        }
        mostRecentReadTimeInSeconds = glucoseReads[glucoseReads.length - 1].x;
        brushRange = getDayRangeFromUpperBound(mostRecentReadTimeInSeconds);
        // 1 snap every day for 7 days
        var snapGuides = getDateSnapGuides(moment.unix(glucoseReads[0].x).toDate(), moment.unix(mostRecentReadTimeInSeconds).toDate());
        var brush = d3.svg.brush()
            .x(x2)
            .extent(brushRange)
            .on("brush", brushed)
            .on("brushend", snapBrush);
        brushElement = context.append("g");
        brushElement.attr("class", "x brush")
            .call(brush)
            .selectAll("rect")
            .attr("y", -6)
            .attr("height", viewfinderHeight + 7);
        // Remove resize handles because we don't want to allow resizing. The brush is fixed to the length of one day
        brushElement.selectAll(".resize").remove();
        // Disable cross-hair cursor to avoid user overriding the brush completely
        brushElement.selectAll(".extent").style("pointer-events", "all");
        brushElement.selectAll(".background").style("pointer-events", "none");
        // Call a first brushed to view the last day only
        brushed();
        toggleNormalRange();
        // Add steady sailor data
        d3.json("/" + pathPrefix + "steadySailor", function(error, data) {
            if (data != undefined) {
                showProfile(data, "sailor_dashboard");
                rawSteadySailor = data.data[0].data;
                steadySailor = rawSteadySailor
                if (steadySailor.length > 0) {
                    offset = steadySailor[0].x - timeRangeLowerBound;
                    steadySailor.forEach(function(element) {
                        element.x = element.x - offset;
                        element.date = parseDate(element.x * 1000);
                    });
                }
                sailorSegments = splitReadsInRangeSegments(steadySailor);
                addToGraph(focus, "steadySailor", context, sailorSegments, steadySailor, y, glucoseLine, false);
            }
        });
        addBackgroundAndHover(focus, glucoseReads, userEventGroups, width, height, x, y, focusCoordinates, chartContainerElement, hoverbox, focusLine);
    });
}

function addToGraph(focus, className, context, segments, glucoseReads, y, glucoseLineFunc, highlightSegments, viewfinderLineFunc) {
    for (var i = 0; i < segments.length; i++) {
        var segment = segments[i];
        var segmentClass = className;

        // if highlightSegments is enabled, append the segment name to the class names
        if (highlightSegments === true) {
            segmentClass = className + " " + segment.range;
        }

        focus.append("path")
            .attr("class", segmentClass)
            .attr("clip-path", "url(#clip)")
            .datum(segment.reads)
            .attr("d", glucoseLineFunc);
    }
    if (viewfinderLineFunc != undefined) {
        context.append("path")
            .datum(glucoseReads)
            .attr("d", viewfinderLineFunc);
    }
}

function addTargetRangeClip(svg, canvasWidth, canvasHeight, id, lowerBound, upperBound, maxDomainValue, y) {
    svg.append("defs").append("clipPath")
        .attr("id", id)
        .append("rect")
        .attr("width", canvasWidth)
        .attr("height", y(maxDomainValue - upperBound) - y(maxDomainValue - lowerBound))
        .attr("y", y(upperBound));
    svg.append("defs").append("clipPath")
        .attr("id", "clip")
        .append("rect")
        .attr("width", canvasWidth)
        .attr("height", canvasHeight);
}

function addBackgroundAndHover(focus, glucoseReads, userEventGroups, width, height, x, y, focusCoordinates, chartContainerElement, hoverbox, focusLine) {
    // We add a background so that we have a svg area to register the mouse
    // movement
    var background = focus.append("rect")
        .attr("id", "background")
        .attr("class", "background")
        .attr("clip-path", "url(#clip)")
        .attr("width", width)
        .attr("height", height);
    chartBackground = background[0][0];
    chartBackground.addEventListener('mousemove', function(event) {
        var rect = this.getBoundingClientRect();
        var left = event.clientX - rect.left - this.clientLeft + this.scrollLeft;
        var top = event.clientY - rect.top - this.clientTop + this.scrollTop;
        var time = x.invert(left);
        coordinates = getHoverCoordinates(glucoseReads, time);
        focusCoordinates[1].x = focusCoordinates[0].x = coordinates.x;
        focus.selectAll("#focusLine").attr("d", focusLine);
        hoverleftposition = event.clientX - chartContainerElement.getBoundingClientRect().left;
        hoverbox.style("left", hoverleftposition + "px");
        // Clear previous content
        hoverbox.text(null);
        hoverbox.append("p")
            .attr("class", "glucose")
            .text(Math.round(coordinates.y) + " mg/dl");
        var userEventGroupIndex = Math.abs(binaryIndexOf.call(userEventGroups, time));
        // We beyond the last event group, check if the last is close enough
        if (userEventGroupIndex >= userEventGroups.length) {
            if (isHoveringEventGroup(userEventGroups[userEventGroups.length - 1], time)) {
                userEventGroup = userEventGroups[userEventGroups.length - 1];
                appendUserEventsToHoverBox(hoverbox, userEventGroup);
            }
        }
        // Look for the event after the current position to see if it's close enough
        else if (isHoveringEventGroup(userEventGroups[userEventGroupIndex], time)) {
            userEventGroup = userEventGroups[userEventGroupIndex];
            appendUserEventsToHoverBox(hoverbox, userEventGroup);
        }
        // Look to see if the one prior is close enough
        else if (userEventGroupIndex > 0 && isHoveringEventGroup(userEventGroups[userEventGroupIndex - 1], time)) {
            userEventGroup = userEventGroups[userEventGroupIndex - 1];
            appendUserEventsToHoverBox(hoverbox, userEventGroup);
        }
    }, false);
    chartBackground.addEventListener('mouseover', function(event) {
        d3.select(".hoverbox").style("display", "block");
        d3.select(".focusLine").style("display", "block");
    });
    chartBackground.addEventListener('mouseout', function(event) {
        d3.select(".hoverbox").style("display", "none");
        d3.select(".focusLine").style("display", "none");
    });
}

function showProfile(dashboardData, sectionName) {
    var sectionPrefix = sectionName + ".";
    var glukitScore = dashboardData.score;
    var profilePicture = dashboardData.picture;
    var firstName = dashboardData.firstName;
    var lastName = dashboardData.lastName;
    var lastSync = dashboardData.lastSync;
    var lowerBound = dashboardData.scoreDetails.LowerBound;
    var upperBound = dashboardData.scoreDetails.UpperBound;
    var joinedOn = dashboardData.joinedOn;
    var img = document.createElement("img");
    img.src = profilePicture;
    img.className = "avatar bubble";
    var src = document.getElementById(sectionPrefix + "profilePic");
    src.appendChild(img);
    document.getElementById(sectionPrefix + "userName").innerHTML = firstName + " " + lastName;
    document.getElementById(sectionPrefix + "lastSync").innerHTML = moment(lastSync, "YYYY-MM-DDHH:mm:ssZ").fromNow();
    if (glukitScore == null) {
        document.getElementById(sectionPrefix + "glukitScore").innerHTML = "?";
    } else {
        document.getElementById(sectionPrefix + "glukitScore").innerHTML = glukitScore;
    }
    document.getElementById(sectionPrefix + "glukitScoreDates").innerHTML = moment(lowerBound, "YYYY-MM-DDHH:mm:ssZ").format('LL') + ' - ' + moment(upperBound, "YYYY-MM-DDHH:mm:ssZ").format('LL');
}

function addTrendToProfile(pathPrefix, dashboardData, sectionName) {
    var sectionPrefix = sectionName + ".";
    var glukitScore = dashboardData.score;
    $.getJSON("/" + pathPrefix + "glukitScores?limit=8", function(data) {
        var mostRecentScore = data[0].Value;
        var referenceScore = data[data.length - 1].Value;
        var trendClass = mostRecentScore < referenceScore ? "icon-up-dir score-trend up" : "icon-down-dir score-trend down";

        document.getElementById(sectionPrefix + "glukitScore").innerHTML = glukitScore + "<i class=\"" + trendClass + "\"></i>";
    });
}

function toggleNormalRange() {
    $("#normalRangePercentage").show();
    $("#lowRangePercentage").hide();
    $("#highRangePercentage").hide();
    //change the color of the text to white and the background to pink
    $("#inTargetRangeText").css("color", "#fff");
    $("#inRange").css("background", "#9e2064");
    //reset the other div colors
    $("#belowTargetRangeText").css("color", "#404041");
    $("#below").css("background", "#eee");
    $("#aboveTargetRangeText").css("color", "#404041");
    $("#above").css("background", "#eee");
    $("#below").show();
    $("#above").show();
    $('.NORMAL').css("stroke", "#2390de");
    $('.LOW').css("stroke", "#ddd");
    $('.HIGH').css("stroke", "#ddd");
}

function highlightLines() {
    $("#inRange").click(toggleNormalRange);
    $("#below").click(function() {
        $("#above").show();
        $("#inRange").show();
        //change colors of selected div
        $("#belowTargetRangeText").css("color", "#fff");
        $("#below").css("background", "#9e2064");
        //reset the other div colors
        $("#inTargetRangeText").css("color", "#404041");
        $("#inRange").css("background", "#eee");
        $("#aboveTargetRangeText").css("color", "#404041");
        $("#above").css("background", "#eee");
        //show and hide percentages
        $("#normalRangePercentage").hide();
        $("#lowRangePercentage").show();
        $("#highRangePercentage").hide();
        //Highlight relevant pieces of the graph
        $('.NORMAL').css("stroke", "#ddd");
        $('.LOW').css("stroke", "#de2333");
        $('.HIGH').css("stroke", "#ddd");
    });
    $("#above").click(function() {
        //change colors of selected div
        $("#above").css("background", "#9e2064");
        $("#aboveTargetRangeText").css("color", "#fff");
        //reset colors of other divs
        $("#below").css("background", "#eee");
        $("#belowTargetRangeText").css("color", "#404041");
        $("#inRange").css("background", "#eee");
        $("#inTargetRangeText").css("color", "#404041");
        $("#below").show();
        $("#inRange").show();
        $("#normalRangePercentage").hide();
        $("#lowRangePercentage").hide();
        $("#highRangePercentage").show();
        $('.NORMAL').css("stroke", "#ddd");
        $('.LOW').css("stroke", "#ddd");
        $('.HIGH').css("stroke", "#de2333");
    });
}