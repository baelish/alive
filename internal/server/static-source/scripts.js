/* jshint esversion: 6 */

// Dashboard module - encapsulates all dashboard functionality
const Dashboard = (function() {
  "use strict";

  // Constants
  const MAX_LOG_ENTRIES = 30;
  const KEEPALIVE_TIMEOUT = 5000;
  const DRAWER_REOPEN_TIMEOUT = 1000;

  // Cached DOM references
  const elements = {
    get tooltip() { return document.getElementById("tooltip"); },
    get bigBox() { return document.getElementById("big-box"); },
    get statusBar() { return document.getElementById("status-bar"); },
    get boxesGrid() { return document.getElementById("boxes-grid"); },
    get drawerContent() { return document.getElementById("drawer-content"); },
    get drawerPanel() { return document.getElementById("side-drawer-panel"); },
    get drawerBackdrop() { return document.getElementById("drawer-backdrop"); }
  };

  // Module state
  let activeBoxId = null;
  let resizingDrawer = false;
  let resizeStartX = 0;
  let resizeStartWidth = 0;
  let dashboardStale = false;
  let pendingDrawerReopenId = null;
  let lastKa;
  let ka;

  // Reload if re-visiting using back/forward buttons
  if (
    String(window.performance.getEntriesByType("navigation")[0].type) ===
    "back_forward"
  ) {
    location.reload();
  }

  // Register with box event source
  const source = new EventSource("/events/");

  source.onmessage = function (rawEvent) {
    const event = JSON.parse(rawEvent.data);
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
        break;
    }
  };

  source.onerror = function(error) {
    console.error("EventSource connection error:", error);
    if (elements.statusBar) {
      elements.statusBar.classList.remove("amber", "green", "grey", "noUpdate", "red");
      elements.statusBar.classList.add("noUpdate", "stale");
      const messageEl = elements.statusBar.getElementsByClassName("message")[0];
      if (messageEl) {
        messageEl.innerHTML = "Connection Lost";
      }
    }
  };

  // Box tooltip
  function boxHover(tip) {
    if (elements.tooltip) {
      elements.tooltip.innerHTML = tip;
      elements.tooltip.style.display = "block";
    }
  }

  function boxOut() {
    if (elements.tooltip) {
      elements.tooltip.innerHTML = "";
      elements.tooltip.style.display = "none";
    }
  }

  // Switch drawer tabs
  function switchDrawerTab(event, tabName) {
    const tabs = document.querySelectorAll(".drawer-tab");
    const tabContents = document.querySelectorAll(".drawer-tab-content");

    tabs.forEach(tab => tab.classList.remove("active"));
    tabContents.forEach(content => content.classList.remove("active"));

    event.currentTarget.classList.add("active");
    const targetTab = document.getElementById(`drawer-tab-${tabName}`);
    if (targetTab) {
      targetTab.classList.add("active");
    }
  }

  // Throttle helper for performance
  function throttle(func, wait) {
    let timeout;
    return function(...args) {
      if (!timeout) {
        timeout = setTimeout(() => {
          timeout = null;
          func.apply(this, args);
        }, wait);
      }
    };
  }

  // Initialize the dashboard drawer
  function initDashboard() {
    document.addEventListener("keydown", event => {
      if (event.key === "Escape") {
        closeDrawer();
      }
    });

    const resizeHandle = document.getElementById("drawer-resize-handle");
    if (!resizeHandle) return;

    const throttledResize = throttle(resizeDrawerToViewport, 100);
    window.addEventListener("resize", throttledResize);

    resizeHandle.addEventListener("pointerdown", event => {
      const panel = elements.drawerPanel;
      if (!panel) return;

      resizingDrawer = true;
      resizeStartX = event.clientX;
      resizeStartWidth = panel.getBoundingClientRect().width;
      resizeHandle.setPointerCapture(event.pointerId);
      document.body.classList.add("drawer-resizing");
      event.preventDefault();
    });

    resizeHandle.addEventListener("pointermove", event => {
      if (!resizingDrawer) return;

      const bounds = drawerWidthBounds();
      const width = Math.max(bounds.min, Math.min(bounds.max, resizeStartWidth + resizeStartX - event.clientX));
      if (elements.drawerPanel) {
        elements.drawerPanel.style.width = `${width}px`;
      }
      rightSizeBigBox("dashboard");
    });

    resizeHandle.addEventListener("pointerup", event => {
      resizingDrawer = false;
      resizeHandle.releasePointerCapture(event.pointerId);
      document.body.classList.remove("drawer-resizing");
    });

    sortAllBoxes();
  }

  function sortAllBoxes() {
    const grid = elements.boxesGrid;
    if (!grid) return;

    const sizeOrder = ["xlarge", "dlarge", "large", "dmedium", "medium", "dsmall", "small", "dmicro", "micro", "dot"];
    const boxElements = [];

    // Collect all box elements (every other element, since templates are interleaved)
    for (let i = 0; i < grid.children.length; i += 2) {
      const boxElement = grid.children[i];
      const detailElement = grid.children[i + 1];
      if (boxElement && boxElement.classList.contains("box")) {
        boxElements.push({
          box: boxElement,
          detail: detailElement,
          rank: boxSizeRank(boxElement, sizeOrder),
          name: boxElement.dataset.sortName || ""
        });
      }
    }

    // Sort by rank (size) first, then by name
    boxElements.sort((a, b) => {
      if (a.rank !== b.rank) {
        return a.rank - b.rank;
      }
      return a.name.localeCompare(b.name);
    });

    // Re-append in sorted order
    boxElements.forEach(item => {
      grid.appendChild(item.box);
      if (item.detail) {
        grid.appendChild(item.detail);
      }
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
  sortAllBoxes();
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
  sortAllBoxes();
}

function clearDashboardBoxes() {
  let grid = document.getElementById("boxes-grid");
  if (grid === null) {
    return;
  }
  grid.classList.add("stale");
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
    const title = box.displayName || box.name;

    const divContent = `
      <div onclick='boxClick(this.id)' onmouseover='boxHover("${box.name}")' onmouseout='boxOut()' id='${box.id}' class='${box.status} ${box.size} box'>
        <p class='title'>${title}</p>
        <p class='box-id'>${box.id}</p>
        <p class='maxTBU'>${box.maxTBU || ''}</p>
        <p class='expireAfter'>${box.expireAfter || ''}</p>
      </div>
    `;

    const grid = elements.boxesGrid;
    if (!grid) return;

    const precedingBox = document.getElementById(after);
    if (precedingBox) {
      precedingBox.insertAdjacentHTML("afterEnd", divContent);
    } else {
      grid.insertAdjacentHTML("beforeend", divContent);
    }

    const newBox = document.getElementById(box.id);
    if (!newBox) return;

    newBox.dataset.sortName = box.name;

    const oldDetail = document.getElementById(`box-detail-${box.id}`);
    if (oldDetail) {
      oldDetail.remove();
    }

  const detail = document.createElement("template");
  detail.id = `box-detail-${box.id}`;
  const maxTBU = hasMetadataValue(box.maxTBU) ? `<dt class="drawer-max-tbu-label">Max TBU</dt><dd class="drawer-max-tbu">${metadataIcon("activity", "Max TBU")}<span class="drawer-metadata-value">${box.maxTBU}</span></dd>` : '';
  const expireAfter = hasMetadataValue(box.expireAfter) ? `<dt class="drawer-expire-after-label">Expires after</dt><dd class="drawer-expire-after">${metadataIcon("trash-2", "Expires after")}<span class="drawer-metadata-value">${box.expireAfter}</span></dd>` : '';
  // Latest update indicators: only show if latest message has metadata
  const latestMessage = box.messages && box.messages.length > 0 ? box.messages[0] : null;
  const latestMaxTBU = latestMessage ? latestMessage.maxTBU : null;
  const latestExpireAfter = latestMessage ? latestMessage.expireAfter : null;
  const maxTBUIndicator = hasMetadataValue(latestMaxTBU) ? `<span class="event-indicator" aria-label="Max TBU: ${latestMaxTBU}" title="Max TBU: ${latestMaxTBU}"><svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg><span>${latestMaxTBU}</span></span>` : '';
  const expireAfterIndicator = hasMetadataValue(latestExpireAfter) ? `<span class="event-indicator" aria-label="Expires after: ${latestExpireAfter}" title="Expires after: ${latestExpireAfter}"><svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg><span>${latestExpireAfter}</span></span>` : '';
  const links = box.links && box.links.length ? '<div class="drawer-section drawer-links"><h3>Links</h3>' + box.links.map(link => {
    return `<a href="${link.url}" target="_blank" rel="noopener noreferrer"><svg class="external-link-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg><span>${link.name}</span></a><br />`;
  }).join("") + '</div>' : '';
  detail.innerHTML =
    '<div class="drawer-header-wrapper ' + box.status + '"><div class="drawer-header"><h2 class="drawer-title">' +
    title +
    '</h2><div class="drawer-actions"><a class="drawer-open-page" href="/box/' + encodeURIComponent(box.id) + '" target="_blank" rel="noopener noreferrer" aria-label="Open full box page" title="Open full box page"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 3h7v7"></path><path d="M10 14 21 3"></path><path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path></svg></a><button class="drawer-close" type="button" aria-label="Close details" onclick="closeDrawer()"><svg class="drawer-icon" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></svg></button></div></div></div>' +
    '<div class="drawer-tabs"><button class="drawer-tab active" onclick="switchDrawerTab(event, \'overview\')">Overview</button><button class="drawer-tab" onclick="switchDrawerTab(event, \'events\')">Events</button></div>' +
    '<div class="drawer-content-area"><div id="drawer-tab-overview" class="drawer-tab-content active">' +
    '<div class="drawer-section drawer-metadata"><h3>Status Details</h3><dl class="drawer-meta"><dt>ID</dt><dd><span class="drawer-metadata-value">' +
    box.id +
    '</span></dd><dt>Size</dt><dd><span class="drawer-metadata-value">' + box.size + '</span></dd>' +
    maxTBU + expireAfter +
    '</dl></div>' + links + '<div class="drawer-section drawer-latest-update"><h3>Latest update</h3><div class="event-log-card ' + box.status + '"><div class="event-top"><strong class="event-status">' + box.status + '</strong><div class="event-leading">' + maxTBUIndicator + expireAfterIndicator + '</div></div><div class="event-message">' +
    box.lastMessage +
    '</div><small>' + box.lastUpdate +
    '</small></div></div></div>' +
    '<div id="drawer-tab-events" class="drawer-tab-content"><div class="drawer-section drawer-events"><h3>Event log</h3><div class="event-log"></div></div></div></div>';
  const eventLog = detail.content.querySelector(".event-log");
  if (eventLog && box.messages) {
    box.messages.forEach(message => {
      const card = createEventCard(message.status, message.message, message.timeStamp);
      addEventIndicator(card, "activity", "Max TBU", message.maxTBU);
      addEventIndicator(card, "trash-2", "Expires after", message.expireAfter);
      eventLog.append(card);
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
  return sizeOrder.findIndex(size => boxElement.classList.contains(size));
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
    window.setTimeout(() => {
      if (pendingDrawerReopenId === id) {
        pendingDrawerReopenId = null;
        closeDrawer();
      }
    }, DRAWER_REOPEN_TIMEOUT);
  }
}

  // Update box
  function updateBox(event) {
    const targetBox = document.getElementById(event.id);

    if (targetBox) {
      changeAlertLevel(targetBox, event.status, event.lastMessage);
    }

    const detail = document.getElementById(`box-detail-${event.id}`);
    if (detail) {
      updateDetailContent(detail.content, event);
    }

    if (activeBoxId === event.id && elements.drawerContent) {
      updateDetailContent(elements.drawerContent, event);
    }

    if (!targetBox) return;

  if (hasMetadataValue(event.maxTBU)) {
    const t = targetBox.getElementsByClassName("maxTBU")[0];
    if (t) {
      if (t.tagName === "TR") {
        t.getElementsByTagName("TD")[0].innerHTML = event.maxTBU;
        t.style.display = "table-row";
      } else {
        t.innerHTML = event.maxTBU;
      }
    }
  }

  if (hasMetadataValue(event.expireAfter)) {
    const t = targetBox.getElementsByClassName("expireAfter")[0];
    if (t) {
      if (t.tagName === "TR") {
        t.getElementsByTagName("TD")[0].innerHTML = event.expireAfter;
        t.style.display = "table-row";
      } else {
        t.innerHTML = event.expireAfter;
      }
    }
  }
}

  function updateDetailContent(content, event) {
    const statusValue = event.status === "noUpdate" ? "red" : event.status || "grey";

    // Update drawer header wrapper with new status color
    const headerWrapper = content.querySelector(".drawer-header-wrapper");
    if (headerWrapper) {
      headerWrapper.classList.remove("amber", "green", "grey", "noUpdate", "red");
      headerWrapper.classList.add(statusValue);
    }

    // Update latest update section with event card style
    const latestUpdateSection = content.querySelector(".drawer-latest-update");
    if (latestUpdateSection) {
      const latestCard = latestUpdateSection.querySelector(".event-log-card");
      if (latestCard) {
        latestCard.classList.remove("amber", "green", "grey", "noUpdate", "red");
        latestCard.classList.add(statusValue);

        const statusLabel = latestCard.querySelector(".event-status");
        if (statusLabel) statusLabel.textContent = statusValue;

        const messageEl = latestCard.querySelector(".event-message");
        if (messageEl) messageEl.textContent = event.lastMessage;

        const timeEl = latestCard.querySelector("small");
        if (timeEl) timeEl.textContent = myTime();

      // Update indicators
      const eventLeading = latestCard.querySelector(".event-leading");
      if (eventLeading !== null) {
        eventLeading.innerHTML = '';
        if (hasMetadataValue(event.maxTBU)) {
          const maxTBUIndicator = document.createElement("span");
          maxTBUIndicator.className = "event-indicator";
          maxTBUIndicator.setAttribute("aria-label", `Max TBU: ${event.maxTBU}`);
          maxTBUIndicator.title = `Max TBU: ${event.maxTBU}`;
          maxTBUIndicator.innerHTML = `<svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg><span>${event.maxTBU}</span>`;
          eventLeading.appendChild(maxTBUIndicator);
        }
        if (hasMetadataValue(event.expireAfter)) {
          const expireAfterIndicator = document.createElement("span");
          expireAfterIndicator.className = "event-indicator";
          expireAfterIndicator.setAttribute("aria-label", `Expires after: ${event.expireAfter}`);
          expireAfterIndicator.title = `Expires after: ${event.expireAfter}`;
          expireAfterIndicator.innerHTML = `<svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg><span>${event.expireAfter}</span>`;
          eventLeading.appendChild(expireAfterIndicator);
        }
      }
    }
  }

    // Legacy support for old structure
    const dot = content.querySelector(".drawer-status-dot");
    const strip = content.querySelector(".drawer-status-strip");
    const message = content.querySelector(".drawer-message");
    const updated = content.querySelector(".drawer-last-updated");

    if (dot) {
      dot.classList.remove("amber", "green", "grey", "noUpdate", "red");
      dot.classList.add(statusValue);
    }
    if (strip) {
      strip.classList.remove("amber", "green", "grey", "noUpdate", "red");
      strip.classList.add(statusValue);
    }
    if (message) message.textContent = event.lastMessage;
    if (updated) updated.textContent = myTime();

    // Only update Status Details metadata if event explicitly changes them
    if (event.maxTBU !== undefined) {
      updateOptionalMetadata(content, "maxTBU", "Max TBU", "drawer-max-tbu-label", "drawer-max-tbu", event.maxTBU);
    }
    if (event.expireAfter !== undefined) {
      updateOptionalMetadata(content, "expireAfter", "Expires after", "drawer-expire-after-label", "drawer-expire-after", event.expireAfter);
    }

    const card = createEventCard(statusValue, event.lastMessage, null);
    addEventIndicator(card, "activity", "Max TBU", event.maxTBU);
    addEventIndicator(card, "trash-2", "Expires after", event.expireAfter);

    const log = content.querySelector(".event-log");
    if (log) {
      log.prepend(card);
      while (log.children.length > MAX_LOG_ENTRIES) {
        log.lastElementChild.remove();
      }
    }
  }

  function createEventCard(status, message, timestamp) {
    const statusValue = status === "noUpdate" ? "red" : status || "grey";
    const card = document.createElement("div");
    card.className = `event-log-card ${statusValue}`;
    card.innerHTML = '<div class="event-body"><div class="event-top"><strong class="event-status"></strong><div class="event-leading"></div></div><div class="event-message"></div><small></small></div>';
    card.getElementsByClassName("event-status")[0].textContent = statusValue;
    card.getElementsByClassName("event-message")[0].textContent = message;
    card.getElementsByTagName("small")[0].textContent = myTime(timestamp);
    return card;
  }

  function addEventIndicator(card, icon, label, value) {
    if (!hasMetadataValue(value)) return;

    const indicator = document.createElement("span");
    indicator.className = "event-indicator";
    indicator.setAttribute("aria-label", `${label}: ${value}`);
    indicator.title = `${label}: ${value}`;
    indicator.innerHTML = icon === "activity"
      ? `<svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg><span>${value}</span>`
      : `<svg class="event-icon" viewBox="0 0 24 24" aria-hidden="true"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg><span>${value}</span>`;

    const eventLeading = card.getElementsByClassName("event-leading")[0];
    if (eventLeading) {
      eventLeading.append(indicator);
    }
  }

  function metadataIcon(icon, label) {
    if (icon === "activity") {
      return `<svg class="metadata-icon" viewBox="0 0 24 24" aria-label="${label}"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>`;
    }
    return `<svg class="metadata-icon" viewBox="0 0 24 24" aria-label="${label}"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>`;
  }

  function hasMetadataValue(value) {
    // Allow any truthy value including "0s" and other duration strings
    return value !== null && value !== undefined && value !== "";
  }

  function updateOptionalMetadata(content, field, label, labelClass, valueClass, value) {
    let valueElement = content.querySelector(`.${valueClass}`);
    let labelElement = content.querySelector(`.${labelClass}`);

    if (!hasMetadataValue(value)) {
      if (valueElement) valueElement.remove();
      if (labelElement) labelElement.remove();
      return;
    }

    if (!valueElement || !labelElement) {
      const metadata = content.querySelector(".drawer-meta");
      if (!metadata) return;

      labelElement = document.createElement("dt");
      labelElement.className = labelClass;
      labelElement.textContent = label;

      const valueContainer = document.createElement("dd");
      valueContainer.className = valueClass;
      valueContainer.insertAdjacentHTML("beforeend", metadataIcon(field === "maxTBU" ? "activity" : "trash-2", label));

      valueElement = document.createElement("span");
      valueElement.className = "drawer-metadata-value";
      valueContainer.append(valueElement);
      metadata.append(labelElement, valueContainer);
    }

    const valueText = valueElement.querySelector(".drawer-metadata-value");
    if (valueText) {
      valueText.textContent = value;
    } else {
      valueElement.textContent = value;
    }
  }

  // Keepalive
  function keepalive() {
    if (dashboardStale) {
      location.reload();
      return;
    }

    const ct = Date.now();
    if (lastKa && lastKa + 60000 < ct) {
      location.reload();
    }
    lastKa = ct;

    const target = elements.statusBar;
    if (!target) return;

    target.classList.remove("amber", "green", "grey", "noUpdate", "red", "stale");
    target.classList.add("green");

    const messageEl = target.getElementsByClassName("message")[0];
    const lastUpdatedEl = target.getElementsByClassName("lastUpdated")[0];

    if (messageEl) messageEl.innerHTML = "Live";
    if (lastUpdatedEl) lastUpdatedEl.textContent = `Last update ${myTimeNoMillis()}`;

    if (ka) {
      clearTimeout(ka);
    }

    ka = setTimeout(() => {
      target.classList.remove("amber", "green", "grey", "noUpdate", "red");
      target.classList.add("noUpdate", "stale");
      if (messageEl) messageEl.innerHTML = "Disconnected";
      clearDashboardBoxes();
    }, KEEPALIVE_TIMEOUT);
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

// Print time without milliseconds (up to seconds only)
function myTimeNoMillis(t) {
  let r;
  if (t != null) {
    r = new Date(t);
  } else {
    r = new Date();
  }
  return r.toISOString().split('.')[0].replace('T', ' ');
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

  // Make big box to fit as many biggest boxes as will fit the current window
  function rightSizeBigBox(pageType) {
    const bigBox = elements.bigBox;
    if (!bigBox) return;

    let widthBox;
    if (pageType === "dashboard") {
      let drawerWidth = 0;
      if (document.body.classList.contains("drawer-open") && elements.drawerPanel) {
        drawerWidth = elements.drawerPanel.getBoundingClientRect().width;
      }
      // Reduce width to prevent status bar overflow (subtract extra margin)
      widthBox = document.documentElement.clientWidth - drawerWidth - 8;
    } else {
      widthBox = window.innerWidth - 30;
    }
    bigBox.style.width = `${widthBox}px`;
  }

  // Public API - expose functions that need to be called from HTML
  return {
    boxHover,
    boxOut,
    boxClick,
    closeDrawer,
    initDashboard,
    keepalive,
    rightSizeBigBox,
    switchDrawerTab
  };
})();

// Make functions available globally for HTML inline event handlers
window.boxHover = Dashboard.boxHover;
window.boxOut = Dashboard.boxOut;
window.boxClick = Dashboard.boxClick;
window.closeDrawer = Dashboard.closeDrawer;
window.initDashboard = Dashboard.initDashboard;
window.keepalive = Dashboard.keepalive;
window.rightSizeBigBox = Dashboard.rightSizeBigBox;
window.switchDrawerTab = Dashboard.switchDrawerTab;
