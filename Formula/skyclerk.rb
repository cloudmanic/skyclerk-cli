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
  version "0.1.9"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.9/skyclerk-darwin-arm64"
    sha256 "c75bdabfd21ded8af7c3687f8a60a567720fd63220c4996da7b8e0eddc0d82db"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.9/skyclerk-darwin-amd64"
    sha256 "97d3178b4daf581981606056b85e3232df3edd43ba91d064995b4a0a214ebe8a"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.9/skyclerk-linux-arm64"
    sha256 "b82cd7d83e5e78f38a2933b40a1ec739802b7950b4d27164819d4a1f98f8343a"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.9/skyclerk-linux-amd64"
    sha256 "90a11d6cf28b3f9ad959639d1c3dc88b8c1e2057abc807a306e9d47c68200bb3"
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "skyclerk"
  end

  test do
    assert_match "skyclerk version", shell_output("#{bin}/skyclerk version")
  end
end
