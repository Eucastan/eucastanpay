package events

// Events Struct
type UserRegisteredEvent struct {
	EventMetadata
	UserID    string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type UserRegistrationFailedEvent struct {
	EventMetadata
	UserID string `json:"id"`
	Reason string `json:"reason"`
}

type KYCCreatedEvent struct {
	EventMetadata
	UserID    string `json:"user_id"`
	KYCStatus string `json:"kyc_status"`
	Timestamp int64  `json:"timestamp"`
}

type UserKYCVerifiedEvent struct {
	EventMetadata
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	KYCStatus string `json:"kyc_status"`
	Timestamp int64  `json:"timestamp"`
}

type AccountCreatedEvent struct {
	EventMetadata
	AccountID   string `json:"id"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	AccountNo   int64  `json:"account_no"`
	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`
	Timestamp   int64  `json:"timestamp"`
}

type CreateAccRequestEvent struct {
	AccountID   string `json:"id"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	AccountNo   int64  `json:"account_no"`
	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`
	Timestamp   int64  `json:"timestamp"`
}

type DepositAccountEvent struct {
	AccountID    string `json:"id"`
	UserID       string `json:"user_id"`
	Amount       int64  `json:"amount"`
	AccountNo    int64  `json:"account_no"`
	AccountType  string `json:"account_type"`
	Reference    string `json:"reference"`
	BalanceAfter int64  `json:"balance_after"`
	Currency     string `json:"currency"`
	Timestamp    int64  `json:"timestamp"`
}

type WithdrawalEvent struct {
	AccountID    string `json:"id"`
	UserID       string `json:"user_id"`
	Amount       int64  `json:"amount"`
	AccountNo    int64  `json:"account_no"`
	AccountType  string `json:"account_type"`
	Reference    string `json:"reference"`
	BalanceAfter int64  `json:"balance_after"`
	Currency     string `json:"currency"`
	Timestamp    int64  `json:"timestamp"`
}

type CreateAccFailedEvent struct {
	EventMetadata
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Reason string `json:"reason"`
}

type DebitCompletedEvent struct {
	EventMetadata
	UserID           string `json:"user_id"`
	FromAccID        string `json:"from_account_id"`
	Reference        string `json:"reference"`
	Amount           int64  `json:"amount"`
	FromBalanceAfter int64  `json:"from_balance_after"`
	Timestamp        int64  `json:"timestamp"`
}

type DebitFailedEvent struct {
	EventMetadata
	UserID    string `json:"user_id"`
	Reference string `json:"reference"`
	Reason    string `json:"reason"`
}

type ReverseInitiatedEvent struct {
	EventMetadata
	UserID     string `json:"user_id"`
	TransferID string `json:"transfer_id"`
	Reference  string `json:"reference"`
	FromAccID  string `json:"from_account_id"`
	FromAccNo  int64  `json:"from_account_no"`
	ToAccID    string `json:"to_account_id"`
	ToAccNo    int64  `json:"to_account_no"`
	Amount     int64  `json:"amount"`
	Timestamp  int64  `json:"timestamp"`
}

type ReverseFailedTransferEvent struct {
	EventMetadata
	UserID    string `json:"user_id"`
	Reference string `json:"reference"`
	AccountID string `json:"account_id"`
	AccountNo int64  `json:"account_no"`
	Amount    int64  `json:"amount"`
}

type CreditCompletedEvent struct {
	EventMetadata
	UserID         string `json:"user_id"`
	ToAccID        string `json:"to_account_id"`
	Reference      string `json:"reference"`
	Amount         int64  `json:"amount"`
	ToBalanceAfter int64  `json:"to_balance_after"`
	Timestamp      int64  `json:"timestamp"`
}

type CreditFailedEvent struct {
	EventMetadata
	UserID    string `json:"user_id"`
	Reference string `json:"reference"`
	Reason    string `json:"reason"`
}

type CreditRequestedEvent struct {
	EventMetadata
	UserID    string `json:"user_id"`
	Reference string `json:"reference"`
	FromAccID string `json:"from_account_id"`
	FromAccNo int64  `json:"from_account_no"`
	ToAccID   string `json:"to_account_id"`
	ToAccNo   int64  `json:"to_account_no"`
	Amount    int64  `json:"amount"`
}

type DebitRequestedEvent struct {
	EventMetadata
	UserID    string `json:"user_id"`
	Reference string `json:"reference"`
	FromAccID string `json:"from_account_id"`
	FromAccNo int64  `json:"from_account_no"`
	ToAccID   string `json:"to_account_id"`
	ToAccNo   int64  `json:"to_account_no"`
	Amount    int64  `json:"amount"`
}

type TransferCompletedEvent struct {
	EventMetadata
	TransferID       string `json:"transfer_id"`
	Reference        string `json:"reference"`
	UserID           string `json:"user_id"`
	FromAccID        string `json:"from_account_id"`
	FromAccNo        int64  `json:"from_account_no"`
	ToAccID          string `json:"to_account_id,omitempty"`
	ToAccNo          int64  `json:"to_account_no,omitempty"`
	Amount           int64  `json:"amount"`
	Currency         string `json:"currency"`
	Description      string `json:"description"`
	FromBalanceAfter int64  `json:"from_balance_after"`
	ToBalanceAfter   int64  `json:"to_balance_after"`
	Timestamp        int64  `json:"timestamp"`
}

type TransferInitiatedEvent struct {
	EventMetadata
	UserID     string `json:"user_id"`
	TransferID string `json:"transfer_id"`
	Reference  string `json:"reference"`
	FromAccID  string `json:"from_account_id"`
	FromAccNo  int64  `json:"from_account_no"`
	ToAccID    string `json:"to_account_id"`
	ToAccNo    int64  `json:"to_account_no"`
	Amount     int64  `json:"amount"`
	Timestamp  int64  `json:"timestamp"`
}

type TransferFailedEvent struct {
	EventMetadata
	TransferID    string `json:"transfer_id"`
	Reference     string `json:"reference"`
	UserID        string `json:"user_id"`
	Amount        int64  `json:"amount"`
	FailureReason string `json:"failure_reason"`
	Timestamp     int64  `json:"timestamp"`
}

type LedgerCreatedEvent struct {
	EventMetadata
	LedgerID      string `json:"ledger_id"`
	Reference     string `json:"reference"`
	UserID        string `json:"user_id"`
	AccountID     string `json:"account_id"`
	Type          string `json:"type"` // credit, debit, transfer
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
	Description   string `json:"description"`
	Timestamp     int64  `json:"timestamp"`
}

type LedgerReconciliationAlertEvent struct {
	EventMetadata
	AccountID      string `json:"account_id"`
	AccountBalance int64  `json:"account_balance"`
	LedgerBalance  int64  `json:"ledger_balance"`
	Difference     int64  `json:"difference"`
	Timestamp      int64  `json:"timestamp"`
}

type AdminActionEvent struct {
	EventMetadata
	AdminID    string `json:"admin_id"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Reason     string `json:"reason"`
	Metadata   string `json:"metadata"`
	Timestamp  int64  `json:"timestamp"`
}

type DLQEvent struct {
	EventMetadata
	OriginalTopic string `json:"original_topic"`
	Error         string `json:"error"`
	RetryCount    int    `json:"retry_count"`
	Payload       string `json:"payload"`
	FailedAt      int64  `json:"failed_at"`
}
