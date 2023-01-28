package api

const (
	// protocol config
	withProtobuf = "default"
	withThrift   = "optional"

	// distributed strategy config
	withConsistentHash = "default"
	withRaft           = "will be supported"

	// eliminate strategy config
	withLRU = "default"
	withLFU = "optional"

	// source config
	withRedisSource = "default"
	withMySQLSource = "will be supported"
)
