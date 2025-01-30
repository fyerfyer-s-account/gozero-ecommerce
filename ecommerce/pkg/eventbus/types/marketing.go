package types

import "time"

type MarketingEventType string

const (
	// Coupon events
	CouponIssued    MarketingEventType = "marketing.coupon.issued"
	CouponUsed      MarketingEventType = "marketing.coupon.used"
	CouponExpired   MarketingEventType = "marketing.coupon.expired"
	CouponCancelled MarketingEventType = "marketing.coupon.cancelled"

	// Promotion events
	PromotionStarted MarketingEventType = "marketing.promotion.started"
	PromotionEnded   MarketingEventType = "marketing.promotion.ended"
	PromotionApplied MarketingEventType = "marketing.promotion.applied"

	// Points events
	PointsEarned   MarketingEventType = "marketing.points.earned"
	PointsUsed     MarketingEventType = "marketing.points.used"
	PointsExpired  MarketingEventType = "marketing.points.expired"
	PointsRefunded MarketingEventType = "marketing.points.refunded"

	MarketingPaymentSuccess MarketingEventType = "marketing.payment.success"
	MarketingPaymentFailed  MarketingEventType = "marketing.payment.failed"
)

// MarketingEvent represents the base marketing event structure
type MarketingEvent struct {
	Type      MarketingEventType `json:"type"`
	UserID    int64              `json:"user_id"`
	Timestamp time.Time          `json:"timestamp"`
}

// CouponEvent represents coupon-related events
type CouponEvent struct {
	MarketingEvent
	CouponID   int64   `json:"coupon_id"`
	CouponCode string  `json:"coupon_code"`
	Amount     float64 `json:"amount"`
	OrderNo    string  `json:"order_no,omitempty"`
}

// PromotionEvent represents promotion-related events
type PromotionEvent struct {
	MarketingEvent
	PromotionID   int64   `json:"promotion_id"`
	PromotionType int32   `json:"promotion_type"`
	Discount      float64 `json:"discount,omitempty"`
	OrderNo       string  `json:"order_no,omitempty"`
}

// PointsEvent represents points-related events
type PointsEvent struct {
	MarketingEvent
	Points  int64  `json:"points"`
	Balance int64  `json:"balance"`
	Source  string `json:"source"`
	OrderNo string `json:"order_no,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// MarketingPaymentSuccessEvent represents marketing actions after payment success
type MarketingPaymentSuccessEvent struct {
	MarketingEvent
	OrderNo      string  `json:"order_no"`
	PaymentNo    string  `json:"payment_no"`
	Amount       float64 `json:"amount"`
	RewardPoints int64   `json:"reward_points,omitempty"`
	CouponID     int64   `json:"coupon_id,omitempty"`
}

// MarketingPaymentFailedEvent represents marketing rollback after payment failure
type MarketingPaymentFailedEvent struct {
	MarketingEvent
	OrderNo   string `json:"order_no"`
	PaymentNo string `json:"payment_no"`
	Reason    string `json:"reason"`
}

// Add validation methods
func (e *CouponEvent) Validate() error {
	if e.CouponID == 0 || e.CouponCode == "" {
		return &NonRetryableError{
			EventError: &EventError{
				Code:    "INVALID_COUPON_EVENT",
				Message: "coupon_id and coupon_code are required",
			},
		}
	}
	return nil
}

func (e *PromotionEvent) Validate() error {
	if e.PromotionID == 0 {
		return &NonRetryableError{
			EventError: &EventError{
				Code:    "INVALID_PROMOTION_EVENT",
				Message: "promotion_id is required",
			},
		}
	}
	return nil
}

func (e *PointsEvent) Validate() error {
	if e.Points == 0 {
		return &NonRetryableError{
			EventError: &EventError{
				Code:    "INVALID_POINTS_EVENT",
				Message: "points value is required",
			},
		}
	}
	return nil
}

func (e *MarketingPaymentSuccessEvent) Validate() error {
	if e.OrderNo == "" || e.PaymentNo == "" {
		return &NonRetryableError{
			EventError: &EventError{
				Code:    "INVALID_MARKETING_PAYMENT_SUCCESS",
				Message: "order_no and payment_no are required",
			},
		}
	}
	return nil
}

func (e *MarketingPaymentFailedEvent) Validate() error {
	if e.OrderNo == "" || e.PaymentNo == "" {
		return &NonRetryableError{
			EventError: &EventError{
				Code:    "INVALID_MARKETING_PAYMENT_FAILED",
				Message: "order_no and payment_no are required",
			},
		}
	}
	return nil
}
