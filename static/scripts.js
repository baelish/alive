function boxClick(id) {
    var t = Date.now()
    var pt = myTime(t)
    var target = document.getElementById(id)
    var width = window.innerWidth
    target.getElementsByClassName("message")[0].innerHTML = "Changed " + width
    changeAlertLevel(target, "amber")
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

function changeAlertLevel(target, level) {
    if ( level.indexOf(["amber", "green", "grey", "red"]) !== -1) { level = "grey" }
    target.classList.remove("amber", "green", "grey", "red")
    target.classList.add(level)
}

function rightSizeBigBox() {
    var availableWidth = Math.floor(window.innerWidth / 512) * 512
    widthBox = (availableWidth >= 1024) ? availableWidth:1024
    document.getElementById('big-box').style.width = widthBox + "px"
}
