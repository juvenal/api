package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type RequestHeader struct {
	APISignature string
	APIVersion   int
	APIFormat    string
	CustomPrefix string
	CustomHeader map[string]string
}

func (rh *RequestHeader) APIInfoCollect(headers http.Request) (err error) {

	rh.CustomHeader = make(map[string]string)

	for header, values := range headers.Header {
		//
		switch  {
		case strings.ToLower(header) == "accept":
			//
			for _, value := range values {
				// Test for API version signature
				if strings.HasPrefix(value, rh.APISignature) {
					apiInfo := strings.Replace(value, rh.APISignature, "", 1)
					if matched, _ := regexp.MatchString(`v[0-9]+\+(json|xml)`, apiInfo); matched == true {
						splitApiInfo := strings.Split(apiInfo, "+")
						if rh.APIVersion, err = strconv.Atoi(strings.Replace(splitApiInfo[0], "v", "", 1)); err == nil {
							rh.APIFormat = splitApiInfo[1]
						} else {
							// Abort check and return error
							return errors.New("No valid API version informed... Aborting!")
						}
					} else {
						// No match for API information, abort and return error
						return errors.New("No API information provided... Aborting!")
					}
				} else {
					// No API version signature found, abort and return error
					return errors.New("No API Accept header provided... Aborting!")
				}
			}
			log.Printf("APIInfoCollect: Request - API version: %d, API format: %s\n",
				rh.APIVersion, rh.APIFormat)

		case strings.HasPrefix(strings.ToLower(header), rh.CustomPrefix):
			//
			for _, value := range values {
				rh.CustomHeader[header] = value
			}
			log.Printf("APIInfoCollect: Custom header - %s: %s\n",
				header, rh.CustomHeader[header])
		}
	}

	return nil
}

func (rh *RequestHeader) APIContentHeader() string {
	return fmt.Sprintf("%sv%d+%s", rh.APISignature, rh.APIVersion, rh.APIFormat)
}

func (rh *RequestHeader) APICustomHeaders() map[string]string {
	return rh.CustomHeader
}

func (rh *RequestHeader) APIAddCustomHeader(key string, value string) {
	rh.CustomHeader[key] = value
}

func (rh *RequestHeader) APIGetCustomHeader(key string) string {
	return rh.CustomHeader[key]
}
