package types

import "github.com/google/uuid"

type UsageRecord struct {
	IdempotentKey  uuid.UUID `json:"idempotentKey"`
	OrganizationID uuid.UUID `json:"organizationID"`
	ProductID      string    `json:"productId"`
	Quantity       int64     `json:"quantity"`
	PublishTime    int64     `json:"publishTime"`
}