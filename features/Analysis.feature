Feature: Analysis
  Scenario: Analyze a project that doesnt exist
    Given an Ion Channel account id account_id
    Given an Ion Channel project id ASDF
    Given an Ion Channel build number 1
    When I run ion-cli
    Then the ion output should contain:
    """
    false
    """
  Scenario: Analyze a project with an account_id that doesnt exist
    Given an Ion Channel account id 123456
    Given an Ion Channel project id 7b9a0a87-fbe6-40c1-aa37-89acc6e5c191
    Given an Ion Channel build number 1
    When I run ion-cli
    Then the ion output should contain:
    """
    false
    """
  Scenario: Analyze a project with an account_id that does exist
    Given an Ion Channel account id account_id
    Given an Ion Channel project id 7b9a0a87-fbe6-40c1-aa37-89acc6e5c191
    Given an Ion Channel build number 1
    When I run ion-cli
    Then the ion output should contain:
    """
    true
    """
