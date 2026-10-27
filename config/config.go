package config

type Config struct {
	AppPort          string
	AppEnv           string
	DBHost           string
	DBPort           string
	DBName           string
	DBUser           string
	DBPassword       string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	JWTSecret        string
	JWTExpiryHours   int
	MidtransServerKey string
	MidtransClientKey string
	MidtransEnv      string
}
