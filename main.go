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

var (
	url   string
	skips = map[string]bool{
		"root":           true,
		"self":           true,
		"subscribe":      true, // action
		"shell":          true, // open shell
		"yaml":           true, // export link
		"icon":           true, // export link
		"readme":         true, // export link
		"app-readme":     true, // export link
		"exportYaml":     true, // export link
		"authConfigs":    true, // Objects don't match schema/are not under the collection
		"dynamicSchemas": true, // Not sure what this is for yet and its a schema its self
		"ldapConfigs":    true, // 404 - No collection.
	}
)

// Collections -
type Collections struct {
	Links map[string]string `json:"links"`
}

// Collection -
type Collection struct {
	*norman.Collection
	Data []norman.Resource `json:"data"`
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
	url, ok := os.LookupEnv("RANCHER_URL")
	if !ok {
		log.Fatal("Set RANCHER_URL")
	}
	log.Debug("Import descriptions")
	descriptions := make(map[string]string)
	yamlDescriptions, err := ioutil.ReadFile("./data/descriptions.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlDescriptions, descriptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Import base")
	swagger := &openapi.OpenAPI{}
	yamlFile, err := ioutil.ReadFile("./data/base.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, swagger)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Initialize swagger maps")
	swagger.Paths = make(map[string]openapi.PathItem)
	swagger.Components.Parameters = make(map[string]openapi.Parameter)

	log.Debug("Get Root Collections")
	collections, err := getCollections(url)
	if err != nil {
		log.Fatal(err)
	}

	for col, link := range collections {
		// Only follow a specific root collection
		only, ok := os.LookupEnv("COLLECTION")
		if ok {
			if col != only {
				log.Debug("Single Collection Only: ", col)
				continue
			}
		}

		err = parseCollection(col, link, "/", url, swagger)
		if err != nil {
			log.Warnf("Failed to parse %s, %s, %s - %v", col, link, "/", err)
		}
	}

	// Render swagger doc
	out, err := json.Marshal(swagger)
	log.Debug(string(out))
	err = ioutil.WriteFile("./build/swagger.json", out, 0644)
}

