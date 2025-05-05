package validation

type Config struct {
	MaxAmount int
}

var ValidationConfig = Config{
	MaxAmount: 1000, // The maximum amount allowed for a payment
}
