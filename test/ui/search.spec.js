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

describe("Search form UI", function () {
  this.timeout(10000);

  let browser;
  let page;
  const baseUrl = process.env.BASE_URL || "http://localhost:8080";

  before(async () => {
    browser = await chromium.launch(launchOptions());
    page = await browser.newPage();
  });

  after(async () => {
    await browser.close();
  });

  it("trims city value and disables submit button on submit", async () => {
    await page.goto(baseUrl + "/");

    await page.evaluate(() => {
      const form = document.querySelector(".search-form");
      if (form) {
        form.addEventListener("submit", (event) => {
          event.preventDefault();
        });
      }
    });

    await page.fill('input[name="city"]', "  Edinburgh  ");
    await page.click('button[type="submit"]');

    const value = await page.inputValue('input[name="city"]');
    const disabled = await page.isDisabled('button[type="submit"]');

    expect(value).to.equal("Edinburgh");
    expect(disabled).to.equal(true);
  });
});
