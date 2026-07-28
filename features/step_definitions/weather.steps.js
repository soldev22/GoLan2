const { Given, When, Then } = require("@cucumber/cucumber");
const { expect } = require("chai");

Given("I open the home page", async function () {
  await this.page.goto(this.baseUrl + "/");
});

When("I search weather for {string}", async function (city) {
  await this.page.fill('input[name="city"]', city);
  await this.page.click('button[type="submit"]');
});

Then("I should see the weather page", async function () {
  await this.page.waitForURL(/\/weather(\?|$)/);
  const panel = this.page.locator('section[aria-label="weather content"]');
  expect(await panel.count()).to.be.greaterThan(0);
});

Then("the city input should contain {string}", async function (city) {
  const value = await this.page.inputValue('input[name="city"]');
  expect(value).to.equal(city);
});
