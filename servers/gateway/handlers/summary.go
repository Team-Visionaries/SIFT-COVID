package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"
const headerContentType = "Content-Type"
const jsonContent = "application/json"
const htmlContent = "text/html"

// properties that map structs
const ogPrefix = "og:"
const imageProp = "image"
const titleProp = "title"
const urlProp = "url"
const siteNameProp = "site_name"
const descrProp = "description"
const authorProp = "author"
const typeProp = "type"
const keywordsProp = "keywords"

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	w.Header().Add(headerCORS, corsAnyOrigin)
	w.Header().Add(headerContentType, jsonContent)
	url := r.URL.Query().Get(urlProp)
	if len(url) == 0 {
		log.Println("URL query parameter is missing")
		http.Error(w, "The url query parameter is missing", http.StatusBadRequest)
		return
	}
	respBody, fetchErr := fetchHTML(url)
	if fetchErr != nil {
		log.Println("Error in fetching html \n", fetchErr)
		http.Error(w, fetchErr.Error(), http.StatusInternalServerError)
		return
	}
	pageSum, summaryErr := extractSummary(url, respBody)
	if summaryErr != nil {
		log.Println("Error creating summary", summaryErr)
		http.Error(w, summaryErr.Error(), http.StatusInternalServerError)
		return
	}
	respBody.Close()
	buffer, err := json.Marshal(pageSum)
	if err != nil {
		log.Println("Error creating json", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buffer)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	resp, err := http.Get(pageURL)

	//if there was an error, report it and exit
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(http.StatusText(resp.StatusCode))
	}
	ctype := resp.Header.Get(headerContentType)
	if !strings.HasPrefix(ctype, htmlContent) {
		return nil, fmt.Errorf("response content type was not text/html")
	}
	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/
	actualURL := pageURL[:strings.Index(pageURL, "com")+3]

	summaryMap, error := getData(htmlStream)
	if error != nil {
		return nil, error
	}

	// populate structs
	ps := &PageSummary{}
	ps.Type = summaryMap[typeProp]
	ps.URL = summaryMap[urlProp]
	if summaryMap[titleProp] == "" {
		ps.Title = summaryMap["main_title"]
	} else {
		ps.Title = summaryMap[titleProp]
	}
	ps.SiteName = summaryMap[siteNameProp]
	ps.Description = summaryMap[descrProp]
	ps.Author = summaryMap[authorProp]

	// get all images
	ps.Images = getImages(summaryMap, actualURL)
	if summaryMap[keywordsProp] != "" {
		ps.Keywords = strings.Split(strings.ReplaceAll(summaryMap[keywordsProp], " ", ""), ",")
	}
	// create and populate icon
	if summaryMap["link:rel"] != "" {
		icon := &PreviewImage{}
		url := summaryMap["link:href"]
		if strings.HasPrefix(url, "/") {
			icon.URL = actualURL + url
		} else {
			icon.URL = url
		}
		icon.Type = summaryMap["link:type"]
		// handle string sizes to extract width and height
		dim := summaryMap["link:sizes"]
		if dim != "" && strings.Contains(dim, "x") {
			i := strings.Index(dim, "x")
			w, we := strconv.Atoi(dim[i+1:])
			h, he := strconv.Atoi(dim[:i])
			if we == nil {
				icon.Width = w
			}
			if he == nil {
				icon.Height = h
			}
		}
		ps.Icon = icon
	}
	return ps, nil
}

// helper function that goes through the file and extracts all necessary
// fields into map for post processing
func getData(htmlStream io.ReadCloser) (map[string]string, error) {
	tokenizer := html.NewTokenizer(htmlStream)
	summaryList := make(map[string]string)
	imageCount := 0
	for {
		tokenType := tokenizer.Next()

		//if it's an error token, we either reached
		//the end of the file, or the HTML was malformed
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			err = tokenizer.Err()
			return nil, err
		}
		// if the head tag has closed we are done populating the
		// struct and done reading the file
		if tokenType == html.EndTagToken && tokenizer.Token().Data == "head" {
			break
		}
		// process the token according to the token type
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			// check if it is of "meta property type"
			if "meta" == token.Data {
				tokenizer.Next()
				attributesList := token.Attr
				attrName := ""
				attrVal := ""
				// add attributes in list
				for _, attr := range attributesList {
					if "property" == attr.Key && strings.HasPrefix(attr.Val, ogPrefix) {
						attrName = strings.TrimLeft(attr.Val, ogPrefix)
					}
					if "name" == attr.Key {
						attrName = attr.Val
					}
					if "content" == attr.Key {
						attrVal = attr.Val
					}

				}
				if len(attrName) != 0 && len(attrVal) != 0 {
					// adding image data depending on number of images
					if imageProp == attrName {
						imageCount++
					}
					if strings.Contains(attrName, imageProp) {
						attrName = attrName + strconv.Itoa(imageCount)
					}
					// add to map
					if summaryList[attrName] == "" {
						summaryList[attrName] = attrVal
					}
				}
			}
			// in case there is no og:title then get main title
			if titleProp == token.Data {
				titleType := tokenizer.Next()
				if titleType == html.TextToken {
					summaryList["main_title"] = tokenizer.Token().Data
				}
			}
			// add link attributes with "link" identifier
			// and is of type icon
			if "link" == token.Data {
				tokenizer.Next()
				attributesList := token.Attr
				linkAttrs := make(map[string]string)
				for _, attr := range attributesList {
					linkAttrs[attr.Key] = attr.Val
				}
				if linkAttrs["rel"] == "icon" {
					for key, element := range linkAttrs {
						summaryList["link:"+key] = element
					}
				}
			}
		}
		// for loop end
	}
	// keep track of number of images for post processing
	summaryList["image_count"] = strconv.Itoa(imageCount)
	return summaryList, nil
}

// helper function that creates the slice of PreviewImage that contains
// all the images in the html file
func getImages(summaryMap map[string]string, pageURL string) []*PreviewImage {
	imgList := []*PreviewImage{}
	imgCount, countError := strconv.Atoi(summaryMap["image_count"])
	if countError == nil && imgCount != 0 {
		for i := 1; i <= imgCount; i++ {
			k := strconv.Itoa(i)
			pv := &PreviewImage{}
			url := summaryMap["image"+k]
			if strings.HasPrefix(url, "/") {
				pv.URL = pageURL + url
			} else {
				pv.URL = url
			}
			pv.SecureURL = summaryMap["image:secure_url"+k]
			pv.Type = summaryMap["image:type"+k]
			width, convErr := strconv.Atoi(summaryMap["image:width"+k])
			if convErr == nil {
				pv.Width = width
			}
			height, convErr := strconv.Atoi(summaryMap["image:height"+k])
			if convErr == nil {
				pv.Height = height
			}
			pv.Alt = summaryMap["image:alt"+k]
			imgList = append(imgList, pv)
		}
		return imgList
	}
	return nil
}
