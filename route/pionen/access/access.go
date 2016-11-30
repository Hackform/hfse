package access

const (
	// 0-127 level tiers
	// 128-255 level tags
	ROOT   uint8 = 0
	ADMIN  uint8 = 4
	MOD    uint8 = 64
	USER   uint8 = 127
	PUBLIC uint8 = 255
)
