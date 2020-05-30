package structs

//Reference is the struct to get master refs
type Reference struct {
	ID     int    `json:"id"`
	Values string `json:"values"`
	Status int    `json:"status"`
}
