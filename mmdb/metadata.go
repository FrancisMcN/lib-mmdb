package mmdb

import (
	"github.com/FrancisMcN/lib-mmdb2/field"
	"time"
)

type Metadata struct {
	NodeCount                uint32
	RecordSize               uint16
	IpVersion                uint16
	DatabaseType             string
	Languages                []string
	BinaryFormatMajorVersion uint16
	BinaryFormatMinorVersion uint16
	BuildEpoch               time.Time
	Description              map[string]string
}

func (m Metadata) Bytes() []byte {

	metadata := field.Map(make(map[field.Field]field.Field))
	metadata[field.String("node_count")] = field.Uint32(m.NodeCount)
	metadata[field.String("record_size")] = field.Uint16(m.RecordSize)
	metadata[field.String("ip_version")] = field.Uint16(m.IpVersion)
	metadata[field.String("database_type")] = field.String(m.DatabaseType)
	langs := make([]field.Field, 0)
	for _, l := range m.Languages {
		langs = append(langs, field.String(l))
	}
	metadata[field.String("languages")] = field.Array(langs)
	metadata[field.String("binary_format_major_version")] = field.Uint16(m.BinaryFormatMajorVersion)
	metadata[field.String("binary_format_minor_version")] = field.Uint16(m.BinaryFormatMinorVersion)
	metadata[field.String("build_epoch")] = field.Uint64(m.BuildEpoch.Unix())
	descriptions := make(map[field.Field]field.Field)
	for k, v := range m.Description {
		descriptions[field.String(k)] = field.String(v)
	}
	metadata[field.String("description")] = field.Map(descriptions)

	return metadata.Bytes()

}

func ParseMetadata(b []byte) Metadata {

	fp := field.FieldParserSingleton()
	fieldMap := fp.Parse(b).(field.Map)

	var nodeCount uint32
	var recordSize uint16
	var ipVersion uint16
	var databaseType string
	var languages []string
	var binaryFormatMajorVersion uint16
	var binaryFormatMinorVersion uint16
	var buildEpoch time.Time
	var description map[string]string

	if v, f := fieldMap[field.String("node_count")]; f && v.Type() == field.Uint32Field {
		nodeCount = uint32(v.(field.Uint32))
	}
	if v, f := fieldMap[field.String("record_size")]; f && v.Type() == field.Uint16Field {
		recordSize = uint16(v.(field.Uint16))
	}
	if v, f := fieldMap[field.String("ip_version")]; f && v.Type() == field.Uint16Field {
		ipVersion = uint16(v.(field.Uint16))
	}
	if v, f := fieldMap[field.String("database_type")]; f && v.Type() == field.StringField {
		databaseType = string(v.(field.String))
	}
	if v, f := fieldMap[field.String("languages")]; f && v.Type() == field.ArrayField {
		languages = make([]string, 0)
		for _, v := range v.(field.Array) {
			languages = append(languages, v.String())
		}
	}
	if v, f := fieldMap[field.String("binary_format_major_version")]; f && v.Type() == field.Uint16Field {
		binaryFormatMajorVersion = uint16(v.(field.Uint16))
	}
	if v, f := fieldMap[field.String("binary_format_minor_version")]; f && v.Type() == field.Uint16Field {
		binaryFormatMinorVersion = uint16(v.(field.Uint16))
	}
	if v, f := fieldMap[field.String("build_epoch")]; f && v.Type() == field.Uint64Field {
		buildEpoch = time.Unix(int64(v.(field.Uint64)), 0)
	}
	if v, f := fieldMap[field.String("description")]; f && v.Type() == field.MapField {
		description = make(map[string]string)
		for k, v := range v.(field.Map) {
			description[k.String()] = v.String()
		}
	}
	m := Metadata{
		NodeCount:                nodeCount,
		RecordSize:               recordSize,
		IpVersion:                ipVersion,
		DatabaseType:             databaseType,
		Languages:                languages,
		BinaryFormatMajorVersion: binaryFormatMajorVersion,
		BinaryFormatMinorVersion: binaryFormatMinorVersion,
		BuildEpoch:               buildEpoch,
		Description:              description,
	}
	return m
}
