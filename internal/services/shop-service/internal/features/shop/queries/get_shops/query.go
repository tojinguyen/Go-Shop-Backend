package getshops

// GetShopsQuery represents the query to get shops by owner
type GetShopsQuery struct {
	OwnerID string `json:"owner_id" validate:"required,uuid"`
}
