package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"startpoint/core/client/validator"
	"startpoint/core/model"
	"startpoint/core/writer"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/renderer"
	"github.com/pb33f/libopenapi/utils"
	"github.com/rs/zerolog/log"
)

func ReadSpec(path string, workspace string) {

	// load an OpenAPI 3 specification from bytes
	specBytes, err := loadSpec(path)
	if err != nil {
		panic(err)
	}

	// create a new document from specification bytes
	document, err := libopenapi.NewDocument(specBytes)

	// if anything went wrong, an error is thrown
	if err != nil {
		panic(fmt.Sprintf("cannot create new document: %e", err))
	}

	var (
		requests []model.RequestMold
		profiles []model.Profile
	)
	if document.GetSpecInfo().SpecType == utils.OpenApi3 {
		requests, profiles = importOpenAPIV3(document, workspace)
	} else if document.GetSpecInfo().SpecType == utils.OpenApi2 {
		// TODO:
		fmt.Print("Importing OpenAPI v2 not yet implemented\n")
	}

	fmt.Printf("Processed entries into %d requests and %d profiles. Next going to save these.\n\n", len(requests), len(profiles))

	for _, profile := range profiles {
		path := filepath.Join(profile.Root, profile.Filename)
		contents := profile.AsDotEnv()
		_, err := writer.WriteFile(path, contents)
		status := "OK"
		if err != nil {
			log.Error().Err(err).Msg("Failed to save profile")
			status = "ERROR"
		}
		fmt.Printf("[%s] %s\n", status, path)
	}

	for _, request := range requests {
		path := filepath.Join(request.Root, request.Filename)
		contents := request.Raw()
		_, err := writer.WriteFile(path, contents)
		status := "OK"
		if err != nil {
			log.Error().Err(err).Msg("Failed to save request")
			status = "ERROR"
		}
		fmt.Printf("[%s] %s\n", status, path)
	}

}

