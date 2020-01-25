// Reload if re-visiting using back/forward buttons.
if (String(window.performance.getEntriesByType("navigation")[0].type) === "back_forward"){
  location.reload();
}

let missedKas = 0

// Register with box event source
let source = new EventSource("/events/");
source.onmessage = function(event) {
    event = JSON.parse(event.data);
    switch(event.type) {
      case "keepalive":
        keepalive();

        break;

      case "updateBox":
        let targetBox = document.getElementById(event.id);
        if (targetBox !== null) { changeAlertLevel(targetBox, event.status, event.lastMessage); }

        break;

      case "deleteBox":
        deleteBox(event.id);

        break;

      case "reloadPage":
        location.reload()
    }
};

function boxHover(tip) {
    let target = document.getElementById("tooltip")
    target.innerHTML = tip
    target.display = "block"
}


function boxOut() {
    let target = document.getElementById("tooltip")
    target.innerHTML = ""
    target.display = "hidden"
}


function boxClick(id) {
    window.location.href = "/box/" + id;
}


function deleteBox(id) {
  let target = document.getElementById(id);
  target.parentNode.removeChild(target);
}


function keepalive() {
  let target = document.getElementById("status-bar");
  changeAlertLevel(target, "green", "");
  if(typeof ka !== "undefined") { clearTimeout(ka); };
  ka = setTimeout(
    function(){
      missedKas++;
      if (missedKas > 5 ) { location.reload(); };
      if ((now() - lastKa) > 60 ) { location.reload(); };
      lastKa = now();
      changeAlertLevel(target, "red", "ERROR: No keepalives for 5s.");
    }
    , 5 * 1000
  );
}


function myTime(t) {
    let r;
    if (t != null) {
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
    if ( ["amber","green","grey","noUpdate","red"].indexOf(status) === -1) { status = "grey" }
    target.classList.remove("amber", "green", "grey", "noUpdate", "red");
    target.classList.add(status);
    target.getElementsByClassName("message")[0].innerHTML = message;
    target.getElementsByClassName("lastUpdated")[0].innerHTML = myTime()
}


function rightSizeBigBox() {
    let availableWidth = Math.floor((window.innerWidth -30) / 512) * 512;
    let widthBox = (availableWidth >= 512) ? availableWidth:512;
    document.getElementById('big-box').style.width = widthBox + "px";
    let fullWidthBoxes = document.getElementsByClassName('fullwidth');
    for (let i = 0; i < fullWidthBoxes.length; i++) {
        fullWidthBoxes[i].style.width = (widthBox - 2) + "px";
    }
}
