package openapi

// OpenAPI - Top level OpenAPI Object
// https://swagger.io/specification/#openapi-object
type OpenAPI struct {
	OpenAPI      string                `yaml:"openapi"`
	Info         Info                  `yaml:"info"`
	Servers      []Server              `yaml:"servers"`
	Paths        map[string]PathItem   `yaml:"paths"`
	Components   Components            `yaml:"components,omitempty"`
	Security     []map[string][]string `yaml:"security,omitempty"`
	Tags         []Tag                 `yaml:"tags,omitempty"`
	ExternalDocs ExternalDocumentation `yaml:"externalDocs,omitempty"`
}

// Info - https://swagger.io/specification/#infoObject
type Info struct {
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version,omitempty"`
}

// Server - https://swagger.io/specification/#serverObject
type Server struct {
	URL         string                    `yaml:"url,omitempty"`
	Description string                    `yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `yaml:"variables,omitempty"`
}

// PathItem - https://swagger.io/specification/#pathsItemObject
type PathItem struct {
	Ref         string      `yaml:"$ref,omitempty"`
	Summary     string      `yaml:"summary,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Get         *Operation  `yaml:"get,omitempty"`
	Post        *Operation  `yaml:"post,omitempty"`
	Put         *Operation  `yaml:"put,omitempty"`
	Delete      *Operation  `yaml:"delete,omitempty"`
	Options     *Operation  `yaml:"options,omitempty"`
	Head        *Operation  `yaml:"head,omitempty"`
	Patch       *Operation  `yaml:"patch,omitempty"`
	Trace       *Operation  `yaml:"trace,omitempty"`
	Servers     []Server    `yaml:"servers,omitempty"`
	Parameters  []Parameter `yaml:"parameters,omitempty"` // or Ref
}

// Components - https://swagger.io/specification/#componentsObject
type Components struct {
	Schemas         map[string]Schema         `yaml:"schemas,omitempty"`         // or Ref
	Responses       map[string]Response       `yaml:"responses,omitempty"`       // or Ref
	Parameters      map[string]Parameter      `yaml:"parameters,omitempty"`      // or Ref
	Examples        map[string]Example        `yaml:"examples,omitempty"`        // or Ref
	RequestBodies   map[string]RequestBody    `yaml:"requestBodies,omitempty"`   // or Ref
	Headers         map[string]Parameter      `yaml:"headers,omitempty"`         // or Ref
	SecuritySchemes map[string]SecurityScheme `yaml:"securitySchemes,omitempty"` // or Ref
	Links           map[string]Link           `yaml:"links,omitempty"`           // or Ref
	Callbacks       map[string]PathItem       `yaml:"callbacks,omitempty"`       // or Ref
}

// SecurityScheme - https://swagger.io/specification/#securitySchemeObject
type SecurityScheme struct {
	Ref              string      `yaml:"$ref,omitempty"`
	Type             string      `yaml:"type,omitempty"`
	Description      string      `yaml:"description,omitempty"`
	Name             string      `yaml:"name,omitempty"`
	In               string      `yaml:"in,omitempty"`
	Scheme           string      `yaml:"scheme,omitempty"`
	BearerFormat     string      `yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `yaml:"flows,omitempty"`
	OpenIDConnectURL string      `yaml:"openIdConnectUrl,omitempty"`
}

// OAuthFlows - https://swagger.io/specification/#oauthFlowsObject
type OAuthFlows struct {
	Implicit          *OAuthFlow `yaml:"implicit,omitempty"`
	Password          *OAuthFlow `yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `yaml:"authorizationCode,omitempty"`
}

// OAuthFlow - https://swagger.io/specification/#oauthFlowObject
type OAuthFlow struct {
	authorizationURL string            `yaml:"authorizationUrl,omitempty"`
	tokenURL         string            `yaml:"tokenUrl,omitempty"`
	refreshURL       string            `yaml:"refreshUrl,omitempty"`
	scopes           map[string]string `yaml:"scopes,omitempty"`
}

