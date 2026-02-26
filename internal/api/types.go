// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

// LoginRequest represents the payload for the OAuth token endpoint.
type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	GrantType string `json:"grant_type"`
	ClientID  string `json:"client_id"`
}

// LoginResponse represents the response from the OAuth token endpoint.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserID      uint   `json:"user_id"`
	TokenType   string `json:"token_type"`
}

// Account represents a Skyclerk account.
type Account struct {
	ID           uint   `json:"id"`
	OwnerID      uint   `json:"owner_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Country      string `json:"country"`
	Locale       string `json:"locale"`
	Currency     string `json:"currency"`
	LastActivity string `json:"last_activity"`
}

// User represents a Skyclerk user.
type User struct {
	ID           uint      `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	LastActivity string    `json:"last_activity"`
	Accounts     []Account `json:"accounts"`
}

// MeResponse represents the current user profile.
type MeResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

// MeUpdateRequest represents the payload for updating the current user.
type MeUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// ChangePasswordRequest represents the payload for changing the user password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// Ledger represents a financial transaction entry.
type Ledger struct {
	ID        uint     `json:"id"`
	AccountID uint     `json:"account_id"`
	AddedByID uint     `json:"added_by_id"`
	Amount    float64  `json:"amount"`
	Date      string   `json:"date"`
	Contact   Contact  `json:"contact"`
	Category  Category `json:"category"`
	Labels    []Label  `json:"labels"`
	Files     []File   `json:"files"`
	Note      string   `json:"note"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// LedgerCreateRequest represents the payload for creating a ledger entry.
type LedgerCreateRequest struct {
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	ContactID  uint    `json:"contact_id"`
	CategoryID uint    `json:"category_id"`
	LabelIDs   []uint  `json:"label_ids,omitempty"`
	Note       string  `json:"note,omitempty"`
}

// LedgerUpdateRequest represents the payload for updating a ledger entry.
type LedgerUpdateRequest struct {
	Amount     float64 `json:"amount,omitempty"`
	Date       string  `json:"date,omitempty"`
	ContactID  uint    `json:"contact_id,omitempty"`
	CategoryID uint    `json:"category_id,omitempty"`
	LabelIDs   []uint  `json:"label_ids,omitempty"`
	Note       string  `json:"note,omitempty"`
}

// LedgerSummary represents a summary of ledger data grouped by year, label, and category.
type LedgerSummary struct {
	Years      []LedgerSummaryYear     `json:"years"`
	Labels     []LedgerSummaryItem     `json:"labels"`
	Categories []LedgerSummaryItem     `json:"categories"`
}

// LedgerSummaryYear represents a year entry in the ledger summary.
type LedgerSummaryYear struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

// LedgerSummaryItem represents a label or category entry in the ledger summary.
type LedgerSummaryItem struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Category represents a transaction category.
type Category struct {
	ID        uint   `json:"id"`
	AccountID uint   `json:"account_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Count     int    `json:"count"`
}

// CategoryCreateRequest represents the payload for creating a category.
type CategoryCreateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// CategoryUpdateRequest represents the payload for updating a category.
type CategoryUpdateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Label represents a transaction label/tag.
type Label struct {
	ID        uint   `json:"id"`
	AccountID uint   `json:"account_id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
}

// LabelCreateRequest represents the payload for creating a label.
type LabelCreateRequest struct {
	Name string `json:"name"`
}

// LabelUpdateRequest represents the payload for updating a label.
type LabelUpdateRequest struct {
	Name string `json:"name"`
}

// Contact represents a business contact.
type Contact struct {
	ID            uint   `json:"id"`
	AccountID     uint   `json:"account_id"`
	Name          string `json:"name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Fax           string `json:"fax"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zip           string `json:"zip"`
	Country       string `json:"country"`
	Website       string `json:"website"`
	AccountNumber string `json:"account_number"`
}

// ContactCreateRequest represents the payload for creating a contact.
type ContactCreateRequest struct {
	Name          string `json:"name"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	Zip           string `json:"zip,omitempty"`
	Country       string `json:"country,omitempty"`
	Website       string `json:"website,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
}

// ContactUpdateRequest represents the payload for updating a contact.
type ContactUpdateRequest struct {
	Name          string `json:"name,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Address       string `json:"address,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	Zip           string `json:"zip,omitempty"`
	Country       string `json:"country,omitempty"`
	Website       string `json:"website,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
}

// File represents an uploaded file/receipt.
type File struct {
	ID              uint   `json:"id"`
	AccountID       uint   `json:"account_id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Size            int64  `json:"size"`
	URL             string `json:"url"`
	Thumb600By600   string `json:"thumb_600_by_600_url"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// Activity represents an account activity log entry.
type Activity struct {
	ID         uint    `json:"id"`
	AccountID  uint    `json:"account_id"`
	UserID     uint    `json:"user_id"`
	LedgerID   uint    `json:"ledger_id"`
	ContactID  uint    `json:"contact_id"`
	LabelID    uint    `json:"label_id"`
	CategoryID uint    `json:"category_id"`
	Action     string  `json:"action"`
	SubAction  string  `json:"sub_action"`
	Amount     float64 `json:"amount"`
	Message    string  `json:"message"`
	CreatedAt  string  `json:"created_at"`
}

// PnlReport represents a profit and loss report.
type PnlReport struct {
	Income    float64        `json:"income"`
	Expense   float64        `json:"expense"`
	Profit    float64        `json:"profit"`
	Breakdown []PnlBreakdown `json:"breakdown,omitempty"`
}

// PnlBreakdown represents a line item in a P&L report.
type PnlBreakdown struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// Invite represents a pending user invitation.
type Invite struct {
	ID        uint   `json:"id"`
	AccountID uint   `json:"account_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Message   string `json:"message"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

// InviteCreateRequest represents the payload for inviting a user.
type InviteCreateRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Message   string `json:"message,omitempty"`
}

// Billing represents account billing information.
type Billing struct {
	ID                 uint   `json:"id"`
	PaymentProcessor   string `json:"payment_processor"`
	Subscription       string `json:"subscription"`
	Status             string `json:"status"`
	TrialExpire        string `json:"trial_expire"`
	CardBrand          string `json:"card_brand"`
	CardLast4          string `json:"card_last4"`
	CardExpMonth       int    `json:"card_exp_month"`
	CardExpYear        int    `json:"card_exp_year"`
	CurrentPeriodStart string `json:"current_period_start"`
	CurrentPeriodEnd   string `json:"current_period_end"`
}
