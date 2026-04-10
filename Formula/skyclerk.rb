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
  version "0.1.7"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.7/skyclerk-darwin-arm64"
    sha256 "367247263052a1f75d8e16e81db4fd0beeea8f64007240f0df270606163ba9f7"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.7/skyclerk-darwin-amd64"
    sha256 "caae6c2836a33cf348ce225859788536d3faf971ebb1f8fb8886e5d9c597797a"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.7/skyclerk-linux-arm64"
    sha256 "ed129c318ad2826d4562010bbeb503fdbb19f4a1685f5e3fc813af27266f94b7"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.7/skyclerk-linux-amd64"
    sha256 "ec24ea2203437ea38ea43a016338bbef0e8ee1ef6e5cca616c6c8c0d1bf3bed9"
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "skyclerk"
  end

  test do
    assert_match "skyclerk version", shell_output("#{bin}/skyclerk version")
  end
end
