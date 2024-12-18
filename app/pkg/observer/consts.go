package observer

const (
	TagName      = "observer"
	TagIgnoreVal = "ignore"
)

var (
	blackListWords = []string{
		//grpc
		"state",
		"sizeCache",
		"unknownFields",
	}
)
