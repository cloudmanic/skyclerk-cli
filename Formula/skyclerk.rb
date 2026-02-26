#
# Date: 2026-02-25
# Copyright (c) 2026. All rights reserved.
#

# Homebrew formula for the Skyclerk CLI. Downloads pre-built binaries
# from GitHub releases for the current platform and architecture.
class Skyclerk < Formula
  desc "CLI for the Skyclerk bookkeeping API"
  homepage "https://github.com/cloudmanic/skyclerk-cli"
  license "MIT"
  version "latest"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/latest/download/skyclerk-darwin-arm64"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/latest/download/skyclerk-darwin-amd64"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/latest/download/skyclerk-linux-arm64"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/latest/download/skyclerk-linux-amd64"
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "skyclerk"
  end

  test do
    assert_match "skyclerk version", shell_output("#{bin}/skyclerk version")
  end
end
