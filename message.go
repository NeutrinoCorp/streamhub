package streamhub

import (
	"strconv"
	"strings"
	"time"
)

// CloudEventsSpecVersion the CloudEvents specification version used by streamhub
const CloudEventsSpecVersion = "1.0"

// Message is a unit of information which holds the primitive message (data) in binary format along multiple
// fields in order to preserve a schema definition within a stream pipeline.
//
// The schema is based on the Cloud Native Computing Foundation (CNCF)'s CloudEvents specification.
//
// For more information, please look: https://github.com/cloudevents/spec
type Message struct {
	ID          string `json:"id"`
	Stream      string `json:"stream"`
	Source      string `json:"source"`
	SpecVersion string `json:"specversion"`
	Type        string `json:"type"`
	Data        []byte `json:"data"`

	// Optional fields

	DataContentType   string `json:"datacontenttype,omitempty"`
	DataSchema        string `json:"dataschema,omitempty"`
	DataSchemaVersion int    `json:"dataschemaversion,omitempty"`
	Timestamp         string `json:"time,omitempty"`
	Subject           string `json:"subject,omitempty"`

	// Streamhub fields
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"`

	// consumer-only fields
	DecodedData interface{} `json:"-"`
	GroupName   string      `json:"-"`
}

// NewMessageArgs arguments required by NewMessage function to operate.
type NewMessageArgs struct {
	SchemaVersion        int
	Data                 []byte
	ID                   string
	Source               string
	Stream               string
	SchemaDefinitionName string
	ContentType          string
	GroupName            string
	Subject              string
}

// NewMessage allocates an immutable Message ready to be transported in a stream.
func NewMessage(args NewMessageArgs) Message {
	strSchemaVersion := strconv.Itoa(args.SchemaVersion)
	return Message{
		ID:                args.ID,
		Stream:            args.Stream,
		Source:            args.Source,
		SpecVersion:       CloudEventsSpecVersion,
		Type:              generateMessageType(args.Source, args.Stream, strSchemaVersion),
		Data:              args.Data,
		DataContentType:   args.ContentType,
		DataSchema:        args.SchemaDefinitionName,
		DataSchemaVersion: args.SchemaVersion,
		Timestamp:         time.Now().UTC().Format(time.RFC3339),
		Subject:           args.Subject,
		GroupName:         args.GroupName,
	}
}

func generateMessageType(source, stream, version string) string {
	var buff strings.Builder
	if source != "" && !strings.HasPrefix(stream, source) {
		buff.WriteString(source)
		buff.WriteString(".")
	}
	buff.WriteString(stream)
	if version != "" {
		buff.WriteString(".v")
		buff.WriteString(version)
	}
	return buff.String()
}
