#
# Date: 2026-02-25
# Copyright (c) 2026. All rights reserved.
#

# Homebrew formula for the Skyclerk CLI. Downloads pre-built binaries
# from GitHub releases for the current platform and architecture.
# This formula is auto-updated by the release workflow.
class Skyclerk < Formula
  desc "CLI for the Skyclerk bookkeeping API"
  homepage "https://github.com/cloudmanic/skyclerk-cli"
  license "MIT"
  version "0.1.6"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.6/skyclerk-darwin-arm64"
    sha256 "0a80d3b35ba993a06326e986ebebecffb32770af3c9235895b068088930ae3c1"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.6/skyclerk-darwin-amd64"
    sha256 "6ed85f36bd4b5c16016dcac00c1240fa8fdf48c2b355d54dbae85a3141276c01"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.6/skyclerk-linux-arm64"
    sha256 "2314d5bbd619bf8c9281784a38625ae213a99443761a0926931ed81e4e36e4e8"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.6/skyclerk-linux-amd64"
    sha256 "47f7d896380863ab740550056ddc341ac2f51a898e0d743d7c99b9021db4f0f1"
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "skyclerk"
  end

  test do
    assert_match "skyclerk version", shell_output("#{bin}/skyclerk version")
  end
end