// Parameter - https://swagger.io/specification/#parameterObject
// Headers - They follow parameter objects
// if this is an "or $ref" we may need to mark all the fields omitempty
type Parameter struct {
	Ref             string               `yaml:"$ref,omitempty"`
	Name            string               `yaml:"name,omitempty"`
	In              string               `yaml:"in,omitempty"`
	Description     string               `yaml:"description,omitempty"`
	Required        bool                 `yaml:"required,omitempty"`
	Deprecated      bool                 `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                 `yaml:"allowEmptyValue,omitempty"`
	Style           string               `yaml:"style,omitempty"`
	Explode         bool                 `yaml:"explode,omitempty"`
	AllowReserved   bool                 `yaml:"allowReserved,omitempty"`
	Schema          *Schema              `yaml:"schema,omitempty"` // or ref
	Example         *interface{}         `yaml:"example,omitempty"`
	Examples        map[string]Example   `yaml:"examples,omitempty"` // or ref
	Content         map[string]MediaType `yaml:"content,omitempty"`
}

// RequestBody - https://swagger.io/specification/#requestBodyObject
type RequestBody struct {
	Ref         string               `yaml:"$ref,omitempty"`
	Description string               `yaml:"description,omitempty"`
	Content     map[string]MediaType `yaml:"content,omitempty"`
	Required    bool                 `yaml:"required,omitempty"`
}

// MediaType - https://swagger.io/specification/#mediaTypeObject
type MediaType struct {
	Schema   *Schema             `yaml:"schema,omitempty"` // or ref
	Example  *interface{}        `yaml:"example,omitempty"`
	Examples map[string]Example  `yaml:"examples,omitempty"` // or ref
	Encoding map[string]Encoding `yaml:"encoding,omitempty"`
}

// Encoding - https://swagger.io/specification/#encodingObject
type Encoding struct {
	ContentType   string               `yaml:"contentType,omitempty"`
	Headers       map[string]Parameter `yaml:"headers,omitempty"` // or ref
	Style         string               `yaml:"style,omitempty"`
	Explode       bool                 `yaml:"explode,omitempty"`
	AllowReserved bool                 `yaml:"allowReserved,omitempty"`
}

// Example - https://swagger.io/specification/#exampleObject
type Example struct {
	Ref           string       `yaml:"$ref,omitempty"`
	Summary       string       `yaml:"summary,omitempty"`
	Description   string       `yaml:"description,omitempty"`
	Value         *interface{} `yaml:"value,omitempty"`
	ExternalValue string       `yaml:"externalValue,omitempty"`
}

// Operation - https://swagger.io/specification/#operationObject
type Operation struct {
	Tags         []string               `yaml:"tags,omitempty"`
	Summary      string                 `yaml:"summary,omitempty"`
	Description  string                 `yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty"`
	OperationID  string                 `yaml:"operationId,omitempty"`
	Parameters   []Parameter            `yaml:"parameters,omitempty"`
	RequestBody  *RequestBody           `yaml:"requestBody,omitempty"`
	Responses    map[string]Response    `yaml:"responses,omitempty"`
	Callbacks    map[string]PathItem    `yaml:"callbacks,omitempty"`
	Deprecated   bool                   `yaml:"deprecated,omitempty"`
	Security     map[string][]string    `yaml:"security,omitempty"`
	Servers      []Server               `yaml:"servers,omitempty"`
}

