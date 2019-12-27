package sandpiper

// Company represents company model
type Company struct {
	Base
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	Owner     User       `json:"owner"`
}
