package openapi

// OpenAPI - Top level OpenAPI Object
// https://swagger.io/specification/#openapi-object
type OpenAPI struct {
	OpenAPI      string                `yaml:"openapi" json:"openapi"`
	Info         Info                  `yaml:"info" json:"info"`
	Servers      []Server              `yaml:"servers" json:"servers"`
	Paths        map[string]PathItem   `yaml:"paths" json:"paths"`
	Components   Components            `yaml:"components,omitempty" json:"components,omitempty"`
	Security     []map[string][]string `yaml:"security,omitempty" json:"security,omitempty"`
	Tags         []Tag                 `yaml:"tags,omitempty" json:"tags,omitempty"`
	ExternalDocs ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

// Info - https://swagger.io/specification/#infoObject
type Info struct {
	Title       string `yaml:"title,omitempty" json:"title,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Version     string `yaml:"version,omitempty" json:"version,omitempty"`
}

// Server - https://swagger.io/specification/#serverObject
type Server struct {
	URL         string                    `yaml:"url,omitempty" json:"url,omitempty"`
	Description string                    `yaml:"description,omitempty" json:"description,omitempty"`
	Variables   map[string]ServerVariable `yaml:"variables,omitempty" json:"variables,omitempty"`
}

// PathItem - https://swagger.io/specification/#pathsItemObject
type PathItem struct {
	Ref         string      `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Summary     string      `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Get         *Operation  `yaml:"get,omitempty" json:"get,omitempty"`
	Post        *Operation  `yaml:"post,omitempty" json:"post,omitempty"`
	Put         *Operation  `yaml:"put,omitempty" json:"put,omitempty"`
	Delete      *Operation  `yaml:"delete,omitempty" json:"delete,omitempty"`
	Options     *Operation  `yaml:"options,omitempty" json:"options,omitempty"`
	Head        *Operation  `yaml:"head,omitempty" json:"head,omitempty"`
	Patch       *Operation  `yaml:"patch,omitempty" json:"patch,omitempty"`
	Trace       *Operation  `yaml:"trace,omitempty" json:"trace,omitempty"`
	Servers     []Server    `yaml:"servers,omitempty" json:"servers,omitempty"`
	Parameters  []Parameter `yaml:"parameters,omitempty" json:"parameters,omitempty"` // or Ref
}

// Components - https://swagger.io/specification/#componentsObject
type Components struct {
	Schemas         map[string]Schema         `yaml:"schemas,omitempty" json:"schemas,omitempty"`                 // or Ref
	Responses       map[string]Response       `yaml:"responses,omitempty" json:"responses,omitempty"`             // or Ref
	Parameters      map[string]Parameter      `yaml:"parameters,omitempty" json:"parameters,omitempty"`           // or Ref
	Examples        map[string]Example        `yaml:"examples,omitempty" json:"examples,omitempty"`               // or Ref
	RequestBodies   map[string]RequestBody    `yaml:"requestBodies,omitempty" json:"requestBodies,omitempty"`     // or Ref
	Headers         map[string]Parameter      `yaml:"headers,omitempty" json:"headers,omitempty"`                 // or Ref
	SecuritySchemes map[string]SecurityScheme `yaml:"securitySchemes,omitempty" json:"securitySchemes,omitempty"` // or Ref
	Links           map[string]Link           `yaml:"links,omitempty" json:"links,omitempty"`                     // or Ref
	Callbacks       map[string]PathItem       `yaml:"callbacks,omitempty" json:"callbacks,omitempty"`             // or Ref
}

// SecurityScheme - https://swagger.io/specification/#securitySchemeObject
type SecurityScheme struct {
	Ref              string      `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Type             string      `yaml:"type,omitempty" json:"type,omitempty"`
	Description      string      `yaml:"description,omitempty" json:"description,omitempty"`
	Name             string      `yaml:"name,omitempty" json:"name,omitempty"`
	In               string      `yaml:"in,omitempty" json:"in,omitempty"`
	Scheme           string      `yaml:"scheme,omitempty" json:"scheme,omitempty"`
	BearerFormat     string      `yaml:"bearerFormat,omitempty" json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `yaml:"flows,omitempty" json:"flows,omitempty"`
	OpenIDConnectURL string      `yaml:"openIdConnectUrl,omitempty" json:"openIdConnectUrl,omitempty"`
}

// OAuthFlows - https://swagger.io/specification/#oauthFlowsObject
type OAuthFlows struct {
	Implicit          *OAuthFlow `yaml:"implicit,omitempty" json:"implicit,omitempty"`
	Password          *OAuthFlow `yaml:"password,omitempty" json:"password,omitempty"`
	ClientCredentials *OAuthFlow `yaml:"clientCredentials,omitempty" json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `yaml:"authorizationCode,omitempty" json:"authorizationCode,omitempty"`
}

// OAuthFlow - https://swagger.io/specification/#oauthFlowObject
type OAuthFlow struct {
	AuthorizationURL string            `yaml:"authorizationUrl,omitempty" json:"authorizationUrl,omitempty"`
	TokenURL         string            `yaml:"tokenUrl,omitempty" json:"tokenUrl,omitempty"`
	RefreshURL       string            `yaml:"refreshUrl,omitempty" json:"refreshUrl,omitempty"`
	Scopes           map[string]string `yaml:"scopes,omitempty" json:"scopes,omitempty"`
}

// Parameter - https://swagger.io/specification/#parameterObject
// Headers - They follow parameter objects
// if this is an "or $ref" we may need to mark all the fields omitempty
type Parameter struct {
	Ref             string               `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Name            string               `yaml:"name,omitempty" json:"name,omitempty"`
	In              string               `yaml:"in,omitempty" json:"in,omitempty"`
	Description     string               `yaml:"description,omitempty" json:"description,omitempty"`
	Required        bool                 `yaml:"required,omitempty" json:"required,omitempty"`
	Deprecated      bool                 `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	AllowEmptyValue bool                 `yaml:"allowEmptyValue,omitempty" json:"allowEmptyValue,omitempty"`
	Style           string               `yaml:"style,omitempty" json:"style,omitempty"`
	Explode         bool                 `yaml:"explode,omitempty" json:"explode,omitempty"`
	AllowReserved   bool                 `yaml:"allowReserved,omitempty" json:"allowReserved,omitempty"`
	Schema          *Schema              `yaml:"schema,omitempty" json:"schema,omitempty"` // or ref
	Example         *interface{}         `yaml:"example,omitempty" json:"example,omitempty"`
	Examples        map[string]Example   `yaml:"examples,omitempty" json:"examples,omitempty"` // or ref
	Content         map[string]MediaType `yaml:"content,omitempty" json:"content,omitempty"`
}

// RequestBody - https://swagger.io/specification/#requestBodyObject
type RequestBody struct {
	Ref         string               `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Description string               `yaml:"description,omitempty" json:"description,omitempty"`
	Content     map[string]MediaType `yaml:"content,omitempty" json:"content,omitempty"`
	Required    bool                 `yaml:"required,omitempty" json:"required,omitempty"`
}

// MediaType - https://swagger.io/specification/#mediaTypeObject
type MediaType struct {
	Schema   *Schema             `yaml:"schema,omitempty" json:"schema,omitempty"` // or ref
	Example  *interface{}        `yaml:"example,omitempty" json:"example,omitempty"`
	Examples map[string]Example  `yaml:"examples,omitempty" json:"examples,omitempty"` // or ref
	Encoding map[string]Encoding `yaml:"encoding,omitempty" json:"encoding,omitempty"`
}

// Encoding - https://swagger.io/specification/#encodingObject
type Encoding struct {
	ContentType   string               `yaml:"contentType,omitempty" json:"contentType,omitempty"`
	Headers       map[string]Parameter `yaml:"headers,omitempty" json:"headers,omitempty"` // or ref
	Style         string               `yaml:"style,omitempty" json:"style,omitempty"`
	Explode       bool                 `yaml:"explode,omitempty" json:"explode,omitempty"`
	AllowReserved bool                 `yaml:"allowReserved,omitempty" json:"allowReserved,omitempty"`
}

// Example - https://swagger.io/specification/#exampleObject
type Example struct {
	Ref           string       `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Summary       string       `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description   string       `yaml:"description,omitempty" json:"description,omitempty"`
	Value         *interface{} `yaml:"value,omitempty" json:"value,omitempty"`
	ExternalValue string       `yaml:"externalValue,omitempty" json:"externalValue,omitempty"`
}

// Operation - https://swagger.io/specification/#operationObject
type Operation struct {
	Tags         []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
	Summary      string                 `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description  string                 `yaml:"description,omitempty" json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	OperationID  string                 `yaml:"operationId,omitempty" json:"operationId,omitempty"`
	Parameters   []Parameter            `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	RequestBody  *RequestBody           `yaml:"requestBody,omitempty" json:"requestBody,omitempty"`
	Responses    map[string]Response    `yaml:"responses,omitempty" json:"responses,omitempty"`
	Callbacks    map[string]PathItem    `yaml:"callbacks,omitempty" json:"callbacks,omitempty"`
	Deprecated   bool                   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Security     map[string][]string    `yaml:"security,omitempty" json:"security,omitempty"`
	Servers      []Server               `yaml:"servers,omitempty" json:"servers,omitempty"`
}

// Schema - https://swagger.io/specification/#schemaObject
type Schema struct {
	Ref string `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	// JSON Schema Fields
	Title            string   `yaml:"title,omitempty" json:"title,omitempty"`
	MultipleOf       *int64   `yaml:"multipleOf,omitempty" json:"multipleOf,omitempty"`
	Maximum          *int64   `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	ExclusiveMaximum *int64   `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	Minimum          *int64   `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	ExclusiveMinimum *int64   `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	MaxLength        *int64   `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	MinLength        *int64   `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	Pattern          string   `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MaxItems         *int64   `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	MinItems         *int64   `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	UniqueItems      bool     `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`
	MaxProperties    *int64   `yaml:"maxProperties,omitempty" json:"maxProperties,omitempty"`
	MinProperties    *int64   `yaml:"minProperties,omitempty" json:"minProperties,omitempty"`
	Required         []string `yaml:"required,omitempty" json:"required,omitempty"`
	Enum             []string `yaml:"enum,omitempty" json:"enum,omitempty"`

	// OpenAPI modified JSON Schema Fields
	Type                 string            `yaml:"type,omitempty" json:"type,omitempty"`
	Description          string            `yaml:"description,omitempty" json:"description,omitempty"`
	Properties           map[string]Schema `yaml:"properties,omitempty" json:"properties,omitempty"`
	AdditionalProperties *Schema           `yaml:"additionalProperties,omitempty" json:"additionalProperties,omitempty"`
	Format               string            `yaml:"format,omitempty" json:"format,omitempty"`
	Default              interface{}       `yaml:"default,omitempty" json:"default,omitempty"`
	Items                *Schema           `yaml:"items,omitempty" json:"items,omitempty"`
	AllOf                []Schema          `yaml:"allOf,omitempty" json:"allOf,omitempty"`
	OneOf                []Schema          `yaml:"oneOf,omitempty" json:"oneOf,omitempty"`
	AnyOf                []Schema          `yaml:"anyOf,omitempty" json:"anyOf,omitempty"`
	Not                  []Schema          `yaml:"not,omitempty" json:"not,omitempty"`

	// Static OpenAPI Fields
	Discriminator *Discriminator         `yaml:"discriminator,omitempty" json:"discriminator,omitempty"`
	Nullable      bool                   `yaml:"nullable,omitempty" json:"nullable,omitempty"`
	ReadOnly      bool                   `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`
	WriteOnly     bool                   `yaml:"writeOnly,omitempty" json:"writeOnly,omitempty"`
	XML           *XML                   `yaml:"xml,omitempty" json:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Example       interface{}            `yaml:"example,omitempty" json:"example,omitempty"`
	Deprecated    bool                   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
}

// Discriminator - https://swagger.io/specification/#discriminatorObject
type Discriminator struct {
	PropertyName string            `yaml:"propertyName,omitempty" json:"propertyName,omitempty"`
	Mapping      map[string]string `yaml:"mapping,omitempty" json:"mapping,omitempty"`
}

// XML - https://swagger.io/specification/#xmlObject
type XML struct {
	Name      string `yaml:"name,omitempty" json:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Prefix    string `yaml:"prefix,omitempty" json:"prefix,omitempty"`
	Attribute bool   `yaml:"attribute,omitempty" json:"attribute,omitempty"`
	Warpped   bool   `yaml:"wrapped,omitempty" json:"wrapped,omitempty"`
}

// Response - https://swagger.io/specification/#responseObject
// if this is an "or $ref" we may need to mark all the fields omitempty
type Response struct {
	Ref         string               `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Description string               `yaml:"description,omitempty" json:"description,omitempty"`
	Headers     map[string]Parameter `yaml:"headers,omitempty" json:"headers,omitempty"`
	Content     map[string]MediaType `yaml:"content,omitempty" json:"content,omitempty"`
	Links       map[string]Link      `yaml:"links,omitempty" json:"links,omitempty"`
}

// Link - https://swagger.io/specification/#linkObject
// if this is an "or $ref" we may need to mark all the fields omitempty
// expression: https://swagger.io/specification/#runtimeExpression
type Link struct {
	Ref          string            `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	OperationRef string            `yaml:"operationRef,omitempty" json:"operationRef,omitempty"`
	OperationID  string            `yaml:"operationID,omitempty" json:"operationID,omitempty"`
	Parameters   map[string]string `yaml:"parameters,omitempty" json:"parameters,omitempty"`   //expression
	RequestBody  string            `yaml:"requestBody,omitempty" json:"requestBody,omitempty"` // expression
	Description  string            `yaml:"description,omitempty" json:"description,omitempty"`
	Server       *Server           `yaml:"server,omitempty" json:"server,omitempty"`
}

// Tag - https://swagger.io/specification/#tagObject
type Tag struct {
	Name         string                 `yaml:"name" json:"name"`
	Description  string                 `yaml:"description,omitempty" json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

// ExternalDocumentation - https://swagger.io/specification/#externalDocumentationObject
type ExternalDocumentation struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	URL         string `yaml:"url" json:"url"`
}

// ServerVariable - https://swagger.io/specification/#serverVariableObject
type ServerVariable struct {
	Enum        []string `yaml:"enum,omitempty" json:"enum,omitempty"`
	Default     string   `yaml:"default" json:"default"`
	Description string   `yaml:"description" json:"description"`
}
