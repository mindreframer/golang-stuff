require "fileutils"

class VersionedFile
  def initialize(file, regex)
    @file = file
    @regex = regex
  end

  def current_version!
    @current_version ||= matched_data![1]
  end

  def bump_version!(type)
    position = case type
               when :major
                 0
               when :minor
                 1
               when :patch
                 2
               end
    @current_version = current_version!.split(".").tap do |v|
      v[position] = v[position].to_i + 1
      # Reset consequent numbers
      ((position + 1)..2).each { |p| v[p] = 0 }
    end.join(".")
  end

  def save!
    text = File.read(@file)
    new_line = matched_data![0].gsub(matched_data![1], @current_version)
    text.gsub!(matched_data![0], new_line)

    File.open(@file, "w") { |f| f.puts text }
  end

  private

  def matched_data!
    @matched_data ||= begin
                        m = @regex.match File.read(@file)
                        raise "No version #{@regex} matched in #{@file}" unless m
                        m
                      end
  end
end

def fullpath(file)
  File.expand_path(file, File.dirname(__FILE__))
end

VERSION_FILES = {
  fullpath("commands/version.go") => /^const Version = "(\d+.\d+.\d+)"$/,
  fullpath("README.md")           => /Current version is \[(\d+.\d+.\d+)\]/,
  fullpath(".goxc.json")          => /"PackageVersion": "(\d+.\d+.\d+)"/,
  fullpath("homebrew/gh.rb")      => /VERSION = "(\d+.\d+.\d+)"/
}

class Git
  class << self
    def dirty?
      !`git status -s`.empty?
    end

    def checkout
      `git checkout .`
    end

    def commit_all(msg)
      `git commit -am "#{msg}"`
    end

    def create_tag(tag, msg)
      `git tag -a #{tag} -m "#{msg}"`
    end
  end
end

namespace :release do
  desc "Current released version"
  task :current do
    vf = VersionedFile.new(*VERSION_FILES.first)
    puts vf.current_version!
  end

  [:major, :minor, :patch].each do |type|
    desc "Release #{type} version"
    task type do
      if Git.dirty?
        puts "Please commit all changes first"
        exit 1
      end

      new_versions = VERSION_FILES.map do |file, regex|
        begin
          vf = VersionedFile.new(file, regex)
          current_version = vf.current_version!
          vf.bump_version!(type)
          vf.save!
          puts "Successfully bump #{file} from #{current_version} to #{vf.current_version!}"
          vf.current_version!
        rescue => e
          Git.checkout
          raise e
        end
      end

      require "set"
      new_versions = new_versions.to_set
      if new_versions.size != 1
        raise "More than one version found among #{VERSION_FILES}"
      end

      new_version = "v#{new_versions.first}"
      msg = "Bump version to #{new_version}"
      Git.commit_all(msg)
      Git.create_tag(new_version, msg)
    end
  end
end

module OS
  class << self
    def type
      if darwin?
        "darwin"
      elsif linux?
        "linux"
      elsif windows?
        "windows"
      else
        raise "Unknown OS type #{RUBY_PLATFORM}"
      end
    end

    def dropbox_dir
      if darwin? || linux?
        File.join ENV["HOME"], "Dropbox"
      elsif windows?
        File.join ENV["DROPBOX_DIR"]
      else
        raise "Unknown OS type #{RUBY_PLATFORM}"
      end
    end

    def windows?
      (/cygwin|mswin|mingw|bccwin|wince|emx/ =~ RUBY_PLATFORM) != nil
    end

    def darwin?
      (/darwin/ =~ RUBY_PLATFORM) != nil
    end

    def linux?
      (/linux/ =~ RUBY_PLATFORM) != nil
    end
  end
end

namespace :build do
  desc "Build for current operating system"
  task :current => [:update_goxc, :remove_build_target, :build_gh, :move_to_dropbox]

  task :update_goxc do
    puts "Updating goxc..."
    result = system "go get -u github.com/laher/goxc"
    raise "Fail to update goxc" unless result
  end

  task :remove_build_target do
    FileUtils.rm_rf fullpath("target")
  end

  task :build_gh do
    puts "Building for #{OS.type}..."
    puts `goxc -wd=. -os=#{OS.type} -c=#{OS.type}`
  end

  task :move_to_dropbox do
    vf = VersionedFile.new(*VERSION_FILES.first)
    build_dir = fullpath("target/#{vf.current_version!}-snapshot")
    dropbox_dir = File.join(OS.dropbox_dir, "Public", "gh")

    FileUtils.cp_r build_dir, dropbox_dir, :verbose => true
  end
end
