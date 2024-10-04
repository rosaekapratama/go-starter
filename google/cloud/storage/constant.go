package storage

const (
	spanUpload            = "common.google.cloud.storage.Upload"
	spanDownload          = "common.google.cloud.storage.Download"
	spanNewStreamUpload   = "common.google.cloud.storage.NewStreamUpload"
	spanNewStreamDownload = "common.google.cloud.storage.NewStreamDownload"

	PublicUrl        = "https://storage.googleapis.com"
	AuthenticatedUrl = "https://storage.googleapis.com"

	errWriterFormat = "errWriter=%s"
	errReaderFormat = "errReader=%s"
)
