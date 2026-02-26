# Skyclerk CLI

A command-line interface for the [Skyclerk](https://app.skyclerk.com) bookkeeping API. Manage your ledger entries, categories, labels, contacts, reports, and more from the terminal.

## Installation

### Homebrew

```bash
brew tap cloudmanic/skyclerk-cli https://github.com/cloudmanic/skyclerk-cli
brew install skyclerk
```

### From Source

```bash
git clone https://github.com/cloudmanic/skyclerk-cli.git
cd skyclerk-cli
make build
```

The binary will be at `./skyclerk`. Move it to your `$PATH` or run `make install` to install to `$GOPATH/bin`.

## Authentication

Before using the CLI, you need to log in with your Skyclerk credentials. You'll need a valid OAuth client ID registered in your Skyclerk instance.

```bash
skyclerk login
```

You'll be prompted for:

1. **Client ID** - Your registered OAuth application client ID
2. **Email** - Your Skyclerk account email
3. **Password** - Your Skyclerk account password (hidden input)

If your account belongs to multiple Skyclerk accounts, you'll be prompted to select a default. Credentials are stored in `~/.config/skyclerk/config.json`.

To log out and revoke your token:

```bash
skyclerk logout
```

## Usage

All commands support `--output json` for machine-readable output and `--account <id>` to override the default account.

### Accounts

```bash
# List all accounts you belong to
skyclerk accounts list

# Show current account details
skyclerk accounts show

# Set the default account
skyclerk accounts use 42
```

### Ledger

```bash
# List ledger entries (default 25 per page)
skyclerk ledger list
skyclerk ledger list --limit 50 --page 2

# Get a single entry
skyclerk ledger get 12345

# Create an entry
skyclerk ledger create --amount -49.99 --date 2026-02-25 --contact-id 10 --category-id 5 --note "Office supplies"

# Update an entry
skyclerk ledger update 12345 --amount -59.99 --note "Updated note"

# Delete an entry
skyclerk ledger delete 12345

# View ledger summary (years, categories, labels with counts)
skyclerk ledger summary
```

### Categories

```bash
# List all categories
skyclerk categories list

# Get a single category
skyclerk categories get 5

# Create a category (type: expense or income)
skyclerk categories create --name "Office Supplies" --type expense

# Update a category
skyclerk categories update 5 --name "Office Equipment"

# Delete a category
skyclerk categories delete 5
```

### Labels

```bash
# List all labels
skyclerk labels list

# Get a single label
skyclerk labels get 3

# Create a label
skyclerk labels create --name "Q1 2026"

# Update a label
skyclerk labels update 3 --name "Q2 2026"

# Delete a label
skyclerk labels delete 3
```

### Contacts

```bash
# List all contacts
skyclerk contacts list
skyclerk contacts list --search "Acme"

# Get a single contact
skyclerk contacts get 10

# Create a contact
skyclerk contacts create --name "Acme Corp" --email "billing@acme.com" --phone "555-1234"

# Update a contact
skyclerk contacts update 10 --email "new@acme.com"

# Delete a contact
skyclerk contacts delete 10
```

### Files

```bash
# Upload a file
skyclerk files upload receipt.jpg

# Upload and associate with a ledger entry
skyclerk files upload receipt.jpg --ledger-id 12345
```

### Activities

```bash
# List recent account activities
skyclerk activities
skyclerk activities --limit 50
```

### Reports

```bash
# Profit and loss report
skyclerk reports pnl --start 2026-01-01 --end 2026-12-31

# Current year P&L
skyclerk reports pnl-current

# P&L by category
skyclerk reports pnl-by-category --start 2026-01-01 --end 2026-06-30

# P&L by label
skyclerk reports pnl-by-label --start 2026-01-01 --end 2026-06-30

# Income by contact
skyclerk reports income-by-contact --start 2026-01-01 --end 2026-12-31

# Expenses by contact
skyclerk reports expenses-by-contact --start 2026-01-01 --end 2026-12-31
```

### User Profile

```bash
# Show your profile
skyclerk me

# Update your profile
skyclerk me update --first-name "Jane" --last-name "Smith" --email "jane@example.com"
```

### User Management

```bash
# List users in the account
skyclerk users list

# Invite a user
skyclerk users invite --email "new@example.com" --first-name "New" --last-name "User"

# List pending invitations
skyclerk users invites

# Cancel an invitation
skyclerk users cancel-invite 5

# Remove a user from the account
skyclerk users remove 42
```

### Configuration

```bash
# Show current config (tokens are masked)
skyclerk config show

# Manually initialize config
skyclerk config init

# Show version
skyclerk version
```

## JSON Output

Every command supports `--output json` for scripting and AI agent integration:

```bash
skyclerk ledger list --output json
skyclerk categories list --output json
skyclerk reports pnl-current --output json
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--output` | Output format: `table` (default) or `json` |
| `--account` | Override the default account ID for this command |

## Development

```bash
# Run tests
make test

# Run tests without verbose output
make test-short

# Generate coverage report
make coverage

# Lint
make lint

# Build for all platforms
make cross-build

# Clean build artifacts
make clean
```

## License

MIT
