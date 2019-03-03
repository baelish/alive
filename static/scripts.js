var timeouts = []

function boxClick(id) {
    var t = Date.now()
    var pt = myTime(t)
    var target = document.getElementById(id)
    alertNoUpdate(id, 60)
    changeAlertLevel(target, "green", "Testing only")
}


function myTime(t) {
    if ( t != null ) {
        r = new Date(t)
    } else {
        r = new Date()
    }
    return r.getFullYear() + pad(r.getMonth() + 1,2) + pad(r.getDay(),2) + "T" + pad(r.getHours(),2) +
        r.getMinutes() + pad(r.getSeconds(),2)
}


function pad(n, width, z) {
  z = z || '0';
  n = n + '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}


function changeAlertLevel(target, level, message) {
    if ( level.indexOf(["amber", "green", "grey", "red"]) !== -1) { level = "grey" }
    target.classList.remove("amber", "green", "grey", "red")
    target.classList.add(level)
    target.getElementsByClassName("message")[0].innerHTML = message
}


function rightSizeBigBox() {
    var availableWidth = Math.floor(window.innerWidth / 512) * 512
    widthBox = (availableWidth >= 1024) ? availableWidth:1024
    document.getElementById('big-box').style.width = widthBox + "px"
}


function alertNoUpdate(id, time) {
    if(typeof timeouts[id] !== "undefined") { clearTimeout(timeouts[id])};
    var target = document.getElementById(id)
    timeouts[id] = setTimeout(function(){changeAlertLevel(target, "red", "ERROR: No updates for " + time + "s.")}, time * 1000)
}
