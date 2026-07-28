(function () {
  "use strict";

  var forms = document.querySelectorAll(".search-form");

  forms.forEach(function (form) {
    var input = form.querySelector('input[name="city"]');
    var button = form.querySelector('button[type="submit"]');
    var loadingLabel = form.getAttribute("data-loading-label") || "Loading...";

    form.addEventListener("submit", function () {
      if (input) {
        input.value = input.value.trim();
      }

      if (button) {
        button.disabled = true;
        button.textContent = loadingLabel;
      }
    });
  });
})();
