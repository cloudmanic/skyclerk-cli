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
  version "0.1.5"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.5/skyclerk-darwin-arm64"
    sha256 "29cac8f309769c5022ca7de944b9085b8223c54d81dfb7d29d5ed0882e865c54"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.5/skyclerk-darwin-amd64"
    sha256 "469e3f0c5abaa40e405793386a45c9b884d36dce17f83cb75e539abbd53133f2"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.5/skyclerk-linux-arm64"
    sha256 "4d35123f165dfe4fedc956b4ab104b526fe2ef93b9e5f293a5b61c5118abf03d"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/skyclerk-cli/releases/download/v0.1.5/skyclerk-linux-amd64"
    sha256 "af2569f09b1533d6afa2218561dc94354256a90089a18b01bccd1e6bbb649d2c"
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "skyclerk"
  end

  test do
    assert_match "skyclerk version", shell_output("#{bin}/skyclerk version")
  end
end
