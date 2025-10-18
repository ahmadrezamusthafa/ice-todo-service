package repository

type StreamRepository interface {
	Publish(streamName string, data map[string]interface{}) (string, error)
}
