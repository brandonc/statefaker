#!/usr/bin/env ruby
require 'json'
require 'open3'
require 'fileutils'
require 'optparse'

# Configuration defaults
options = {
  cwd: Dir.pwd,
  iterations: 5,
  statefaker_bin: nil,
  resources: 1000,
  outputs: 0,
  verbose: false
}

OptionParser.new do |opts|
  opts.banner = "Usage: push_multiple_states.rb [options]"

  opts.on("-c", "--cwd DIR", "Directory containing Terraform config (default: current dir)") do |v|
    options[:cwd] = v
  end

  opts.on("-n", "--iterations N", Integer, "Number of state versions to push (default: 5)") do |v|
    options[:iterations] = v
  end

  opts.on("-b", "--bin PATH", "Path to statefaker binary (default: searches PATH/cwd)") do |v|
    options[:statefaker_bin] = v
  end

  opts.on("-r", "--resources N", Integer, "Number of resources to generate per state (default: 1000)") do |v|
    options[:resources] = v
  end

  opts.on("-o", "--outputs N", Integer, "Number of outputs to generate per state (default: 0)") do |v|
    options[:outputs] = v
  end

  opts.on("-v", "--verbose", "Run verbosely") do |v|
    options[:verbose] = v
  end
end.parse!

# Helper to find executable in PATH
def find_executable(name)
  ENV['PATH'].split(File::PATH_SEPARATOR).each do |path|
    exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(';') : ['']
    exts.each do |ext|
      exe = File.join(path, "#{name}#{ext}")
      return exe if File.executable?(exe) && !File.directory?(exe)
    end
  end
  nil
end

# 1. Locate statefaker binary
statefaker = options[:statefaker_bin]
if statefaker.nil?
  # specific to this repo structure
  possible_paths = [
    "./statefaker",
    "../statefaker",
    "../../statefaker",
    File.join(File.dirname(__FILE__), "../statefaker")
  ]
  
  statefaker = possible_paths.find { |p| File.exist?(p) }
  
  if statefaker.nil?
    # Try to find in PATH
    statefaker = find_executable("statefaker")
  end

  # If still not found, try to build it if we are in the source repo
  if statefaker.nil? && File.exist?("go.mod") && File.read("go.mod").include?("module github.com/brandonc/statefaker")
    puts "Building statefaker..."
    system("go build -o statefaker ./cmd/statefaker")
    statefaker = "./statefaker"
  end
end

if statefaker.nil? || !File.exist?(statefaker)
  puts "Error: Could not find 'statefaker' binary."
  puts "Please build it with `go build -o statefaker ./cmd/statefaker` or provide path with --bin"
  exit 1
end

statefaker = File.expand_path(statefaker)
puts "Using statefaker binary: #{statefaker}" if options[:verbose]

# 2. Get initial state from Terraform
# We need to run this in the target directory (cwd)
target_cwd = File.expand_path(options[:cwd])
unless File.directory?(target_cwd)
  puts "Error: Target directory #{target_cwd} does not exist."
  exit 1
end

puts "Inspecting remote state in #{target_cwd}..."
stdout, stderr, status = Open3.capture3("terraform", "state", "pull", :chdir => target_cwd)

unless status.success?
  puts "Error: Failed to pull Terraform state. Ensure you are authenticated and initialized in #{target_cwd}."
  puts stderr
  exit 1
end

begin
  remote_state = JSON.parse(stdout)
rescue JSON::ParserError
  puts "Error: invalid JSON from terraform state pull"
  exit 1
end

current_serial = remote_state["serial"]
lineage = remote_state["lineage"]

if current_serial.nil? || lineage.nil?
  puts "Error: Could not determine serial or lineage from remote state."
  exit 1
end

puts "Current Remote Serial: #{current_serial}"
puts "Remote Lineage: #{lineage}"
puts "Plan: Push #{options[:iterations]} new versions starting at serial #{current_serial + 1}"

# 3. Loop and push
temp_file = File.join(Dir.pwd, "temp_huge.tfstate")

(1..options[:iterations]).each do |i|
  # Fetch latest serial to handle remote updates/migrations
  stdout, stderr, status = Open3.capture3("terraform", "state", "pull", :chdir => target_cwd)
  if status.success?
    remote_state = JSON.parse(stdout)
    current_serial = remote_state["serial"]
  end
  
  next_serial = current_serial + 1
  puts "\n--- Iteration #{i}/#{options[:iterations]} (Target Serial: #{next_serial}) ---"

  # Generate
  # We run statefaker where the binary is, or use absolute path
  # It writes to temp_file
  cmd = "#{statefaker} -resources #{options[:resources]} -outputs #{options[:outputs]} > #{temp_file}"
  puts "Generating state with #{cmd}..." if options[:verbose]
  
  unless system(cmd)
    puts "Error running statefaker generation command"
    exit 1
  end

  # Patch
  puts "Patching serial/lineage..." if options[:verbose]
  json_data = File.read(temp_file)
  state = JSON.parse(json_data)
  
  # Update critical fields
  state["version"] = 4
  state["serial"] = next_serial
  state["lineage"] = lineage
  
  File.write(temp_file, JSON.pretty_generate(state))

  # Push
  puts "Pushing to remote workspace..."
  # We must run `terraform state push` inside the directory with the backend config
  # But the file path must be accessible. `temp_file` is absolute path so it's fine.
  stdout, stderr, status = Open3.capture3("terraform", "state", "push", temp_file, :chdir => target_cwd)

  if status.success?
    puts "✅ Successfully pushed serial #{next_serial}"
  else
    puts "❌ Failed to push serial #{next_serial}"
    # Check for stale serial error (race condition)
    if stderr.include?("stale") || stderr.include?("higher serial")
      puts "⚠️  Remote serial changed underneath us. Re-fetching..."
      # In a real generalized script, we might want to loop and retry here.
      # For now, just exiting.
    end
    puts "STDOUT: #{stdout}"
    puts "STDERR: #{stderr}"
    exit 1
  end
  
  # Optional: sleep to not hammer the API too hard
  sleep 1
end

# Cleanup
File.delete(temp_file) if File.exist?(temp_file)
puts "\nDone! Pushed #{options[:iterations]} state versions."