func importOpenAPIV3(document libopenapi.Document, workspace string) ([]model.RequestMold, []model.Profile) {

	// TODO: authentication
	// basic: add Authorization: Basic {basic_auth} header
	// bearer: add Authorization: Bearer {bearer_token} header
	// oauth2: add Authorization: Bearer {bearer_token} header
	// - also with some flows, it might be possible to create prev_req (see https://swagger.io/docs/specification/authentication/oauth2/)
	// openid: add Authorization: Bearer {bearer_token} header
	// - https://swagger.io/docs/specification/authentication/openid-connect-discovery/
	// - kind of tough nut: requires to load openid-configuration and parse it...
	// cookie: add Cookie header with a cookie named XYZ (defined in spec)
	// - https://swagger.io/docs/specification/authentication/cookie-authentication/
	// api key: add header with name XYZ (defined in spec) and {api_key} var
	// - https://swagger.io/docs/specification/authentication/api-keys/

	v3Model, errors := document.BuildV3Model()

	// if anything went wrong when building the v3 model, a slice of errors will be returned
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported",
			len(errors)))
	}

	// get a count of the number of paths and schemas.
	paths := v3Model.Model.Paths.PathItems
	schemas := v3Model.Model.Components.Schemas
	servers := v3Model.Model.Servers
	securitySchemes := v3Model.Model.Components.SecuritySchemes
	jsonMockGenerator := renderer.NewMockGenerator(renderer.JSON)

	// print the number of paths and schemas in the document
	fmt.Printf("There are %d paths, %d schemas and %d servers in the document\n", paths.Len(), schemas.Len(), len(servers))

	fmt.Print("\n")

	profiles := handleServersV3(servers, workspace)

	requests := []model.RequestMold{}
	for pathPairs := paths.First(); pathPairs != nil; pathPairs = pathPairs.Next() {
		pathItem := pathPairs.Value()
		// TODO: handle server override from path
		// operationValue.Servers
		// maybe in this case, just put value (if single) as operation url?

		for operation := pathItem.GetOperations().First(); operation != nil; operation = operation.Next() {

			// TODO: maybe this should come as a flag so that user can define whether to generate yaml or other kind of requests?
			yamlRequest := model.YamlRequest{
				Url:     fmt.Sprintf("{url}%s", pathPairs.Key()),
				Method:  strings.ToUpper(operation.Key()),
				Headers: model.Headers{},
			}

			operationValue := operation.Value()

			headers, queryParams, variables := handleOperationParametersV3(operationValue.Parameters, jsonMockGenerator)

			for _, sec := range operationValue.Security {

				requirements := sec.Requirements
				for reqPairs := requirements.First(); reqPairs != nil; reqPairs = reqPairs.Next() {
					securityScheme := findSecurityScheme(reqPairs.Key(), securitySchemes)
					if securityScheme.Scheme == "basic" {
						headers["Authorization"] = []string{fmt.Sprintf("Basic {%s}", reqPairs.Key())}
						variables[reqPairs.Key()] = ""
					}
					// TODO: other security scheme types
				}

			}

			// add discovered variables to profiles
			for varName, varValue := range variables {
				for _, profile := range profiles {
					profile.Variables[varName] = varValue
				}
			}
			// add discovered query params to URL
			yamlRequest.Url = addQueryParamsToUrl(queryParams, yamlRequest.Url)

			requestBody := operationValue.RequestBody
			if requestBody != nil {
				// Take the first of many possible contents
				if requestBody.Content != nil {
					contentPairs := requestBody.Content.First()
					contentValue := contentPairs.Value()
					contentType := contentPairs.Key()
					example := generateMockExample(contentValue, contentType, jsonMockGenerator)

					if asByteArr, ok := example.([]byte); ok {
						yamlRequest.Body = string(asByteArr)
					} else if asMap, ok := example.(map[string]interface{}); ok {
						yamlRequest.Body = asMap
					}

					headers["Content-Type"] = []string{contentType}
				}
			}

			yamlRequest.Headers = yamlRequest.Headers.FromMap(headers)

			requestName := ""
			if operationValue.OperationId != "" {
				requestName = operationValue.OperationId
			} else {
				requestName = fmt.Sprintf("%s %s", yamlRequest.Method, pathPairs.Key())
			}

			filename := fmt.Sprintf("%s.yaml", sanitizeFileName(requestName))

			requestMold := model.RequestMold{
				Root:        workspace,
				Filename:    filename,
				ContentType: model.CONTENT_TYPE_YAML,
				Name:        requestName,
				Yaml:        &yamlRequest,
			}

			requests = append(requests, requestMold)
		}
	}

	return requests, profiles
}

func generateMockExample(contentValue *v3.MediaType, contentType string, jsonMockGenerator *renderer.MockGenerator) interface{} {
	example, err := jsonMockGenerator.GenerateMock(contentValue.Schema.Schema(), "")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate mock from schema")
		return []byte{}
	}
	if example == nil || string(example) == "<nil>" {
		return []byte{}
	}
	// special treatment to some content types
	switch contentType {
	case "application/xml":
		// there is no simple way to marshal data to xml/html
		// so at least for now, this is what you get
		return []byte("<xml></xml>")
	case "text/html":
		// there is no simple way to marshal data to xml/html
		// so at least for now, this is what you get
		return []byte("<html></html>")
	case "application/x-www-form-urlencoded", "multipart/form-data":
		return convertToMap(example)
	default:
		return example
	}

}

func convertToMap(example []byte) map[string]interface{} {
	var exampleMap map[string]interface{}
	err := json.Unmarshal(example, &exampleMap)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to umarshal %s to map", example)
	}
	if exampleMap == nil {
		exampleMap = make(map[string]interface{})
	}
	return exampleMap
}

func findSecurityScheme(securityName string, securitySchemes *orderedmap.Map[string, *v3.SecurityScheme]) *v3.SecurityScheme {
	for secSchemePair := securitySchemes.First(); secSchemePair != nil; secSchemePair = secSchemePair.Next() {
		if securityName == secSchemePair.Key() {
			return secSchemePair.Value()
		}
	}
	return nil
}

