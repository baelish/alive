var source = new EventSource("/events/");
source.onmessage = function(event) {
    var event = JSON.parse(event.data)
    switch(event.type) {
      case "keepalive":
        keepalive()

        break;

      case "updateBox":
        var targetBox = document.getElementById(event.id);
        changeAlertLevel(targetBox, event.status, event.lastMessage);

        break;

      case "deleteBox":
        deleteBox(event.id)

        break;

      case "reloadPage":
        location.reload()
    }
};


function boxClick(id) {
}

function deleteBox(id) {
  var target = document.getElementById(id)
  target.parentNode.removeChild(target)
}

function keepalive() {
  var target = document.getElementById("status-bar")
  changeAlertLevel(target, "green", "")
  if(typeof ka !== "undefined") { clearTimeout(ka) }
  ka = setTimeout(function(){changeAlertLevel(target, "red", "ERROR: No keepalives for 5s.")}, 5 * 1000)
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


function changeAlertLevel(target, status, message) {
    if ( ["amber","green","grey","noUpdate","red"].indexOf(status) == -1) { status = "grey" }
    target.classList.remove("amber", "green", "grey", "noUpdate", "red")
    target.classList.add(status)
    target.getElementsByClassName("message")[0].innerHTML = message
    target.getElementsByClassName("lastUpdated")[0].innerHTML = myTime()
}


function rightSizeBigBox() {
    var availableWidth = Math.floor((window.innerWidth -30) / 512) * 512
    widthBox = (availableWidth >= 512) ? availableWidth:512
    document.getElementById('big-box').style.width = widthBox + "px"
    document.getElementById('status-bar').style.width =( widthBox -2 ) + "px"
}