func parseCollection(col string, link string, base string, url string, swagger *openapi.OpenAPI) error {
	// Skip "weird/broken" collections
	if skips[col] {
		log.Debug("Skipped: ", col)
		return nil
	}
	log.Infof("Parse Collection: %s -> %s - %s", col, link, base)

	collection, err := getCollection(link)
	if err != nil {
		log.Errorf("Failed to get collection: %v", err)
		return err
	}

	if collection.Type != "collection" {
		log.Debugf("%s is not a collection, skipping - %s %s", col, link, base)
		return nil
	}

	parameters := make([]openapi.Parameter, 0)

	// (╯°□°）╯︵ ┻━┻ some schemas are under the /{collection}/{id}/schemas
	// Base schema path on createTypes.
	schemaRootRegex := regexp.MustCompile(fmt.Sprintf("(?i)^%s(.*)/%s$", url, col))
	_, ok := collection.CreateTypes[collection.ResourceType]
	if !ok {
		return fmt.Errorf("%s, Collection doesn't have CreateTypes", collection.ResourceType)
	}
	schemaRootSlice := schemaRootRegex.FindStringSubmatch(collection.CreateTypes[collection.ResourceType])
	log.Debugf("schema return: %v", schemaRootSlice)
	schemaRoot := schemaRootSlice[1]

	log.Debug("resourceType for collection: ", collection.ResourceType)
	rSchema, err := getSchema(fmt.Sprintf("%v%s/schemas/%s", url, schemaRoot, collection.ResourceType))
	if err != nil {
		return fmt.Errorf("Failed to get Schema for %s/%s - %v", schemaRoot, collection.ResourceType, err)
	}

	// populate swagger schema objects
	translateSchema(rSchema, url, swagger)

	// set schema for collection
	createCollectionSchema(col, collection, swagger)

	// set previous parameters
	searchPrams := regexp.MustCompile("{(\\w+)}")
	previousPrams := searchPrams.FindAllStringSubmatch(base, -1)
	if len(previousPrams) > 0 {
		// previousPrams = previousPrams[1:]
		for _, p := range previousPrams {
			param := openapi.Parameter{
				Ref: fmt.Sprintf("#/components/parameters/%s", p[1]),
			}
			parameters = append(parameters, param)
		}
	}

	// /{collection}
	colParameters := make([]openapi.Parameter, 0)
	if len(parameters) > 0 {
		colParameters = parameters
	} else {
		colParameters = nil
	}
	colPathItem := openapi.PathItem{
		Parameters: colParameters,
	}

	for _, method := range rSchema.CollectionMethods {
		if method == "GET" {
			colPathItem.Get = createCollection("GET", col, collection)
		} else if method == "POST" {
			colPathItem.Post = createCollection("POST", col, collection)
		} else {
			log.Error("Unknown Collection Method: ", method)
		}
	}
	swagger.Paths[fmt.Sprintf("%s%s", base, col)] = colPathItem

	// Resource /{collection}/{id}
	newPramID := fmt.Sprintf("%sId", collection.ResourceType)
	createPathParameter(newPramID, swagger)
	param := openapi.Parameter{
		Ref: fmt.Sprintf("#/components/parameters/%s", newPramID),
	}
	parameters = append(parameters, param)

	resourcePathItem := openapi.PathItem{
		Parameters: parameters,
	}
	for _, method := range rSchema.ResourceMethods {
		if method == "GET" {
			resourcePathItem.Get = createResource("GET", collection.ResourceType)
		} else if method == "PUT" {
			resourcePathItem.Put = createResource("PUT", collection.ResourceType)
		} else if method == "DELETE" {
			resourcePathItem.Delete = createResource("DELETE", collection.ResourceType)
		} else {
			log.Error("Unknown Resource Method: ", method)
		}
	}
	swagger.Paths[fmt.Sprintf("%s%s/{%s}", base, col, newPramID)] = resourcePathItem

	if len(collection.Data) > 0 {
		// take the first one
		subBase := fmt.Sprintf("%s%s/{%s}/", base, col, newPramID)

		for subCol, subColLink := range collection.Data[0].Links {
			if subColLink != collection.Data[0].Links["self"] {
				err = parseCollection(subCol, subColLink, subBase, url, swagger)
				if err != nil {
					log.Warnf("Failed to parse %s, %s, %s - %v", subCol, subColLink, subBase, err)
				}
			}
		}
	}
	return nil
}

