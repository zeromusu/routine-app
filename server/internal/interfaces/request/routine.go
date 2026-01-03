package request

type CreateRoutineRequest struct {
	Title    string `json:"title" binding:"required"`
	Interval string `json:"interval" binding:"required"`
}
