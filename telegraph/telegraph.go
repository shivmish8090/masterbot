package telegraph

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

// Telegraph API Base URL
const TelegraphApi = "https://api.telegra.ph/"

// AccountMap stores access tokens with their flood-wait expiration times (default: 0, means available)
var AccountMap map[string]int64

// TelegraphAccountResponse represents the response from the Telegraph createAccount API
type TelegraphAccountResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		AccessToken string `json:"access_token"`
	} `json:"result"`
}

// TelegraphResponse represents the response from the Telegraph createPage API
type TelegraphResponse struct {
	OK     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result struct {
		URL string `json:"url"`
	} `json:"result"`
}

// createAccount creates a new Telegraph account and returns an access token
func createAccount(shortName string) (string, error) {
	payload := map[string]string{
		"short_name":  shortName,
		"author_name": "Anonymous",
		"author_url":  "https://t.me/ViyomBot",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(TelegraphApi + "createAccount")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(jsonData)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return "", err
	}

	var result TelegraphAccountResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	if !result.OK {
		return "", errors.New("failed to create account")
	}

	return result.Result.AccessToken, nil
}

// init initializes 10 Telegraph accounts
func init() {
	AccountMap = make(map[string]int64)
	for i := 0; i < 10; i++ {
		token, err := createAccount("EcoBot" + strconv.Itoa(i+1))
		if err == nil {
			AccountMap[token] = 0 // No flood wait initially
		}
	}
}

// getAvailableToken returns an access token that is not under a flood wait
func getAvailableToken() (string, error) {
	now := time.Now().Unix()
	for token, waitTime := range AccountMap {
		if waitTime == 0 || waitTime <= now {
			return token, nil
		}
	}
	return "", errors.New("no available accounts due to flood wait")
}

// extractFloodWait extracts flood wait time from error message
func extractFloodWait(errorMsg string) int64 {
	re := regexp.MustCompile(`FLOOD_WAIT_(\d+)`)
	matches := re.FindStringSubmatch(errorMsg)
	if len(matches) > 1 {
		seconds, err := strconv.ParseInt(matches[1], 10, 64)
		if err == nil {
			return seconds
		}
	}
	return 0
}

// CreatePage creates a Telegraph page while handling flood wait errors
func CreatePage(content, firstName string) (string, error) {
	for {
		accessToken, err := getAvailableToken()
		if err != nil {
			return "", err
		}

		// Convert content to Telegraph node format
		telegraphContent := []map[string]interface{}{
			{"tag": "p", "children": []string{content}},
		}

		// Prepare request payload
		payload := map[string]interface{}{
			"access_token": accessToken,
			"title":        "Eco Message",
			"author_name":  firstName,
			"author_url":   "https://t.me/ViyomBot",
			"content":      telegraphContent,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.SetRequestURI(TelegraphApi + "createPage")
		req.Header.SetMethod("POST")
		req.Header.SetContentType("application/json")
		req.SetBody(jsonData)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		client := fasthttp.Client{}
		if err := client.Do(req, resp); err != nil {
			return "", err
		}

		var result TelegraphResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return "", err
		}

		if result.OK {
			return result.Result.URL, nil
		}

		// Check for flood wait error
		floodWaitTime := extractFloodWait(result.Error)
		if floodWaitTime > 0 {
			AccountMap[accessToken] = time.Now().Unix() + floodWaitTime // Store expiration time
			continue                                                    // Try another account
		}

		// Other errors, return immediately
		return "", errors.New(result.Error)
	}
}