func handleOperationParametersV3(parameters []*v3.Parameter, mg *renderer.MockGenerator) (map[string][]string, map[string]string, map[string]string) {

	headers := map[string][]string{}
	queryParams := map[string]string{}
	variables := map[string]string{}

	// TODO: handle server override from path
	// operationValue.Servers
	// maybe in this case, just put value (if single) as operation url?
	for _, param := range parameters {
		example := []byte{}
		schema := param.Schema.Schema()
		if param.Example != nil {
			example = []byte(param.Example.Value)
		} else {
			mock, err := mg.GenerateMock(schema, "")
			if err != nil {
				log.Warn().Err(err).Msg("Failed to create mock")
				example = []byte{}
			} else if string(mock) != "<nil>" {
				example = mock
			}
		}
		if example == nil {
			example = []byte{}
		}
		// FIXME: what if another operation has a parameter with the same name? check if one already exists and scope it OR make the scope beforehand using the operationId?
		variables[param.Name] = string(example)
		templatedVar := fmt.Sprintf("{%s}", param.Name)
		if param.In == "header" {
			headers[param.Name] = []string{templatedVar}
		} else if param.In == "query" {
			queryParams[param.Name] = templatedVar
		} else if param.In == "cookie" {
			cookies, has := headers["Cookie"]
			if !has {
				cookies = []string{}
			}

			// NOTE: cookie header values differ from other headers with how they are delimited: they use semi-colon ";" instead of comma ","
			// In order to trick YAML marshalling not to add commas, we work with single array item and do our own delimitation
			cookies = append(cookies, fmt.Sprintf("%s=%s", param.Name, templatedVar))
			cookiesStr := strings.Join(cookies, ";")
			headers["Cookie"] = []string{cookiesStr}
		}
	}
	return headers, queryParams, variables
}

func handleServersV3(servers []*v3.Server, workspace string) []model.Profile {
	// TODO: environment, region
	profiles := []model.Profile{}
	for idx, server := range servers {
		name := fmt.Sprintf("profile-%d", idx+1)
		variables := map[string]string{}

		// TODO: if URL is relative, try to use spec location as base url (in case given path is a url)
		// otherwise, create BASE_URL variable and leave it empty
		// If there are no servers in spec file, create default profile
		variables["url"] = server.URL

		serverVariables := server.Variables
		for serverVarsPair := serverVariables.First(); serverVarsPair != nil; serverVarsPair = serverVarsPair.Next() {
			varName := serverVarsPair.Key()
			variable := serverVarsPair.Value()
			varValue := ""
			if variable != nil {
				if variable.Default != "" {
					varValue = variable.Default
				} else if len(variable.Enum) > 0 {
					varValue = variable.Enum[0]
				}
			}
			variables[varName] = varValue
		}

		profile := model.Profile{
			Name:      name,
			Root:      workspace,
			Filename:  fmt.Sprintf(".env.%s", name),
			Variables: variables,
		}
		profiles = append(profiles, profile)
	}

	if len(profiles) == 0 {
		name := "default"
		variables := map[string]string{}

		// TODO: if given path is an url, set it here
		// otherwise, create BASE_URL variable and leave it empty
		variables["url"] = ""
		profile := model.Profile{
			Name:      name,
			Root:      workspace,
			Filename:  ".env",
			Variables: variables,
		}
		profiles = append(profiles, profile)
	}
	return profiles
}

func addQueryParamsToUrl(queryParams map[string]string, baseUrl string) string {
	firstQueryP := true
	url := baseUrl
	for queryPName, queryPVal := range queryParams {
		queryKeyword := "&"
		if firstQueryP {
			queryKeyword = "?"
			firstQueryP = false
		}
		url = fmt.Sprintf("%s%s%s=%s", url, queryKeyword, queryPName, queryPVal)
	}
	return url
}

func loadSpec(path string) ([]byte, error) {

	if validator.IsValidUrl(path) {
		r := resty.New().R()
		resp, err := r.Get(path)
		if err != nil {
			return nil, err
		}
		return resp.Body(), nil
	} else {
		file, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

}

func sanitizeFileName(fileName string) string {
	// Define a regular expression to match invalid file name characters
	reg := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

	// Replace invalid characters with an underscore
	safeFileName := reg.ReplaceAllString(fileName, "_")

	// Additional replacement for Windows reserved names (optional)
	reservedNames := regexp.MustCompile(`^(CON|PRN|AUX|NUL|COM\d|LPT\d)(\..*)?$`)
	safeFileName = reservedNames.ReplaceAllString(safeFileName, "reserved_$1")

	return safeFileName
}
