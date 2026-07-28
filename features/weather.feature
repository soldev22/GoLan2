Feature: Weather search

  Scenario: User searches for a city from home page
    Given I open the home page
    When I search weather for "Edinburgh"
    Then I should see the weather page
    And the city input should contain "Edinburgh"
