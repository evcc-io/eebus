package ship

import "github.com/thoas/go-funk"

const ProtocolHandshakeFormatJSON MessageProtocolFormatType = "JSON-UTF8"

// IsSupported validates if format is supported
func (m MessageProtocolFormatsType) IsSupported(format MessageProtocolFormatType) bool {
	return funk.Contains(m.Format, format)
}
