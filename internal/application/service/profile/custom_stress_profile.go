package profile

type CustomStressProfile struct {
}

type CustomProfileRequest struct {
	Method   string
	Url      string
	BodyType string
	Body     string
	Headers  map[string]string
}

type CustomProfileConfig struct {
}
