package bigmapdiff

// GetContext -
type GetContext struct {
	Network      string
	Ptr          *int64
	Query        string
	Size         int64
	Offset       int64
	Level        *int64
	CurrentLevel *int64
	Contract     string

	To int64
}