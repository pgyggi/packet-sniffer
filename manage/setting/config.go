package setting

//
type Conf struct {
	Device string
	Kafka  Kafka
}

//
type Kafka struct {
	Hosts string
	Topic string
}
