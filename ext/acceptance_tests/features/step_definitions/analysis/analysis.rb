# Givens
Given(/^I have an invalid project$/) do
  @expected_project_id = 'bad-project-id'
end

Given(/^I have a valid project$/) do
  @expected_project_id = '7b9a0a87-fbe6-40c1-aa37-89acc6e5c191'
end

Given(/^I have an invalid account$/) do
  @expected_account_id = 'bad-account-id'
end

Given(/^I have a valid account$/) do
  @expected_account_id = '123456'
end

# Whens
When(/^I run an analysis$/) do
  @analysis = `ion-ci #{@expected_project_id} #{@expected_account_id} 1`.chomp
end

# Thens
Then(/^the analysis should fail$/) do
  expect(@analysis).to include('false')
end

Then(/^the analysis should succeed$/) do
  expect(@analysis).to include('true')
end
