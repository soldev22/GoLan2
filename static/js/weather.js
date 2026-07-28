(function () {
  "use strict";

  var forecastDates = document.querySelectorAll("[data-date]");

  forecastDates.forEach(function (node) {
    var rawDate = node.getAttribute("data-date");

    if (!rawDate) {
      return;
    }

    var parsed = new Date(rawDate + "T00:00:00");
    if (Number.isNaN(parsed.getTime())) {
      return;
    }

    var label = new Intl.DateTimeFormat([], {
      weekday: "short",
      month: "short",
      day: "numeric"
    }).format(parsed);

    node.textContent = label;
  });
})();
