Feature: Empty events
  Scenario: List events for month
    When I request event list for month
    Then I will receive 1 event counts in response

  Scenario: List events for week
    When I request event list for week
    Then I will receive 0 event counts in response

  Scenario: List events for day
    When I request event list for day
    Then I will receive 0 event counts in response

  Scenario: Event create without owner
    When I create event with data:
		"""
		{
			"title": "event title",
			"text": "event text",
			"start_time": "2021-01-10T14:00:00Z",
			"end_time": "2021-01-10T15:00:00Z"
		}
		"""
    Then The response code should be 200
    And I receive data with error code = "owner is required"

  Scenario: Event create without title
    When I create event with data:
		"""
		{
			"owner": "event owner",
			"text": "event text",
			"start_time": "2021-01-10T14:00:00Z",
			"end_time": "2021-01-10T15:00:00Z"
		}
		"""
    Then The response code should be 200
    And I receive data with error code = "title is required"

  Scenario: Event create without start-time
    When I create event with data:
		"""
		{
			"owner": "event owner",
			"title": "event title",
			"text": "event text",
			"end_time": "2021-01-10T15:00:00Z"
		}
		"""
    Then The response code should be 200
    And I receive data with error code = "start date is required"

  Scenario: Event create without end-time
    When I create event with data:
		"""
		{
			"owner": "event owner",
			"title": "event title",
			"text": "event text",
			"start_time": "2021-01-10T14:00:00Z"
		}
		"""
    Then The response code should be 200
    And I receive data with error code = "end date is required"

  Scenario: Event create with incorrect start time
    When I create event with data:
		"""
		{
			"owner": "event owner",
			"title": "event title",
			"text": "event text",
			"start_time": "2010-01-10T14:00:00Z",
			"end_time": "2021-01-10T15:00:00Z"
		}
		"""
    Then The response code should be 200
    And I receive data with error code = "start_date is incorrect"

  Scenario: Event create with incorrect period
    When I create event with data:
		"""
		{
			"owner": "event owner",
			"title": "event title",
			"text": "event text",
			"start_time": "2021-01-10T14:00:00Z",
			"end_time": "2021-01-10T13:00:00Z"
		}
		"""
    Then The response code should be 200
    And I receive data with error code = "end_date is less than start_date"

  Scenario: Event create with correct data
    When I create event with data:
		"""
		{
			"owner": "event owner",
			"title": "event title",
			"text": "event text",
			"start_time": "2021-01-10T14:00:00Z",
			"end_time": "2021-01-10T15:00:00Z"
		}
		"""
    Then The response code should be 200
    And I receive data with error code = ""

  Scenario: Event create for coming day
    When I create event for coming day, with title "event for day"
    Then The response code should be 200
    And I receive data with error code = ""

  Scenario: Event create for coming week
    When I create event for coming week, with title "event for week"
    Then The response code should be 200
    And I receive data with error code = ""

  Scenario: Event create for coming month
    When I create event for coming month, with title "event for month"
    Then The response code should be 200
    And I receive data with error code = ""

  Scenario: List events for day
    When I request event list for day
    Then I will receive 1 event counts in response
    And The response will contain event with title "event for day"

  Scenario: List events for week
    When I request event list for week
    Then I will receive 2 event counts in response
    And The response will contain event with title "event for day"
    And The response will contain event with title "event for week"

  Scenario: List events for month
    When I request event list for month
    Then I will receive 3 event counts in response
    And The response will contain event with title "event for day"
    And The response will contain event with title "event for week"
    And The response will contain event with title "event for month"
