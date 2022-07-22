package ship

import (
	"github.com/samber/lo"
)

const ProtocolHandshakeFormatJSON MessageProtocolFormatType = "JSON-UTF8"

// IsSupported validates if format is supported
func (m MessageProtocolFormatsType) IsSupported(format MessageProtocolFormatType) bool {
	return lo.Contains(m.Format, format)
}
