const { Before, After, setWorldConstructor, setDefaultTimeout } = require("@cucumber/cucumber");
const { chromium } = require("playwright");

const stepTimeoutMs = Number(process.env.STEP_TIMEOUT_MS || "10000");
setDefaultTimeout(Number.isNaN(stepTimeoutMs) ? 10000 : stepTimeoutMs);

function launchOptions() {
  const headless = process.env.HEADLESS !== "false";
  const slowMo = Number(process.env.SLOWMO_MS || "0");
  return {
    headless,
    slowMo: Number.isNaN(slowMo) ? 0 : slowMo
  };
}

class CustomWorld {
  constructor() {
    this.browser = null;
    this.page = null;
    this.baseUrl = process.env.BASE_URL || "http://localhost:8080";
  }
}

setWorldConstructor(CustomWorld);

Before(async function () {
  this.browser = await chromium.launch(launchOptions());
  this.page = await this.browser.newPage();
});

After(async function () {
  if (this.browser) {
    await this.browser.close();
  }
});
