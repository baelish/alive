/* jshint esversion: 6 */

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
      if ( window.location.pathname === "/" || window.location.pathname === `/box/${event.id}`) {
        updateBox(event);
      }

      break;

    case "deleteBox":
      if ( window.location.pathname === "/" || window.location.pathname === `/box/${event.id}`) {
        deleteBox(event.id);
      }

      break;

    case "createBox":
      if ( window.location.pathname === "/" ) {
        createBox(event.after, event.box);
      }

      break;

    case "reloadPage":
      location.reload();
  }
};


// Box tooltip
function boxHover(tip) {
  let target = document.getElementById("tooltip");
  target.innerHTML = tip;
  target.display = "block";
}


function boxOut() {
  let target = document.getElementById("tooltip");
  target.innerHTML = "";
  target.display = "hidden";
}


// Box click
function boxClick(id) {
  window.location.href = "/box/" + id;
}


// Create box
function createBox(after, box) {
  let title;
  if (box.displayname) {
    title = box.displayname;
  } else {
    title = box.name;
  }

  let divContent = `
    <div onclick='boxClick(this.id)' onmouseover='boxHover("${box.name}")' onmouseout='boxOut()' id='${box.id}' class='${box.status} ${box.size} box'>
        <p class='title'>${title}</p>
        <p class='message'>${box.lastMessage}</p>
        <p class='lastUpdated'>${box.lastUpdate}</p>
        <p class='maxTBU'>${box.maxTBU}</p>
        <p class='expireAfter'>${box.expireAfter}</p>
    </div>
  `;

  let precedingBox = document.getElementById(after);
  precedingBox.insertAdjacentHTML("afterEnd", divContent);
}


// Remove box
function deleteBox(id) {
  let target = document.getElementById(id);
  target.parentNode.removeChild(target);
}


// Update box
function updateBox(event) {
  let targetBox = document.getElementById(event.id);

  if (targetBox !== null) {
    changeAlertLevel(targetBox, event.status, event.lastMessage);
  }

  if (event.maxTBU) {
    let row = targetBox.getElementsByClassName("maxTBU")[0];
    row.getElementsByTagName('td')[0].innerHTML = event.maxTBU;
    if (event.maxTBU === "0") {row.style.display = "none";} else {row.style.display = "table-row";}
  }

  if (event.expireAfter) {
    let row = targetBox.getElementsByClassName("expireAfter")[0];
    row.getElementsByTagName('td')[0].innerHTML = event.expireAfter;
    if (event.expireAfter === "0") {row.style.display = "none";} else {row.style.display = "table-row";}
  }
}


// keepalive
let lastKa;
function keepalive() {
  let ct = new Date().getTime();
  if ( lastKa && (lastKa + 60000) < ct ) { location.reload();}
  lastKa = ct;
  let target = document.getElementById("status-bar");
  target.classList.remove("amber","green","grey","noUpdate","red");
  target.classList.add("green");
  target.getElementsByClassName("message")[0].innerHTML = "";
  if(typeof ka !== "undefined") { clearTimeout(ka); }
  ka = setTimeout(
    function(){
      target.classList.remove("amber","green","grey","noUpdate","red");
      target.classList.add("noUpdate");
      target.getElementsByClassName("message")[0].innerHTML = "ERROR: No keepalives since " + myTime(lastKa) + ".";
    }, 5 * 1000
  );
}


// Print time in my preferred format
function myTime(t) {
  let r;
  if (t != null) {
    r = new Date(t);
  } else {
    r = new Date();
  }
  return r.toISOString();
}


// Pad a string
function pad(n, width, z) {
  z = z || '0';
  n = n + '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}


// Change alert level of a box
function changeAlertLevel(target, status, message) {
  if ( ["amber","green","grey","noUpdate","red"].indexOf(status) === -1) { status = "grey";}
  target.classList.remove("amber", "green", "grey", "noUpdate", "red");
  target.classList.add(status);
  target.getElementsByClassName("message")[0].innerHTML = message;
  target.getElementsByClassName("lastUpdated")[0].innerHTML = myTime();
  let pMessages = target.getElementsByClassName("previousMessages");
  if (typeof(pMessages[0] != 'undefined') && pMessages[0] != null) {
    pMessages[0].insertAdjacentHTML('afterbegin', "<li>" + myTime() + ": " + status.toUpperCase() + " (" + message + ")</li>");
  }
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
