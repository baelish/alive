// Reload if re-visiting using back/forward buttons.
if (String(window.performance.getEntriesByType("navigation")[0].type) === "back_forward"){
  location.reload();
}


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

        if (targetBox !== null) {
          changeAlertLevel(targetBox, event.status, event.lastMessage);
        }

        if (event.maxTBU !== "") {
          let row = targetBox.getElementsByClassName("maxTBU")[0]
          row.getElementsByTagName('td')[0].innerHTML = event.maxTBU;
          if (event.maxTBU === "0") {row.style.display = "none"} else {row.style.display = "table-row"}
        }

        if (event.expireAfter !== "") {
          let row = targetBox.getElementsByClassName("expireAfter")[0]
          row.getElementsByTagName('td')[0].innerHTML = event.expireAfter;
          if (event.expireAfter === "0") {row.style.display = "none"} else {row.style.display = "table-row"}
        }

        break;

      case "deleteBox":
        deleteBox(event.id);

        break;

      case "reloadPage":
        location.reload()
    }
};


// Box tooltip
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


// Box click
function boxClick(id) {
    window.location.href = "/box/" + id;
}


// Remove box
function deleteBox(id) {
  let target = document.getElementById(id);
  target.parentNode.removeChild(target);
}


// keepalive
let lastKa;
function keepalive() {
  let ct = new Date().getTime()
  if ( lastKa && (lastKa + 60000) < ct ) { location.reload() }
  lastKa = ct
  let target = document.getElementById("status-bar");
  changeAlertLevel(target, "green", "");
  if(typeof ka !== "undefined") { clearTimeout(ka); };
  ka = setTimeout(
    function(){
        changeAlertLevel(target, "red", "ERROR: No keepalives since " + myTime(lastKa) + ".")
    }, 5 * 1000)
}


// Print time in my preferred format
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


// Pad a string
function pad(n, width, z) {
  z = z || '0';
  n = n + '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}


// Change alert level of a box
function changeAlertLevel(target, status, message) {
    if ( ["amber","green","grey","noUpdate","red"].indexOf(status) === -1) { status = "grey" }
    target.classList.remove("amber", "green", "grey", "noUpdate", "red");
    target.classList.add(status);
    target.getElementsByClassName("message")[0].innerHTML = message;
    target.getElementsByClassName("lastUpdated")[0].innerHTML = myTime()
}


// Make big box to fit as many biggest boxes as will fit the current window.
function rightSizeBigBox() {
    let availableWidth = Math.floor((window.innerWidth -30) / 512) * 512;
    let widthBox = (availableWidth >= 512) ? availableWidth:512;
    document.getElementById('big-box').style.width = widthBox + "px";
    let fullWidthBoxes = document.getElementsByClassName('fullwidth');
    for (let i = 0; i < fullWidthBoxes.length; i++) {
        fullWidthBoxes[i].style.width = (widthBox - 2) + "px";
    }
}
