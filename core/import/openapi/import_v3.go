package openapi

import (
	"fmt"
	"startpoint/core/model"
	"strings"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/renderer"
	"github.com/rs/zerolog/log"
)

func ImportOpenAPIV3(document libopenapi.Document, workspace string) ([]model.RequestMold, []model.Profile) {

	fmt.Print("Building given document as OpenAPI Spec v3...")
	v3Model, errors := document.BuildV3Model()

	if len(errors) > 0 {
		fmt.Print("FAILED\n")
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported",
			len(errors)))
	}
	fmt.Print("DONE\n")

	paths := v3Model.Model.Paths.PathItems
	schemas := v3Model.Model.Components.Schemas
	servers := v3Model.Model.Servers
	securitySchemes := v3Model.Model.Components.SecuritySchemes
	jsonMockGenerator := renderer.NewMockGenerator(renderer.JSON)

	fmt.Printf("\nThere are %d paths, %d schemas and %d servers in the document.\n", paths.Len(), schemas.Len(), len(servers))

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
					} else if securityScheme.Scheme == "bearer" {
						headers["Authorization"] = []string{fmt.Sprintf("Bearer {%s}", reqPairs.Key())}
						variables[reqPairs.Key()] = ""
					} else if securityScheme.Type == "apiKey" {
						if securityScheme.In == "cookie" {
							addCookieToHeaders(securityScheme.Name, fmt.Sprintf("{%s}", securityScheme.Name), headers)
							variables[securityScheme.Name] = ""
						} else if securityScheme.In == "header" {
							headers[securityScheme.Name] = []string{fmt.Sprintf("{%s}", securityScheme.Name)}
							variables[securityScheme.Name] = ""
						} else if securityScheme.In == "query" {
							queryParams[securityScheme.Name] = fmt.Sprintf("{%s}", securityScheme.Name)
							variables[securityScheme.Name] = ""
						}
					} else if securityScheme.Type == "openIdConnect" {
						// TODO: check if there is a way to generate prev_req
						headers["Authorization"] = []string{fmt.Sprintf("Bearer {%s}", reqPairs.Key())}
						variables[reqPairs.Key()] = ""
					} else if securityScheme.Type == "oauth2" {
						// TODO: check if there is a way to generate prev_req
						headers["Authorization"] = []string{fmt.Sprintf("Bearer {%s}", reqPairs.Key())}
						variables[reqPairs.Key()] = ""
					}
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
			addCookieToHeaders(param.Name, templatedVar, headers)
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
