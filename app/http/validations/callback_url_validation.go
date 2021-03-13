package validations

import (
	"github.com/mvdan/xurls"
	"github.com/vandario/govms-ipg/app/helpers"
	"strings"
)

func CallbackUrlValidation(theUrlsField string, requestedCallbackUrl string) bool {
	prettyUrls := strings.Replace(theUrlsField, "\\", "", -1)
	urlsSlice := xurls.Relaxed().FindAllString(prettyUrls, -1)
	var theUrls []string

	for _, urls := range urlsSlice {
		removedTrueFromUrl := strings.Replace(urls, "\",true", "", -1)
		removedFalseFromUrl := strings.Replace(removedTrueFromUrl, "\",false", "", -1)
		theUrls = append(theUrls, removedFalseFromUrl)
	}

	var validatedUrls []int
	for _, singleUrl := range theUrls {
		if singleUrl == requestedCallbackUrl {
			validatedUrls = append(validatedUrls, 1)
		} else {
			validatedUrls = append(validatedUrls, 0)
		}
	}

	if helpers.SliceSummation(validatedUrls) == 0 {
		return false
	} else {
		return true
	}
}