// Schema - https://swagger.io/specification/#schemaObject
type Schema struct {
	Ref string `yaml:"$ref,omitempty"`
	// JSON Schema Fields
	Title            string   `yaml:"title,omitempty"`
	MultipleOf       *int64   `yaml:"multipleOf,omitempty"`
	Maximum          *int64   `yaml:"maximum,omitempty"`
	ExclusiveMaximum *int64   `yaml:"exclusiveMaximum,omitempty"`
	Minimum          *int64   `yaml:"minimum,omitempty"`
	ExclusiveMinimum *int64   `yaml:"exclusiveMinimum,omitempty"`
	MaxLength        *int64   `yaml:"maxLength,omitempty"`
	MinLength        *int64   `yaml:"minLength,omitempty"`
	Pattern          string   `yaml:"pattern,omitempty"`
	MaxItems         *int64   `yaml:"maxItems,omitempty"`
	MinItems         *int64   `yaml:"minItems,omitempty"`
	UniqueItems      bool     `yaml:"uniqueItems,omitempty"`
	MaxProperties    *int64   `yaml:"maxProperties,omitempty"`
	MinProperties    *int64   `yaml:"minProperties,omitempty"`
	Required         []string `yaml:"required,omitempty"`
	Enum             []string `yaml:"enum,omitempty"`

	// OpenAPI modified JSON Schema Fields
	Type                 string            `yaml:"type,omitempty"`
	Description          string            `yaml:"description,omitempty"`
	Properties           map[string]Schema `yaml:"properties,omitempty"`
	AdditionalProperties *Schema           `yaml:"additionalProperties,omitempty"`
	Format               string            `yaml:"format,omitempty"`
	Default              interface{}       `yaml:"default,omitempty"`
	Items                *Schema           `yaml:"items,omitempty"`
	AllOf                []Schema          `yaml:"allOf,omitempty"`
	OneOf                []Schema          `yaml:"oneOf,omitempty"`
	AnyOf                []Schema          `yaml:"anyOf,omitempty"`
	Not                  []Schema          `yaml:"not,omitempty"`

	// Static OpenAPI Fields
	Discriminator *Discriminator         `yaml:"discriminator,omitempty"`
	Nullable      bool                   `yaml:"nullable,omitempty"`
	ReadOnly      bool                   `yaml:"readOnly,omitempty"`
	WriteOnly     bool                   `yaml:"writeOnly,omitempty"`
	XML           *XML                   `yaml:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `yaml:"externalDocs,omitempty"`
	Example       interface{}            `yaml:"example,omitempty"`
	Deprecated    bool                   `yaml:"deprecated,omitempty"`
}

// Discriminator - https://swagger.io/specification/#discriminatorObject
type Discriminator struct {
	PropertyName string            `yaml:"propertyName,omitempty"`
	Mapping      map[string]string `yaml:"mapping,omitempty"`
}

// XML - https://swagger.io/specification/#xmlObject
type XML struct {
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
	Prefix    string `yaml:"prefix,omitempty"`
	Attribute bool   `yaml:"attribute,omitempty"`
	Warpped   bool   `yaml:"wrapped,omitempty"`
}

// Response - https://swagger.io/specification/#responseObject
// if this is an "or $ref" we may need to mark all the fields omitempty
type Response struct {
	Ref         string               `yaml:"$ref,omitempty"`
	Description string               `yaml:"description,omitempty"`
	Headers     map[string]Parameter `yaml:"headers,omitempty"`
	Content     map[string]MediaType `yaml:"content,omitempty"`
	Links       map[string]Link      `yaml:"links,omitempty"`
}

// Link - https://swagger.io/specification/#linkObject
// if this is an "or $ref" we may need to mark all the fields omitempty
// expression: https://swagger.io/specification/#runtimeExpression
type Link struct {
	Ref          string            `yaml:"$ref,omitempty"`
	OperationRef string            `yaml:"operationRef,omitempty"`
	OperationID  string            `yaml:"operationID,omitempty"`
	Parameters   map[string]string `yaml:"parameters,omitempty"`  //expression
	RequestBody  string            `yaml:"requestBody,omitempty"` // expression
	Description  string            `yaml:"description,omitempty"`
	Server       *Server           `yaml:"server,omitempty"`
}

// Tag - https://swagger.io/specification/#tagObject
type Tag struct {
	Name         string                 `yaml:"name"`
	Description  string                 `yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty"`
}

// ExternalDocumentation - https://swagger.io/specification/#externalDocumentationObject
type ExternalDocumentation struct {
	Description string `yaml:"description,omitempty"`
	URL         string `yaml:"url"`
}

// ServerVariable - https://swagger.io/specification/#serverVariableObject
type ServerVariable struct {
	Enum        []string `yaml:"enum,omitempty"`
	Default     string   `yaml:"default"`
	Description string   `yaml:"description"`
}