func translateSchema(rancherSchema norman.Schema, url string, swagger *openapi.OpenAPI) {
	properties := make(map[string]openapi.Schema)
	required := make([]string, 0)
	name := rancherSchema.ID
	resourceFields := rancherSchema.ResourceFields

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
			Default:   resourceValue.Default,
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
		if resourceValue.Create == false && resourceValue.Update == false {
			p.ReadOnly = true
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

		// regex to find Schema base path
		findSchemaBase := regexp.MustCompile("^/v3([/\\w]*)")

		// remap types to valid types for openapi
		isValid := regexp.MustCompile("^(string|boolean|object|array)$")
		isIntOrString := regexp.MustCompile("^intOrString$")
		isArrayString := regexp.MustCompile("^array\\[string\\]$")
		isArrayInt := regexp.MustCompile("^array\\[int\\]$")
		isArrayEnum := regexp.MustCompile("^array\\[enum\\]$")
		isRefArray := regexp.MustCompile("^array\\[(\\w+)\\]$")
		isMapString := regexp.MustCompile("^map\\[string\\]$")
		isMapBase64 := regexp.MustCompile("^map\\[base64\\]$")
		isRefMap := regexp.MustCompile("^map\\[(\\w+)\\]$")
		isRefID := regexp.MustCompile("^reference\\[([a-zA-Z0-9/]+)\\]$")
		isArrayRefID := regexp.MustCompile("^array\\[reference\\[([a-zA-Z0-9/]+)\\]\\]$")
		isEnum := regexp.MustCompile("^enum$")
		isDNSLabel := regexp.MustCompile("^(dnsLabel|hostname|dnsLabelRestricted)$")
		isDate := regexp.MustCompile("^date$")
		isPassword := regexp.MustCompile("^password$")
		isInt := regexp.MustCompile("^int$")
		isBase64 := regexp.MustCompile("^base64$")

		switch {
		case isValid.MatchString(resourceValue.Type):
			p.Type = resourceValue.Type

		case isDate.MatchString(resourceValue.Type):
			p.Type = "string"
			p.Format = "date-time"

		case isPassword.MatchString(resourceValue.Type):
			p.Type = "string"
			p.Format = "password"

		case isIntOrString.MatchString(resourceValue.Type):
			p.OneOf = []openapi.Schema{
				openapi.Schema{
					Type: "string",
				},
				openapi.Schema{
					Type: "integer",
				},
			}

		case isInt.MatchString(resourceValue.Type):
			p.Type = "integer"

		case isBase64.MatchString(resourceValue.Type):
			p.Type = "string"
			desc = append(desc, "Base64 encoded string")

		case isEnum.MatchString(resourceValue.Type):
			p.Type = "string"

		case isDNSLabel.MatchString(resourceValue.Type):
			p.Type = "string"
			p.Pattern = "^(\\w|[A-Za-z0-9-\\.]*\\w)$"
			desc = append(desc, "Must be valid Hostname")

		case isArrayString.MatchString(resourceValue.Type):
			p.Type = "array"
			desc = append(desc, "Array of Strings")
			p.Items = &openapi.Schema{
				Type: "string",
			}

		case isArrayInt.MatchString(resourceValue.Type):
			p.Type = "array"
			desc = append(desc, "Array of Integers")
			p.Items = &openapi.Schema{
				Type: "integer",
			}

		case isArrayEnum.MatchString(resourceValue.Type):
			p.Type = "array"
			p.Items = &openapi.Schema{
				Type: "string",
			}
			desc = append(desc, "Array of Valid Options")

		case isRefArray.MatchString(resourceValue.Type):
			// Will be a ref to other resource, resolve the other resource
			refSchemaName := isRefArray.FindStringSubmatch(resourceValue.Type)[1]
			schemaBase := findSchemaBase.FindStringSubmatch(rancherSchema.Version.Path)[1]
			subSchema, err := getSchema(fmt.Sprintf("%s%s/schemas/%s", url, schemaBase, refSchemaName))
			if err != nil {
				log.Errorf("Failed to get Schema for base:%s name:%s ref:%s url:%s - %v", schemaBase, name, refSchemaName, url, err)
			} else {
				// turtles all the way down
				translateSchema(subSchema, url, swagger)

				p.Type = "array"
				p.Items = &openapi.Schema{
					Ref: fmt.Sprintf("#/components/schemas/%s", subSchema.ID),
				}
			}

		case isMapString.MatchString(resourceValue.Type):
			p.Type = "object"
			example := make(map[string]string)
			example["key"] = "value"
			p.Example = example

		case isMapBase64.MatchString(resourceValue.Type):
			p.Type = "object"
			example := make(map[string]string)
			example["key"] = "base64 encoded string"
			p.Example = example

		case isRefMap.MatchString(resourceValue.Type):
			refSchemaName := isRefMap.FindStringSubmatch(resourceValue.Type)[1]
			schemaBase := findSchemaBase.FindStringSubmatch(rancherSchema.Version.Path)[1]
			subSchema, err := getSchema(fmt.Sprintf("%s%s/schemas/%s", url, schemaBase, refSchemaName))
			if err != nil {
				log.Errorf("Failed to get Schema for base:%s name:%s ref:%s url:%s - %v", schemaBase, name, refSchemaName, url, err)
			} else {
				// turtles all the way down
				translateSchema(subSchema, url, swagger)

				p.Type = "object"
				p.AdditionalProperties = &openapi.Schema{
					Ref: fmt.Sprintf("#/components/schemas/%s", subSchema.ID),
				}
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

		default:
			// Should be schema object
			schemaBase := findSchemaBase.FindStringSubmatch(rancherSchema.Version.Path)[1]

			subSchema, err := getSchema(fmt.Sprintf("%s%s/schemas/%s", url, schemaBase, resourceValue.Type))
			if err != nil {
				log.Errorf("Failed to get Schema for base:%s name:%s ref:%s url:%s - %v", schemaBase, name, resourceValue.Type, url, err)
			} else {

				// turtles all the way down
				translateSchema(subSchema, url, swagger)

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
				if subSchema.ID == "" {
					log.Error("id is empty")
				}
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

func createResource(method string, resourceType string) *openapi.Operation {
	schema := &openapi.Schema{
		Ref: fmt.Sprintf("#/components/schemas/%s", resourceType),
	}

	content := make(map[string]openapi.MediaType)
	content["application/json"] = openapi.MediaType{
		Schema: schema,
	}

	resp := make(map[string]openapi.Response)
	request := &openapi.RequestBody{}

	if method == "GET" {
		request = nil
		resp["200"] = openapi.Response{
			Description: fmt.Sprintf("Returns '%s' object.", resourceType),
			Content:     content,
		}
	}
	if method == "PUT" {
		request.Description = fmt.Sprintf("Update `%s` object.", resourceType)
		request.Content = content
		resp["200"] = openapi.Response{
			Description: fmt.Sprintf("Returns '%s' object.", resourceType),
			Content:     content,
		}
	}
	if method == "DELETE" {
		request = nil
		resp["204"] = openapi.Response{
			Description: fmt.Sprint("Delete Successful"),
		}
	}

	return &openapi.Operation{
		Description: fmt.Sprintf("`%s` Resource", resourceType),
		Responses:   resp,
		RequestBody: request,
	}
}

func createCollection(method string, col string, collection *Collection) *openapi.Operation {
	schema := &openapi.Schema{
		Ref: fmt.Sprintf("#/components/schemas/%s", collection.ResourceType),
	}

	content := make(map[string]openapi.MediaType)
	content["application/json"] = openapi.MediaType{
		Schema: schema,
	}

	request := &openapi.RequestBody{}
	resp := make(map[string]openapi.Response)

	if method == "GET" {
		request = nil
		resp["200"] = openapi.Response{
			Description: fmt.Sprintf("Returns list of '%s'", col),
			Content:     content,
		}
	}
	if method == "POST" {
		request.Description = fmt.Sprintf("Create a new `%s` object.", collection.ResourceType)
		request.Content = content
		resp["200"] = openapi.Response{
			Description: fmt.Sprintf("Returns new `%s` object.", collection.ResourceType),
			Content:     content,
		}
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

func createCollectionSchema(name string, collection *Collection, swagger *openapi.OpenAPI) {
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

func getCollection(link string) (*Collection, error) {
	collectionResponse, err := httpGet(link)
	if err != nil {
		return nil, err
	}

	collection := &Collection{}
	err = json.Unmarshal(collectionResponse, collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func httpGet(url string) ([]byte, error) {
	token := os.Getenv("RANCHER_TOKEN")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	goodStatus := regexp.MustCompile("^2\\d\\d")
	if !goodStatus.MatchString(resp.Status) {
		return nil, fmt.Errorf("%s %s returned %s", req.Method, req.URL, resp.Status)
	}

	defer resp.Body.Close()
	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return jsonBody, nil
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
