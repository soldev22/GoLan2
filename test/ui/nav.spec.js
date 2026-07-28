const { chromium } = require("playwright");
const { expect } = require("chai");

function launchOptions() {
  const headless = process.env.HEADLESS !== "false";
  const slowMo = Number(process.env.SLOWMO_MS || "0");
  return {
    headless,
    slowMo: Number.isNaN(slowMo) ? 0 : slowMo
  };
}

describe("Navigation UI", function () {
  this.timeout(10000);

  let browser;
  let context;
  let page;
  const baseUrl = process.env.BASE_URL || "http://localhost:8080";

  before(async () => {
    browser = await chromium.launch(launchOptions());
    context = await browser.newContext({ viewport: { width: 390, height: 844 } });
    page = await context.newPage();
  });

  after(async () => {
    await context.close();
    await browser.close();
  });

  it("toggles nav expanded state", async () => {
    await page.goto(baseUrl + "/");
    await page.click(".nav-toggle");
    const expanded = await page.getAttribute(".nav-toggle", "aria-expanded");
    expect(expanded).to.equal("true");
  });
});
