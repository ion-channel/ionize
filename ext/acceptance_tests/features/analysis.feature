Feature: Analysis
  As a user I want to be able to initiate an analysis of my project
  So that I can view the results

  @expected_failure
  Scenario: Analyze a invalid project
    Given I have an invalid project
    When I run an analysis
    Then the analysis should fail

  Scenario: Analyze a valid project with an invalid account
    Given I have an invalid account
    And I have a valid project
    When I run an analysis
    Then the analysis should fail

  @expected_failure
  Scenario: Analyze a valid project with a valid account
    Given I have a valid account
    And I have a valid project
    When I run an analysis
    Then the analysis should succeed
