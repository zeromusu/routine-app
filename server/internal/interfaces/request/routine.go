package request

type CreateRoutineRequest struct {
	Title    string `json:"title" binding:"required"`
	Interval string `json:"interval" binding:"required"`
}

type UpdateRoutineRequest struct {
	Title    string `json:"title,omitempty"`
	Interval string `json:"interval,omitempty"`
}
