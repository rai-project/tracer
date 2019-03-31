//go:generate go get github.com/a-h/generate/...
//go:generate schema-generate -schemaKeyRequired=false -p schema -o schema_generated.go file-format-schema.json
//go:generate sed -i -e 's/\$schemaReceived/schemaReceived/g' schema_generated.go
package schema
