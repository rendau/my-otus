Feature: Empty events
  Scenario: List events for month
    When I request event list for month
    Then I will receive 0 event counts in response

  Scenario: List events for week
    When I request event list for week
    Then I will receive 0 event counts in response

  Scenario: List events for day
    When I request event list for day
    Then I will receive 0 event counts in response
