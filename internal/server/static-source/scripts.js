/* jshint esversion: 6 */

// Reload if re-visiting using back/forward buttons.
if (
  String(window.performance.getEntriesByType("navigation")[0].type) ===
  "back_forward"
) {
  location.reload();
}

// Register with box event source
let source = new EventSource("/events/");
source.onmessage = function (event) {
  event = JSON.parse(event.data);
  switch (event.type) {
    case "keepalive":
      keepalive();

      break;

    case "updateBox":
      if (
        window.location.pathname === "/" ||
        window.location.pathname === `/box/${event.id}`
      ) {
        updateBox(event);
      }

      break;

    case "deleteBox":
      if (
        window.location.pathname === "/" ||
        window.location.pathname === `/box/${event.id}`
      ) {
        deleteBox(event.id);
      }

      break;

    case "createBox":
      if (window.location.pathname === "/") {
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

let activeBoxId = null;
let resizingDrawer = false;
let resizeStartX = 0;
let resizeStartWidth = 0;
let dashboardStale = false;
let pendingDrawerReopenId = null;

// Initialize the dashboard drawer.
function initDashboard() {
  document.addEventListener("keydown", function (event) {
    if (event.key === "Escape") {
      closeDrawer();
    }
  });

  let resizeHandle = document.getElementById("drawer-resize-handle");
  window.addEventListener("resize", resizeDrawerToViewport);
  resizeHandle.addEventListener("pointerdown", function (event) {
    let panel = document.getElementById("side-drawer-panel");
    resizingDrawer = true;
    resizeStartX = event.clientX;
    resizeStartWidth = panel.getBoundingClientRect().width;
    resizeHandle.setPointerCapture(event.pointerId);
    document.body.classList.add("drawer-resizing");
    event.preventDefault();
  });

  resizeHandle.addEventListener("pointermove", function (event) {
    if (!resizingDrawer) {
      return;
    }
    let bounds = drawerWidthBounds();
    let maxWidth = bounds.max;
    let minWidth = bounds.min;
    let width = Math.max(minWidth, Math.min(maxWidth, resizeStartWidth + resizeStartX - event.clientX));
    document.getElementById("side-drawer-panel").style.width = width + "px";
    rightSizeBigBox("dashboard");
  });

  resizeHandle.addEventListener("pointerup", function (event) {
    resizingDrawer = false;
    resizeHandle.releasePointerCapture(event.pointerId);
    document.body.classList.remove("drawer-resizing");
  });
}

function drawerWidthBounds() {
  return {
    min: window.innerWidth * 0.3,
    max: window.innerWidth * 0.4,
  };
}

function resizeDrawerToViewport() {
  let panel = document.getElementById("side-drawer-panel");
  if (panel === null || !panel.classList.contains("is-open")) {
    return;
  }
  let bounds = drawerWidthBounds();
  let width = panel.getBoundingClientRect().width;
  width = Math.max(bounds.min, Math.min(bounds.max, width));
  panel.style.width = width + "px";
  rightSizeBigBox("dashboard");
}

// Box click
function boxClick(id) {
  let detail = document.getElementById("box-detail-" + id);
  if (detail === null) {
    return;
  }

  if (activeBoxId !== null) {
    let previous = document.getElementById(activeBoxId);
    if (previous !== null) {
      previous.classList.remove("selected");
    }
  }

  activeBoxId = id;
  document.getElementById(id).classList.add("selected");
  document.getElementById("drawer-content").innerHTML = detail.innerHTML;
  document.getElementById("side-drawer-panel").classList.add("is-open");
  document.getElementById("side-drawer-panel").setAttribute("aria-hidden", "false");
  document.getElementById("drawer-backdrop").classList.add("is-visible");
  document.body.classList.add("drawer-open");
  rightSizeBigBox("dashboard");
}

function closeDrawer() {
  if (activeBoxId !== null) {
    let target = document.getElementById(activeBoxId);
    if (target !== null) {
      target.classList.remove("selected");
    }
  }
  activeBoxId = null;
  document.getElementById("drawer-content").innerHTML = "";
  document.getElementById("side-drawer-panel").classList.remove("is-open");
  document.getElementById("side-drawer-panel").setAttribute("aria-hidden", "true");
  document.getElementById("drawer-backdrop").classList.remove("is-visible");
  document.body.classList.remove("drawer-open");
  rightSizeBigBox("dashboard");
}

function clearDashboardBoxes() {
  let grid = document.getElementById("boxes-grid");
  if (grid === null) {
    return;
  }
  grid.innerHTML = "";
  activeBoxId = null;
  let drawerContent = document.getElementById("drawer-content");
  let drawerPanel = document.getElementById("side-drawer-panel");
  let drawerBackdrop = document.getElementById("drawer-backdrop");
  if (drawerContent !== null) drawerContent.innerHTML = "";
  if (drawerPanel !== null) drawerPanel.classList.remove("is-open");
  if (drawerBackdrop !== null) drawerBackdrop.classList.remove("is-visible");
  document.body.classList.remove("drawer-open");
  dashboardStale = true;
}

// Create box
function createBox(after, box) {
  let title;
  if (box.displayName) {
    title = box.displayName;
  } else {
    title = box.name;
  }

  let divContent = `
    <div onclick='boxClick(this.id)' onmouseover='boxHover("${box.name}")' onmouseout='boxOut()' id='${box.id}' class='${box.status} ${box.size} box'>
        <p class='title'>${title}</p>
      <p class='box-id'>${box.id}</p>
        <p class='maxTBU'>${box.maxTBU}</p>
        <p class='expireAfter'>${box.expireAfter}</p>
    </div>
  `;

  let grid = document.getElementById("boxes-grid");
  if (grid === null) {
    return;
  }
  let precedingBox = document.getElementById(after);
  if (precedingBox !== null) {
    precedingBox.insertAdjacentHTML("afterEnd", divContent);
  } else {
    grid.insertAdjacentHTML("beforeend", divContent);
  }
  let newBox = document.getElementById(box.id);
  newBox.dataset.sortName = box.name;
  let oldDetail = document.getElementById("box-detail-" + box.id);
  if (oldDetail !== null) {
    oldDetail.remove();
  }

  let detail = document.createElement("template");
  detail.id = "box-detail-" + box.id;
  let maxTBU = box.maxTBU ? '<dt class="drawer-max-tbu-label">Max TBU</dt><dd class="drawer-max-tbu"><span class="drawer-metadata-value">' + box.maxTBU + '</span>' + metadataIcon("activity", "Max TBU") + '</dd>' : '';
  let expireAfter = box.expireAfter ? '<dt class="drawer-expire-after-label">Expires after</dt><dd class="drawer-expire-after"><span class="drawer-metadata-value">' + box.expireAfter + '</span>' + metadataIcon("trash-2", "Expires after") + '</dd>' : '';
  let links = box.links && box.links.length ? '<div class="drawer-section drawer-links"><h3>Links</h3>' + box.links.map(function (link) {
    return '<a href="' + link.url + '" target="_blank" rel="noopener noreferrer"><svg class="external-link-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg><span>' + link.name + '</span></a><br />';
  }).join("") + '</div>' : '';
  detail.innerHTML =
    '<div class="drawer-status-strip ' +
    box.status +
    '"></div><div class="drawer-header"><h2 class="drawer-title">' +
    title +
    '</h2><div class="drawer-actions"><a class="drawer-open-page" href="/box/' + encodeURIComponent(box.id) + '" target="_blank" rel="noopener noreferrer" aria-label="Open full box page" title="Open full box page"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg></a><button class="drawer-close" type="button" aria-label="Close details" onclick="closeDrawer()"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></svg></button></div></div>' +
    '<div class="drawer-section drawer-metadata"><dl class="drawer-meta"><dt>ID</dt><dd><span class="drawer-metadata-value">' +
    box.id +
    '</span></dd><dt>Size</dt><dd><span class="drawer-metadata-value">' + box.size + '</span></dd>' +
    maxTBU + expireAfter +
    '</dl></div>' + links + '<div class="drawer-section drawer-latest-update"><h3>Latest update</h3><p class="drawer-message">' +
    box.lastMessage +
    '</p><time class="drawer-last-updated">' + box.lastUpdate +
    '</time></div><div class="drawer-section drawer-events"><h3>Event log</h3><div class="event-log"></div></div>';
  let eventLog = detail.content.querySelector(".event-log");
  if (eventLog !== null && box.messages) {
    box.messages.forEach(function (message) {
      eventLog.append(createEventCard(message.status, message.message, message.timeStamp));
    });
  }
  newBox.after(detail);
  placeBoxInSizeOrder(newBox, detail);
  if (pendingDrawerReopenId === box.id) {
    pendingDrawerReopenId = null;
    boxClick(box.id);
  }
}

function placeBoxInSizeOrder(boxElement, detailElement) {
  let grid = document.getElementById("boxes-grid");
  let sizeOrder = ["xlarge", "dlarge", "large", "dmedium", "medium", "dsmall", "small", "dmicro", "micro", "dot"];
  let boxRank = boxSizeRank(boxElement, sizeOrder);
  if (boxRank === -1) {
    return;
  }

  boxElement.remove();
  detailElement.remove();

  let low = 0;
  let high = Math.floor(grid.children.length / 2);
  while (low < high) {
    let middle = Math.floor((low + high) / 2);
    let candidate = grid.children[middle * 2];
    let candidateRank = boxSizeRank(candidate, sizeOrder);
    if (candidateRank < boxRank ||
      (candidateRank === boxRank && candidate.dataset.sortName.localeCompare(boxElement.dataset.sortName) < 0)) {
      low = middle + 1;
    } else {
      high = middle;
    }
  }

  let candidate = grid.children[low * 2];
  if (candidate !== undefined) {
    grid.insertBefore(boxElement, candidate);
    grid.insertBefore(detailElement, candidate);
  } else {
    grid.append(boxElement, detailElement);
  }
}

function boxSizeRank(boxElement, sizeOrder) {
  return sizeOrder.findIndex(function (size) {
    return boxElement.classList.contains(size);
  });
}

// Remove box
function deleteBox(id) {
  let target = document.getElementById(id);
  if (target !== null) {
    target.parentNode.removeChild(target);
  }
  let detail = document.getElementById("box-detail-" + id);
  if (detail !== null) {
    detail.remove();
  }
  if (activeBoxId === id) {
    pendingDrawerReopenId = id;
    window.setTimeout(function () {
      if (pendingDrawerReopenId === id) {
        pendingDrawerReopenId = null;
        closeDrawer();
      }
    }, 1000);
  }
}

// Update box
function updateBox(event) {
  let targetBox = document.getElementById(event.id);

  if (targetBox !== null) {
    changeAlertLevel(targetBox, event.status, event.lastMessage);
  }

  let detail = document.getElementById("box-detail-" + event.id);
  if (detail !== null) {
    updateDetailContent(detail.content, event);
  }

  if (activeBoxId === event.id) {
    updateDetailContent(document.getElementById("drawer-content"), event);
  }

  if (targetBox === null) {
    return;
  }

  if (event.maxTBU) {
    let t = targetBox.getElementsByClassName("maxTBU")[0];
    if (t.tagName === "TR") {
      t.getElementsByTagName("TD")[0].innerHTML = event.maxTBU;
      if (event.maxTBU === "0s") {
        t.style.display = "none";
      } else {
        t.style.display = "table-row";
      }
    } else {
      t.innerHTML = event.maxTBU;
    }
  }

  if (event.expireAfter) {
    let t = targetBox.getElementsByClassName("expireAfter")[0];
    if (t.tagName === "TR") {
      t.getElementsByTagName("TD")[0].innerHTML = event.expireAfter;
      if (event.expireAfter === "0s") {
        t.style.display = "none";
      } else {
        t.style.display = "table-row";
      }
    } else {
      t.innerHTML = event.expireAfter;
    }
  }
}

function updateDetailContent(content, event) {
  let statusValue = event.status === "noUpdate" ? "red" : event.status || "grey";
  let dot = content.querySelector(".drawer-status-dot");
  let strip = content.querySelector(".drawer-status-strip");
  let message = content.querySelector(".drawer-message");
  let updated = content.querySelector(".drawer-last-updated");

  if (dot !== null) {
    dot.classList.remove("amber", "green", "grey", "noUpdate", "red");
    dot.classList.add(statusValue);
  }
  if (strip !== null) {
    strip.classList.remove("amber", "green", "grey", "noUpdate", "red");
    strip.classList.add(statusValue);
  }
  if (message !== null) message.textContent = event.lastMessage;
  if (updated !== null) updated.textContent = myTime();
  updateOptionalMetadata(content, "maxTBU", "Max TBU", "drawer-max-tbu-label", "drawer-max-tbu", event.maxTBU);
  updateOptionalMetadata(content, "expireAfter", "Expires after", "drawer-expire-after-label", "drawer-expire-after", event.expireAfter);

  let card = createEventCard(statusValue, event.lastMessage, null);
  addEventIndicator(card, "activity", "Max TBU", event.maxTBU);
  addEventIndicator(card, "trash-2", "Expires after", event.expireAfter);
  let log = content.querySelector(".event-log");
  if (log !== null) {
    log.prepend(card);
    while (log.children.length > 30) {
      log.lastElementChild.remove();
    }
  }
}

function createEventCard(status, message, timestamp) {
  let statusValue = status === "noUpdate" ? "red" : status || "grey";
  let card = document.createElement("div");
  card.className = "event-log-card " + statusValue;
  card.innerHTML = "<div class=\"event-body\"><div class=\"event-top\"><strong class=\"event-status\"></strong><div class=\"event-leading\"></div></div><div class=\"event-message\"></div><small></small></div>";
  card.getElementsByClassName("event-status")[0].textContent = statusValue;
  card.getElementsByClassName("event-message")[0].textContent = message;
  card.getElementsByTagName("small")[0].textContent = myTime(timestamp);
  return card;
}

function addEventIndicator(card, icon, label, value) {
  if (!value) {
    return;
  }
  let indicator = document.createElement("span");
  indicator.className = "event-indicator";
  indicator.setAttribute("aria-label", label + ": " + value);
  indicator.title = label + ": " + value;
  indicator.innerHTML = icon === "activity"
    ? '<svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg><span>' + value + "</span>"
    : '<svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg><span>' + value + "</span>";
  card.getElementsByClassName("event-leading")[0].append(indicator);
}

function metadataIcon(icon, label) {
  if (icon === "activity") {
    return '<svg class="metadata-icon" viewBox="0 0 24 24" aria-label="' + label + '"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>';
  }
  return '<svg class="metadata-icon" viewBox="0 0 24 24" aria-label="' + label + '"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>';
}

function updateOptionalMetadata(content, field, label, labelClass, valueClass, value) {
  if (!value) {
    return;
  }

  let valueElement = content.querySelector("." + valueClass);
  let labelElement = content.querySelector("." + labelClass);
  if (valueElement === null || labelElement === null) {
    let metadata = content.querySelector(".drawer-meta");
    if (metadata === null) {
      return;
    }
    labelElement = document.createElement("dt");
    labelElement.className = labelClass;
    labelElement.textContent = label;
    let valueContainer = document.createElement("dd");
    valueContainer.className = valueClass;
    valueElement = document.createElement("span");
    valueElement.className = "drawer-metadata-value";
    valueContainer.append(valueElement);
    valueContainer.insertAdjacentHTML("beforeend", metadataIcon(field === "maxTBU" ? "activity" : "trash-2", label));
    metadata.append(labelElement, valueContainer);
  }

  let valueText = valueElement.querySelector(".drawer-metadata-value");
  if (valueText === null) {
    valueText = valueElement;
  }
  valueText.textContent = value;
  let hidden = value === "0s";
  labelElement.style.display = hidden ? "none" : "";
  valueElement.style.display = hidden ? "none" : "";
}

// keepalive
let lastKa;
function keepalive() {
  if (dashboardStale) {
    location.reload();
    return;
  }

  let ct = new Date().getTime();
  if (lastKa && lastKa + 60000 < ct) {
    location.reload();
  }
  lastKa = ct;
  let target = document.getElementById("status-bar");
  target.classList.remove("amber", "green", "grey", "noUpdate", "red");
  target.classList.add("green");
  target.getElementsByClassName("message")[0].innerHTML = "";
  target.getElementsByClassName("lastUpdated")[0].textContent = myTime();
  if (typeof ka !== "undefined") {
    clearTimeout(ka);
  }
  ka = setTimeout(function () {
    target.classList.remove("amber", "green", "grey", "noUpdate", "red");
    target.classList.add("noUpdate");
    target.getElementsByClassName("message")[0].innerHTML =
      "ERROR: No keepalives since " + myTime(lastKa) + ".";
    clearDashboardBoxes();
  }, 5 * 1000);
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
  z = z || "0";
  n = n + "";
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

// Change alert level of a box
function changeAlertLevel(target, status, message) {
  if (["amber", "green", "grey", "noUpdate", "red"].indexOf(status) === -1) {
    status = "grey";
  }
  target.classList.remove("amber", "green", "grey", "noUpdate", "red");
  target.classList.add(status);
  let tileMessage = target.getElementsByClassName("message")[0];
  if (tileMessage !== undefined) {
    tileMessage.innerHTML = message;
  }
  let tileLastUpdated = target.getElementsByClassName("lastUpdated")[0];
  if (tileLastUpdated !== undefined) {
    tileLastUpdated.innerHTML = myTime();
  }
  let pMessages = target.getElementsByClassName("previousMessages");
  if (typeof (pMessages[0] != "undefined") && pMessages[0] != null) {
    pMessages[0].insertAdjacentHTML(
      "afterbegin",
      "<li>" +
        myTime() +
        ": " +
        status.toUpperCase() +
        " (" +
        message +
        ")</li>",
    );
  }
}

// Make big box to fit as many biggest boxes as will fit the current window.
function rightSizeBigBox(pageType) {
  let widthBox;
  if (pageType === "dashboard") {
    let drawerWidth = 0;
    if (document.body.classList.contains("drawer-open")) {
      drawerWidth = document.getElementById("side-drawer-panel").getBoundingClientRect().width;
    }
    let availableWidth = window.innerWidth - drawerWidth - 30;
    widthBox = document.body.classList.contains("drawer-open")
      ? Math.max(0, availableWidth)
      : Math.floor(availableWidth / 512) * 512;
    if (widthBox < 512 && !document.body.classList.contains("drawer-open")) {
      widthBox = 512;
    }
  } else {
    widthBox = window.innerWidth - 30;
  }
  document.getElementById("big-box").style.width = widthBox + "px";
  let fullWidthBoxes = document.getElementsByClassName("fullwidth");
  for (let i = 0; i < fullWidthBoxes.length; i++) {
    fullWidthBoxes[i].style.width = widthBox - 2 + "px";
  }
}
