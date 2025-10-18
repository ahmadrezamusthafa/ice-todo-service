package repository

type StreamRepository interface {
	Publish(streamName string, data map[string]interface{}) (string, error)
	GetStream(streamName string, count int64, start string) ([]map[string]interface{}, error)
}
