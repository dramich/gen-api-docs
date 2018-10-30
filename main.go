package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	openapi "github.com/rancher/gen-api-docs/openapi/v3.0.1"
	norman "github.com/rancher/norman/types"
	log "github.com/sirupsen/logrus"
)

const (
	url = "https://rancher.localhost/v3"
)

var (
	skips = map[string]bool{
		"root":           true,
		"self":           true,
		"subscribe":      true,
		"authConfigs":    true, // Objects don't match schema/are not under the collection
		"dynamicSchemas": true, // Not sure what this is for yet and its a schema in its self
		"ldapConfigs":    true, // No schema for this collection.
	}
)

// Collections -
type Collections struct {
	Links map[string]string `json:"links"`
}

func init() {
	val, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		level, _ := log.ParseLevel(val)
		log.SetLevel(level)
	}
}

// Main
func main() {
	// import descriptions
	descriptions := make(map[string]string)
	yamlDescriptions, err := ioutil.ReadFile("./data/descriptions.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlDescriptions, descriptions)
	if err != nil {
		log.Fatal(err)
	}

	// import base
	swagger := &openapi.OpenAPI{}
	yamlFile, err := ioutil.ReadFile("./data/base.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, swagger)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize swagger maps
	swagger.Paths = make(map[string]openapi.PathItem)
	swagger.Components.Parameters = make(map[string]openapi.Parameter)
	// swagger.Components.Responses = make(map[string]openapi.Response)

	// get all the collections
	collections, err := getCollections(url)
	if err != nil {
		log.Fatal(err)
	}

	// iterate through the collections
	for col, link := range collections {
		// Only specified collection
		only, ok := os.LookupEnv("COLLECTION")
		if ok {
			if col != only {
				continue
			}
		}
		// Skip "weird" collections
		if skips[col] {
			log.Debug("Skipped: ", col)
			continue
		}

		log.Infof("%s -> %s", col, link)
		collection, err := getCollection(link)
		if err != nil {
			log.Errorf("Failed to get collection Resource Type: %v", err)
		}

		log.Debug("resourceType for collection: ", collection.ResourceType)
		rSchema, err := getSchema(fmt.Sprintf("%v/schemas/%v", url, collection.ResourceType))
		if err != nil {
			log.Errorf("Failed to get Schema for %v - %v", collection.ResourceType, err)
		}

		// populate swagger schema objects
		translateSchema(rSchema.ResourceFields, rSchema.ID, swagger)

		// set schema for collection
		setCollectionSchema(col, collection, swagger)

		// populate swagger path objects
		// /{collection}
		colPathItem := openapi.PathItem{}
		for _, method := range rSchema.CollectionMethods {
			if method == "GET" {
				colPathItem.Get = createCollectionGet(col)
			} else if method == "POST" {
				colPathItem.Post = createCollectionPost(col, collection)
			} else {
				log.Error("Unknown Collection Method: ", method)
			}
		}
		swagger.Paths[fmt.Sprint("/", col)] = colPathItem

		// Resource /{collection}/{id}
		createPathParameter("id", swagger)
		param := openapi.Parameter{
			Ref: "#/components/parameters/id",
		}
		params := make([]openapi.Parameter, 0)
		params = append(params, param)
		resourcePathItem := openapi.PathItem{
			Parameters: params,
		}
		for _, method := range rSchema.ResourceMethods {
			if method == "GET" {
				resourcePathItem.Get = createResourceGet(collection.ResourceType)
			} else if method == "PUT" {
				resourcePathItem.Put = createResourcePut(collection.ResourceType)
			} else if method == "DELETE" {
				resourcePathItem.Delete = createResourceDelete(collection.ResourceType)
			} else {
				log.Error("Unknown Resource Method: ", method)
			}
		}
		swagger.Paths[fmt.Sprintf("/%s/{id}", col)] = resourcePathItem
	}

	// scrape links for

	// Render swagger doc
	out, err := yaml.Marshal(swagger)
	log.Debug(string(out))
	err = ioutil.WriteFile("./swagger/swagger.out.yml", out, 0644)

}

func createResourceDelete(resourceType string) *openapi.Operation {
	resp := make(map[string]openapi.Response)
	resp["204"] = openapi.Response{
		Description: fmt.Sprint("Delete Successful"),
	}

	return &openapi.Operation{
		Description: fmt.Sprintf("Delete `%s` Resource", resourceType),
		Responses:   resp,
	}
}

func createResourceGet(resourceType string) *openapi.Operation {
	schema := &openapi.Schema{
		Ref: fmt.Sprintf("#/components/schemas/%s", resourceType),
	}

	content := make(map[string]openapi.MediaType)
	content["application/json"] = openapi.MediaType{
		Schema: schema,
	}

	resp := make(map[string]openapi.Response)
	resp["200"] = openapi.Response{
		Description: fmt.Sprintf("Returns '%s' object.", resourceType),
		Content:     content,
	}

	parameters := make([]openapi.Parameter, 0)
	p := openapi.Parameter{
		Ref: fmt.Sprintf("#/components/parameters/id"),
	}
	parameters = append(parameters, p)

	return &openapi.Operation{
		Description: fmt.Sprintf("`%s` Resource", resourceType),
		Responses:   resp,
		Parameters:  parameters,
	}
}

func createResourcePut(resourceType string) *openapi.Operation {
	schema := &openapi.Schema{
		Ref: fmt.Sprintf("#/components/schemas/%s", resourceType),
	}

	content := make(map[string]openapi.MediaType)
	content["application/json"] = openapi.MediaType{
		Schema: schema,
	}

	resp := make(map[string]openapi.Response)
	resp["200"] = openapi.Response{
		Description: fmt.Sprintf("Returns '%s' object.", resourceType),
		Content:     content,
	}

	request := &openapi.RequestBody{
		Description: fmt.Sprintf("Update `%s` object.", resourceType),
		Content:     content,
	}

	parameters := make([]openapi.Parameter, 0)
	p := openapi.Parameter{
		Ref: fmt.Sprintf("#/components/parameters/id"),
	}
	parameters = append(parameters, p)

	return &openapi.Operation{
		Description: fmt.Sprintf("`%s` Resource", resourceType),
		Responses:   resp,
		RequestBody: request,
		Parameters:  parameters,
	}
}

func createCollectionGet(col string) *openapi.Operation {
	schema := &openapi.Schema{
		Ref: fmt.Sprintf("#/components/schemas/%s", col),
	}

	content := make(map[string]openapi.MediaType)
	content["application/json"] = openapi.MediaType{
		Schema: schema,
	}

	resp := make(map[string]openapi.Response)
	resp["200"] = openapi.Response{
		Description: fmt.Sprintf("Returns list of '%s'", col),
		Content:     content,
	}

	return &openapi.Operation{
		Description: fmt.Sprintf("`%s` Collection", col),
		Responses:   resp,
	}
}

func createCollectionPost(col string, collection *norman.Collection) *openapi.Operation {
	schema := &openapi.Schema{
		Ref: fmt.Sprintf("#/components/schemas/%s", collection.ResourceType),
	}

	content := make(map[string]openapi.MediaType)
	content["application/json"] = openapi.MediaType{
		Schema: schema,
	}

	request := &openapi.RequestBody{
		Description: fmt.Sprintf("Create a new `%s` object.", collection.ResourceType),
		Content:     content,
	}

	resp := make(map[string]openapi.Response)
	resp["200"] = openapi.Response{
		Description: fmt.Sprintf("Returns new `%s` object.", collection.ResourceType),
		Content:     content,
	}

	return &openapi.Operation{
		Description: fmt.Sprintf("`%s` Collection", col),
		Responses:   resp,
		RequestBody: request,
	}

}

func createPathParameter(name string, swagger *openapi.OpenAPI) {
	_, ok := swagger.Components.Parameters[name]
	if !ok {
		swagger.Components.Parameters[name] = openapi.Parameter{
			Name:     name,
			In:       "path",
			Required: true,
			Schema: &openapi.Schema{
				Type: "string",
			},
		}
	}
}

func setCollectionSchema(name string, collection *norman.Collection, swagger *openapi.OpenAPI) {
	log.Debug(collection.Type)
	colAllOf := make([]openapi.Schema, 0)
	resTypeRef := openapi.Schema{
		Ref: "#/components/schemas/collection",
	}
	colAllOf = append(colAllOf, resTypeRef)
	colData := openapi.Schema{
		Type: "array",
		Items: &openapi.Schema{
			Ref: fmt.Sprintf("#/components/schemas/%s", collection.ResourceType),
		},
	}
	colProp := make(map[string]openapi.Schema)
	colProp["data"] = colData

	colSchema := openapi.Schema{
		Type:       "object",
		AllOf:      colAllOf,
		Properties: colProp,
	}
	swagger.Components.Schemas[name] = colSchema
}

func translateSchema(resourceFields map[string]norman.Field, name string, swagger *openapi.OpenAPI) {
	properties := make(map[string]openapi.Schema)
	required := make([]string, 0)

	// Skip Schema if it already exists
	_, ok := swagger.Components.Schemas[name]
	if ok {
		log.Debug(name, " Schema Already Exists")
		return
	}

	// Required
	for resourceName, resourceValue := range resourceFields {
		if resourceValue.Required {
			required = append(required, resourceName)
		}
	}
	// Properties
	for resourceName, resourceValue := range resourceFields {
		desc := make([]string, 0)

		p := &openapi.Schema{
			Default: resourceValue.Default,
			// Description: resourceValue.Description,
			Enum:      resourceValue.Options,
			Maximum:   resourceValue.Max,
			MaxLength: resourceValue.MaxLength,
			Minimum:   resourceValue.Min,
			MinLength: resourceValue.MinLength,
			Pattern:   resourceValue.ValidChars,
			Nullable:  resourceValue.Nullable,
			// other values that I'm not sure what to do with yet.
			// resourceValue.CodeName
			// resourceValue.DynamicField
			// resourceValue.InvalidChars
		}

		usage := ""
		// Skip read-only Properties
		if resourceValue.Create == false && resourceValue.Update == false {
			p.ReadOnly = true
			// continue
		}
		if resourceValue.Update || resourceValue.Create {
			usage = fmt.Sprint("Allowed in Methods:")
		}
		if resourceValue.Create {
			usage = fmt.Sprint(usage, " `POST`")
		}
		if resourceValue.Update {
			usage = fmt.Sprint(usage, " `PUT`")
		}
		if usage != "" {
			desc = append(desc, usage)
		}

		// Populate existing Description
		if resourceValue.Description != "" {
			desc = append(desc, resourceValue.Description)
		}

		// remap types to valid types for openapi
		isValid := regexp.MustCompile("^(string|boolean|object|array)$")
		isArrayString := regexp.MustCompile("^array\\[string\\]$")
		isArrayEnum := regexp.MustCompile("^array\\[enum\\]$")
		isRefArray := regexp.MustCompile("^array\\[(\\w+)\\]$")
		isMapString := regexp.MustCompile("^map\\[string\\]$")
		isRefMap := regexp.MustCompile("^map\\[(\\w+)\\]$")
		isRefID := regexp.MustCompile("^reference\\[([a-zA-Z0-9/]+)\\]$")
		isArrayRefID := regexp.MustCompile("^array\\[reference\\[([a-zA-Z0-9/]+)\\]\\]$")
		isEnum := regexp.MustCompile("^enum$")
		isDNSLabel := regexp.MustCompile("^dnsLabel$")
		isDate := regexp.MustCompile("^date$")
		isPassword := regexp.MustCompile("^password$")
		isInt := regexp.MustCompile("^int$")

		switch {
		case isValid.MatchString(resourceValue.Type):
			p.Type = resourceValue.Type

		case isDate.MatchString(resourceValue.Type):
			p.Type = "string"
			p.Format = "date-time"

		case isPassword.MatchString(resourceValue.Type):
			p.Type = "string"
			p.Format = "password"

		case isInt.MatchString(resourceValue.Type):
			p.Type = "integer"

		case isArrayString.MatchString(resourceValue.Type):
			p.Type = "array"
			p.Description = "Array of Strings"
			p.Items = &openapi.Schema{
				Type: "string",
			}

		case isArrayEnum.MatchString(resourceValue.Type):
			p.Type = "array"
			p.Items = &openapi.Schema{
				Type: "string",
			}
			p.Description = "Array of Valid Options"

		case isRefArray.MatchString(resourceValue.Type):
			// Will be a ref to other resource, resolve the other resource
			refSchemaName := isRefArray.FindStringSubmatch(resourceValue.Type)[1]
			subSchema, err := getSchema(fmt.Sprintf("%v/schemas/%v", url, refSchemaName))
			if err != nil {
				log.Errorf("Failed to get Schema for %v", refSchemaName)
				log.Error(err)
			}
			// turtles all the way down
			subSchemaName := fmt.Sprintf("%s", subSchema.ID)
			translateSchema(subSchema.ResourceFields, subSchemaName, swagger)

			p.Type = "array"
			p.Items = &openapi.Schema{
				Ref: fmt.Sprintf("#/components/schemas/%s", subSchema.ID),
			}

		case isMapString.MatchString(resourceValue.Type):
			p.Type = "object"
			example := make(map[string]string)
			example["key"] = "value"
			p.Example = example

		case isRefMap.MatchString(resourceValue.Type):
			refSchemaName := isRefMap.FindStringSubmatch(resourceValue.Type)[1]
			subSchema, err := getSchema(fmt.Sprintf("%v/schemas/%v", url, refSchemaName))
			if err != nil {
				log.Errorf("Failed to get Schema for %v", refSchemaName)
				log.Error(err)
			}
			// turtles all the way down
			subSchemaName := fmt.Sprintf("%s", subSchema.ID)
			translateSchema(subSchema.ResourceFields, subSchemaName, swagger)

			p.Type = "object"
			p.AdditionalProperties = &openapi.Schema{
				Ref: fmt.Sprintf("#/components/schemas/%s", subSchema.ID),
			}

		case isRefID.MatchString(resourceValue.Type):
			ref := isRefID.FindStringSubmatch(resourceValue.Type)[1]
			p.Type = "string"
			desc = append(desc, fmt.Sprintf("Id of %s", ref))

		case isArrayRefID.MatchString(resourceValue.Type):
			ref := isArrayRefID.FindStringSubmatch(resourceValue.Type)[1]
			p.Type = "array"
			p.Items = &openapi.Schema{
				Type: "string",
			}
			desc = append(desc, fmt.Sprintf("Array of Ids of %s", ref))

		case isEnum.MatchString(resourceValue.Type):
			p.Type = "string"

		case isDNSLabel.MatchString(resourceValue.Type):
			p.Type = "string"
			p.Pattern = "^(\\w|[A-Za-z0-9-\\.]*\\w)$"
			desc = append(desc, "Must be valid Hostname")

		default:
			// Should be schema object
			subSchema, err := getSchema(fmt.Sprintf("%v/schemas/%v", url, resourceValue.Type))
			if err != nil {
				log.Errorf("Failed to get Schema for %v", resourceValue.Type)
				log.Error(err)
			} else {
				// turtles all the way down
				subSchemaName := fmt.Sprintf("%s", subSchema.ID)
				translateSchema(subSchema.ResourceFields, subSchemaName, swagger)

				// reset other fields
				p.Default = nil
				p.Description = ""
				p.Enum = nil
				p.Maximum = nil
				p.MaxLength = nil
				p.Minimum = nil
				p.MinLength = nil
				p.Pattern = ""
				p.Nullable = false
				p.Ref = fmt.Sprintf("#/components/schemas/%s", subSchema.ID)
				desc = []string{}
			}
		}

		p.Description = strings.Join(desc, "; ")

		properties[resourceName] = *p
	}

	schemaObject := openapi.Schema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}

	swagger.Components.Schemas[name] = schemaObject

}

func httpGet(url string) ([]byte, error) {

	username := os.Getenv("RANCHER_ACCESS_KEY")
	password := os.Getenv("RANCHER_SECRET_KEY")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return jsonBody, nil
}

func getSchema(link string) (norman.Schema, error) {
	schema := norman.Schema{}

	schemaResponse, err := httpGet(link)
	if err != nil {
		return schema, err
	}

	err = json.Unmarshal(schemaResponse, &schema)
	if err != nil {
		log.Error(string(schemaResponse))
		log.Error(err)
		return schema, err
	}

	return schema, nil
}

func getCollections(link string) (map[string]string, error) {
	collectionsResponse, err := httpGet(link)
	if err != nil {
		return nil, err
	}

	collections := Collections{}
	err = json.Unmarshal(collectionsResponse, &collections)
	if err != nil {
		return nil, err
	}

	return collections.Links, nil
}

func getCollection(link string) (*norman.Collection, error) {
	collectionResponse, err := httpGet(link)
	if err != nil {
		return nil, err
	}

	collection := &norman.Collection{}
	err = json.Unmarshal(collectionResponse, collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func printPretty(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonData)
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
