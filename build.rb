#!/usr/bin/env ruby

require 'fileutils'

FR_VERSION = '0.0.1'.freeze
GO_VERSION = '1.9'.freeze

OSES_ARCHES = {
  'linux' => %w(amd64 386 arm arm64),
  'windows' => %w(amd64 386),
  'darwin' => %w(amd64)
}.freeze

def docker_run(os, arch)
  <<-EOS.gsub(/\s+/, ' ').gsub(/[\s\t]*\n/, ' ').strip
    docker run
    --rm
    --volume "#{Dir.pwd}:/usr/src/findref"
    --workdir "/usr/src/findref"
    --env GOOS=#{os}
    --env GOARCH=#{arch}
    golang:#{GO_VERSION} go build
  EOS
end

def main
  OSES_ARCHES.each do |os, arches|
    arches.each do |arch|
      dest_dir = "release/#{FR_VERSION}/#{os}/#{arch}"
      puts "Building findref v#{FR_VERSION} for #{os} #{arch}..."
      puts "Running: #{docker_run(os, arch)}"
      system(docker_run(os, arch))
      FileUtils.mkdir_p(dest_dir)
      fr = os == 'windows' ? 'findref.exe' : 'findref'
      FileUtils.mv(fr, "#{dest_dir}/")
    end
  end
end

main
