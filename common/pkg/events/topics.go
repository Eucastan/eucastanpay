package events

const (
	// User Events
	TopicUserRegistered  = "user.registered"
	TopicUserKYCCreated  = "user.kyc.created"
	TopicUserKYCVerified = "user.kyc.verified"

	// Account Events
	TopicAccountCreated   = "account.created"
	TopicCreateAccRequest = "create.account.request"
	TopicDepositAccount   = "deposit.account"
	TopicWithdrawal       = "withdraw.debit"

	// Transfer Events
	TopicTransferInitiated = "transfer.initiated"
	TopicReverseInitiated  = "reverse.initiated"
	TopicDebitRequested    = "debit.requested"
	TopicDebitCompleted    = "debit.completed"
	TopicCreditRequested   = "credit.requested"
	TopicCreditCompleted   = "credit.completed"
	TopicTransferCompleted = "transfer.completed"
	TopicTransferRetry     = "transfer.retry"
	TopicLedgerCreated     = "ledger.created"
	TopicAuditCreated      = "audit.created"

	// Admin events
	TopicAdminActionTaken     = "admin.action.taken"
	TopicUserStatusChanged    = "user.status.changed"
	TopicAccountStatusChanged = "account.status.changed"
	TopicTransferReversed     = "transfer.reversed"

	// Failed Events
	TopicUserRegistrationFailed = "user.register.failed"
	TopicCreateAccFailed        = "create.account.failed"
	TopicTransferFailed         = "transfer.failed"
	TopicCreditFailed           = "credit.failed"
	TopicDebitFailed            = "debit.failed"
	TopicAuditFailed            = "audit.failed"

	// Ledger Events
	TopicLedgerReconciliationAlert = "ledger.reconciliation.alert"

	// DLQ
	TopicTransferDLQ     = "transfer.dlq"
	TopicAccountDLQ      = "account.dlq"
	TopicAdminDLQ        = "admin.dlq"
	TopicLedgerDLQ       = "ledger.dlq"
	TopicAuditDLQ        = "audit.dlq"
	TopicNotificationDLQ = "notification.dlq"
)
