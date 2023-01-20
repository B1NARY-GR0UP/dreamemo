package api

// Dreamemo Have we ever been sober
type Dreamemo struct {
	// TODO: 确定各个层可提供的选项，例如分布式层提供一致性哈希和 raft 两种选项，默认一致性哈希
	options *Options
}
