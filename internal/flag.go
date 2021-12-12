package internal

const (
	serviceFlag    = "service"
	threadFlag     = "threads"
	durationFlag   = "duration"
	outputJsonPath = "output-json-path"

	redisHostFlag    = "host"
	redisPortFlag    = "port"
	redisClientsFlag = "clients"
)

var (
	service  string
	thread   int64
	duration int64

	redisHost    string
	redisPort    string
	redisClients int64
)
