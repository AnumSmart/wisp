package match

import "time"

type Match struct {
	ID                   string    `json:"id"`
	InitiatorID          string    `json:"initiator_id"`  // Кто инициировал
	TargetID             string    `json:"target_id"`     // Кого нашли
	CompatibilityPercent int       `json:"compatibility"` // 0-100%
	Status               string    `json:"status"`        // "pending"/"accepted"/"rejected"
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type MatchAction struct {
	ID        string    `json:"id"`
	MatchID   string    `json:"match_id"`
	Action    string    `json:"action"` // "like", "skip", "report"
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
