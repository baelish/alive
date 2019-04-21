var timeouts = []

var source = new EventSource("/events/");
source.onmessage = function(event) {
    var eventDetails = event.data.split(',')
    switch(eventDetails[0]) {
      case "updateBox":
        var targetBox = document.getElementById(eventDetails[1]);
        changeAlertLevel(targetBox, eventDetails[2], eventDetails[3]);
        break;
      case "reloadPage":
        location.reload()
    }
};


function boxClick(id) {
}


function myTime(t) {
    if ( t != null ) {
        r = new Date(t)
    } else {
        r = new Date()
    }
    return r.getFullYear() + "-" + pad(r.getMonth() + 1,2) + "-" +
        pad(r.getDate(),2) + "T" + pad(r.getHours(),2) + ":" +
        pad(r.getMinutes(),2) + ":" + pad(r.getSeconds(),2)
}


function pad(n, width, z) {
  z = z || '0';
  n = n + '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}


function changeAlertLevel(target, level, message) {
    if ( ["amber","green","grey","red"].indexOf(level) == -1) { level = "grey" }
    target.classList.remove("amber", "green", "grey", "red")
    target.classList.add(level)
    target.getElementsByClassName("message")[0].innerHTML = message
    target.getElementsByClassName("lastUpdated")[0].innerHTML = myTime()
}


function rightSizeBigBox() {
    var availableWidth = Math.floor((window.innerWidth -30) / 512) * 512
    widthBox = (availableWidth >= 512) ? availableWidth:512
    document.getElementById('big-box').style.width = widthBox + "px"
    document.getElementById('status-bar').style.width =( widthBox -2 ) + "px"
}


function alertNoUpdate(id, time) {
    if(typeof timeouts[id] !== "undefined") { clearTimeout(timeouts[id])};
    var target = document.getElementById(id)
    timeouts[id] = setTimeout(function(){changeAlertLevel(target, "red", "ERROR: No updates for " + time + "s.")}, time * 1000)
}
