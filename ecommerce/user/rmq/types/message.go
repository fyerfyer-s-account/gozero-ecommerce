package types

import "time"

type EventType string

const (
	// User account events
	EventTypeUserRegistered      EventType = "user.registered"
	EventTypeUserProfileUpdated  EventType = "user.profile.updated"
	EventTypeUserPasswordChanged EventType = "user.password.changed"
	EventTypeUserPasswordReset   EventType = "user.password.reset"

	// Address events
	EventTypeAddressCreated EventType = "user.address.created"
	EventTypeAddressUpdated EventType = "user.address.updated"
	EventTypeAddressDeleted EventType = "user.address.deleted"

	// Wallet events
	EventTypeWalletRecharged    EventType = "user.wallet.recharged"
	EventTypeTransactionCreated EventType = "user.transaction.created"
)

type UserEvent struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	UserID    int64     `json:"userId"`
	Timestamp int64     `json:"timestamp"`
	Data      any       `json:"data"`
	Metadata  Metadata  `json:"metadata"`
}

type Metadata struct {
	TraceID string `json:"traceId,omitempty"`
	Source  string `json:"source"`
	Version string `json:"version"`
}

type ProfileData struct {
	Username    string  `json:"username,omitempty"`
	Nickname    string  `json:"nickname,omitempty"`
	Avatar      string  `json:"avatar,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Email       string  `json:"email,omitempty"`
	Gender      string  `json:"gender,omitempty"`
	MemberLevel int32   `json:"memberLevel,omitempty"`
	Balance     float64 `json:"balance,omitempty"`
}

type AddressData struct {
	ID            int64  `json:"id"`
	ReceiverName  string `json:"receiverName"`
	ReceiverPhone string `json:"receiverPhone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detailAddress"`
	IsDefault     bool   `json:"isDefault"`
}

type WalletData struct {
	Amount        float64 `json:"amount"`
	Balance       float64 `json:"balance"`
	FrozenAmount  float64 `json:"frozenAmount,omitempty"`
	TransactionID string  `json:"transactionId,omitempty"`
	Type          int32   `json:"type,omitempty"`
}

func NewUserEvent(eventType EventType, userID int64, data any) *UserEvent {
	return &UserEvent{
		ID:        GenerateEventID(),
		Type:      eventType,
		UserID:    userID,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
		Metadata: Metadata{
			Source:  "user.service",
			Version: "1.0",
		},
	}
}

func GenerateEventID() string {
	return time.Now().Format("20060102150405.000") + "-" + RandomString(8)
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
