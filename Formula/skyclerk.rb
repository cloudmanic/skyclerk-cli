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
  version "0.1.8"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.8/skyclerk-darwin-arm64"
    sha256 "ba52102f393833d39e00acc7f49e7723d4e8287936c80cdbe484a3eee28f02fc"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.8/skyclerk-darwin-amd64"
    sha256 "533ba76ebff45f3c48191357737c0c795c7e77146619b0dc50c30c7d41f08a72"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.8/skyclerk-linux-arm64"
    sha256 "6042d632fe735bdd97efeb44be6aa5991ca0f91ed41f60ce2e077dff79313c78"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.8/skyclerk-linux-amd64"
    sha256 "ce88721b134b73395c0343283fb4dee4e04a57ebca03df8e9caaa998932d83a2"
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "skyclerk"
  end

  test do
    assert_match "skyclerk version", shell_output("#{bin}/skyclerk version")
  end
end
