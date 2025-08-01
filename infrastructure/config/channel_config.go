package config

type ChannelConfig struct {
	BalanceChannelBufferSize int `env:"BALANCE_CHANNEL_BUFFER_SIZE" envDefault:"1000"`
}
