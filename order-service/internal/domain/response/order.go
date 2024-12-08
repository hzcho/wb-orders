package response

type GetOrderIDs struct {
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	IDs   []string `json:"ids"`
}
