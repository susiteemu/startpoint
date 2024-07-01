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
	"github.com/pb33f/libopenapi/utils"
	"github.com/rs/zerolog/log"
)

func ReadSpec(path string, workspace string) {

	specBytes, err := loadSpec(path)
	if err != nil {
		panic(err)
	}

	document, err := libopenapi.NewDocument(specBytes)

	if err != nil {
		panic(fmt.Sprintf("cannot create new document: %e", err))
	}

	var (
		requests []model.RequestMold
		profiles []model.Profile
	)
	if document.GetSpecInfo().SpecType == utils.OpenApi3 {
		requests, profiles = ImportOpenAPIV3(document, workspace)
	} else if document.GetSpecInfo().SpecType == utils.OpenApi2 {
		// TODO:
		fmt.Print("Importing from OpenAPI v2 is not supported\n")
		return
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

func addCookieToHeaders(cookieName string, cookieVal string, headers map[string][]string) {
	cookies, has := headers["Cookie"]
	if !has {
		cookies = []string{}
	}

	// NOTE: cookie header values differ from other headers with how they are delimited: they use semi-colon ";" instead of comma ","
	// In order to trick YAML marshalling not to add commas, we work with single array item and do our own delimitation
	cookies = append(cookies, fmt.Sprintf("%s=%s", cookieName, cookieVal))
	cookiesStr := strings.Join(cookies, ";")
	headers["Cookie"] = []string{cookiesStr}
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
		fmt.Print("Given location seems to be an URL. Proceeding to download the file...")
		r := resty.New().R()
		resp, err := r.Get(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to download from %s", path)
			fmt.Print("FAILED.\n")
			return nil, err
		}
		fmt.Print("DONE\n")
		return resp.Body(), nil
	} else {
		fmt.Print("Given location seems to be a local file. Proceeding to read it...")
		file, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read %s", path)
			fmt.Print("FAILED.\n")
			return nil, err
		}
		fmt.Print("DONE\n")
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
