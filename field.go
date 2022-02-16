package paack

type Field struct {
	Name  string      `json:"name"`
	Type  FieldType   `json:"type"`
	Value interface{} `json:"value"`
}

type FieldType string

const (
	FieldTypeString   FieldType = "string"
	FieldTypeCurrency FieldType = "currency"
	FieldTypeCountry  FieldType = "country"
	FieldTypeNumber   FieldType = "number"
	FieldTypeBoolean  FieldType = "boolean"
	FieldTypeUuid     FieldType = "uuid"
	FieldTypeUri      FieldType = "uri"
	FieldTypeDate     FieldType = "date"
	FieldTypeTime     FieldType = "time"
	FieldTypeObject   FieldType = "object"
	FieldTypeArray    FieldType = "array"
	FieldTypeNil      FieldType = "null"
)
