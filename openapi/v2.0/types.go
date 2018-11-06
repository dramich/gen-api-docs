package openapi

// OpenAPI -
type OpenAPI struct {
	Swagger             string                     `yaml:"swagger,omitempty" json:"swagger,omitempty"`
	Info                *Info                      `yaml:"info,omitempty" json:"info,omitempty"`
	Host                string                     `yaml:"host,omitempty" json:"host,omitempty"`
	BasePath            string                     `yaml:"basePath,omitempty" json:"basePath,omitempty"`
	Schemes             []string                   `yaml:"schemes,omitempty" json:"schemes,omitempty"`
	Consumes            []string                   `yaml:"consumes,omitempty" json:"consumes,omitempty"`
	Produces            []string                   `yaml:"produces,omitempty" json:"produces,omitempty"`
	Paths               map[string]*Path           `yaml:"paths,omitempty" json:"paths,omitempty"`
	Definitions         map[string]*Schema         `yaml:"definitions,omitempty" json:"definitions,omitempty"`
	Parameters          map[string]*Parameter      `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses           map[string]*Response       `yaml:"responses,omitempty" json:"responses,omitempty"`
	SecurityDefinitions map[string]*SecurityScheme `yaml:"securityDefinitions,omitempty" json:"securityDefinitions,omitempty"`
	Security            []map[string][]string      `yaml:"security,omitempty" json:"security,omitempty"`
	Tags                []*Tag                     `yaml:"tags,omitempty" json:"tags,omitempty"`
	ExternalDocs        *ExternalDocumentation     `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

// Tag -
type Tag struct {
	Name         string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Description  string                 `yaml:"description,omitempty" json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

// SecurityScheme -
type SecurityScheme struct {
	Type             string            `yaml:"type,omitempty" json:"type,omitempty"`
	Description      string            `yaml:"description,omitempty" json:"description,omitempty"`
	Name             string            `yaml:"name,omitempty" json:"name,omitempty"`
	In               string            `yaml:"in,omitempty" json:"in,omitempty"`
	Flow             string            `yaml:"flow,omitempty" json:"flow,omitempty"`
	AuthorizationURL string            `yaml:"authorizationUrl,omitempty" json:"authorizationUrl,omitempty"`
	TokenURL         string            `yaml:"tokenUrl,omitempty" json:"tokenUrl,omitempty"`
	Scopes           map[string]string `yaml:"scopes,omitempty" json:"scopes,omitempty"`
}

// Path -
type Path struct {
	Ref        string       `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	GET        *Operation   `yaml:"get,omitempty" json:"get,omitempty"`
	PUT        *Operation   `yaml:"put,omitempty" json:"put,omitempty"`
	POST       *Operation   `yaml:"post,omitempty" json:"post,omitempty"`
	DELETE     *Operation   `yaml:"delete,omitempty" json:"delete,omitempty"`
	OPTIONS    *Operation   `yaml:"options,omitempty" json:"options,omitempty"`
	HEAD       *Operation   `yaml:"head,omitempty" json:"head,omitempty"`
	PATCH      *Operation   `yaml:"patch,omitempty" json:"patch,omitempty"`
	Parameters []*Parameter `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

// Operation -
type Operation struct {
	Tags         []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
	Summary      string                 `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description  string                 `yaml:"description,omitempty" json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	OperationID  string                 `yaml:"operationID,omitempty" json:"operationID,omitempty"`
	Consumes     []string               `yaml:"consumes,omitempty" json:"consumes,omitempty"`
	Produces     []string               `yaml:"produces,omitempty" json:"produces,omitempty"`
	Parameters   map[string]*Parameter  `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses    map[string]*Response   `yaml:"responses,omitempty" json:"responses,omitempty"`
	Schemes      []string               `yaml:"schemes,omitempty" json:"schemes,omitempty"`
	Deprecated   bool                   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Security     map[string][]string    `yaml:"security,omitempty" json:"security,omitempty"`
}

// Response -
type Response struct {
	Ref         string                  `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Description string                  `yaml:"description,omitempty" json:"description,omitempty"`
	Schema      *Schema                 `yaml:"schema,omitempty" json:"schema,omitempty"`
	Headers     map[string]*Header      `yaml:"headers,omitempty" json:"headers,omitempty"`
	Examples    map[string]*interface{} `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// Header -
type Header struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Type        string `yaml:"type,omitempty" json:"type,omitempty"`
}

// Parameter -
type Parameter struct {
	Ref              string         `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Name             string         `yaml:"name,omitempty" json:"name,omitempty"`
	In               string         `yaml:"in,omitempty" json:"in,omitempty"`
	Description      string         `yaml:"description,omitempty" json:"description,omitempty"`
	Required         bool           `yaml:"required,omitempty" json:"required,omitempty"`
	Schema           *Schema        `yaml:"schema,omitempty" json:"schema,omitempty"`
	Type             string         `yaml:"type,omitempty" json:"type,omitempty"`
	Format           string         `yaml:"format,omitempty" json:"format,omitempty"`
	AllowEmptyValue  bool           `yaml:"allowEmptyValue,omitempty" json:"allowEmptyValue,omitempty"`
	Items            *Items         `yaml:"items,omitempty" json:"items,omitempty"`
	CollectionFormat string         `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Default          *interface{}   `yaml:"default,omitempty" json:"default,omitempty"`
	Maximum          int            `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	ExclusiveMaximum bool           `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	Minimum          int            `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	ExclusiveMinimum bool           `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	MaxLength        int            `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	MinLength        int            `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	Pattern          string         `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MaxItems         int            `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	MinItems         int            `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	UniqueItems      bool           `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`
	Enum             []*interface{} `yaml:"enum,omitempty" json:"enum,omitempty"`
	MultipleOf       int            `yaml:"multipleOf,omitempty" json:"multipleOf,omitempty"`
}

// Items -
type Items struct {
	Type             string         `yaml:"type,omitempty" json:"type,omitempty"`
	Format           string         `yaml:"format,omitempty" json:"format,omitempty"`
	Items            *Items         `yaml:"items,omitempty" json:"items,omitempty"`
	CollectionFormat string         `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Default          *interface{}   `yaml:"default,omitempty" json:"default,omitempty"`
	Maximum          int            `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	ExclusiveMaximum bool           `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	Minimum          int            `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	ExclusiveMinimum bool           `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	MaxLength        int            `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	MinLength        int            `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	Pattern          string         `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MaxItems         int            `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	MinItems         int            `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	UniqueItems      bool           `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`
	Enum             []*interface{} `yaml:"enum,omitempty" json:"enum,omitempty"`
	MultipleOf       int            `yaml:"multipleOf,omitempty" json:"multipleOf,omitempty"`
}

// Schema -
type Schema struct {
	Ref                  string                 `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Format               string                 `yaml:"format,omitempty" json:"format,omitempty"`
	Title                string                 `yaml:"title,omitempty" json:"title,omitempty"`
	Description          string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Default              *interface{}           `yaml:"default,omitempty" json:"default,omitempty"`
	MultipleOf           int                    `yaml:"multipleOf,omitempty" json:"multipleOf,omitempty"`
	Maximum              int                    `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	Minimum              int                    `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	MaxLength            int                    `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	MinLength            int                    `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	Pattern              string                 `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MaxItems             int                    `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	MinItems             int                    `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	UniqueItems          bool                   `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`
	MaxProperties        int                    `yaml:"maxProperties,omitempty" json:"maxProperties,omitempty"`
	MinProperties        int                    `yaml:"minProperties,omitempty" json:"minProperties,omitempty"`
	Required             []string               `yaml:"required,omitempty" json:"required,omitempty"`
	Enum                 []string               `yaml:"enum,omitempty" json:"enum,omitempty"`
	Type                 string                 `yaml:"type,omitempty" json:"type,omitempty"`
	Items                []*Schema              `yaml:"items,omitempty" json:"items,omitempty"`
	AllOf                []*Schema              `yaml:"allOf,omitempty" json:"allOf,omitempty"`
	Properties           map[string]*Schema     `yaml:"properties,omitempty" json:"properties,omitempty"`
	AdditionalProperties *Schema                `yaml:"additionalProperties,omitempty" json:"additionalProperties,omitempty"`
	Discriminator        string                 `yaml:"discriminator,omitempty" json:"discriminator,omitempty"`
	ReadOnly             bool                   `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`
	XML                  string                 `yaml:"xml,omitempty" json:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Example              *interface{}           `yaml:"example,omitempty" json:"example,omitempty"`
}

// ExternalDocumentation -
type ExternalDocumentation struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	URL         string `yaml:"url,omitempty" json:"url,omitempty"`
}

// Info -
type Info struct {
	Title          string   `yaml:"title,omitempty" json:"title,omitempty"`
	Description    string   `yaml:"description,omitempty" json:"description,omitempty"`
	TermsOfService string   `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	Contact        *Contact `yaml:"contact,omitempty" json:"contact,omitempty"`
	License        *License `yaml:"license,omitempty" json:"license,omitempty"`
	Version        string   `yaml:"version,omitempty" json:"version,omitempty"`
}

// Contact -
type Contact struct {
	Name  string `yaml:"name,omitempty" json:"name,omitempty"`
	URL   string `yaml:"url,omitempty" json:"url,omitempty"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

// License -
type License struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	URL  string `yaml:"URL,omitempty" json:"URL,omitempty"`
}
