package api

// AddObservation represents an observation associated with an entity.
type AddObservation struct {
	EntityName string   `json:"entityName" validate:"required"`
	Contents   []string `json:"contents" validate:"required,min=1"`
}
type AddObservationsRequest struct {
	Observations []AddObservation `json:"observations"  validate:"required,min=1"`
}

type AddObservationsResponse struct {
}

type DeleteObservationsRequest struct {
}

type DeleteObservationsResponse struct {
	Success bool `json:"success"`
}
