require 'rspec'
When(/^I run ion\-cli$/) do
  @output = `./main #{@project_id} #{@account_id} #{@build_id}`.chomp
end

Then(/^the ion output should contain:$/) do |string|
  expect(@output).to include(string)
end
Given(/^an Ion Channel account id (\S+)$/) do |account_id|
  @account_id = account_id
end
Given(/^an Ion Channel project id (\S+)$/) do |project_id|
  @project_id = project_id
end
Given(/^an Ion Channel build number (\S+)$/) do |build_id|
  @build_id = build_id
end
