require 'rspec/expectations'

def env_or_default key, default
  if ENV.has_key?(key)
    return ENV[key]
  else
    return default
  end
end

World(RSpec::Matchers)
