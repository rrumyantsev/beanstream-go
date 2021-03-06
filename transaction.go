package beanstream

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func ProcessBody(httpMethod string, url string, merchId string, apiKey string, data interface{}, responseType interface{}) (interface{}, error) {

	jsonData, _ := json.Marshal(data)
	//fmt.Println("--> Request: ", string(jsonData))
	//fmt.Println("Url: ", url)
	passcode := GenerateAuthCode(merchId, apiKey)
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
	//fmt.Println("Authorization: " + passcode)
	req.Header.Set("Authorization", "Passcode "+passcode)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("<-- Response:", string(body))
	//fmt.Println("response Status:", resp.Status)

	// handle errors
	if resp.StatusCode != 200 {
		return nil, handleError(resp, body)
	}

	err = json.Unmarshal([]byte(body), &responseType)
	if err != nil {
		return nil, &BeanstreamApiException{resp.StatusCode, 0, 0, err.Error(), "Error parsing Json response", nil}
	}

	//fmt.Printf("responseType: %T : %v\n", responseType, responseType)
	return responseType, nil
}

func ProcessMultiPart(httpMethod string, url string, merchId string, apiKey string, responseType interface{}, jsonCriteria interface{}, batchFile string) (interface{}, error) {

	jsonData, _ := json.Marshal(jsonCriteria)
	fmt.Println("--> Request: ", string(jsonData))
	fmt.Println("Url: ", url)
	passcode := GenerateAuthCode(merchId, apiKey)

	// multipart/form-data section:
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// add the json 'criteria'
	fw, err := w.CreateFormField("criteria")
	if err != nil {
		panic(err)
	}
	if _, err = fw.Write(jsonData); err != nil {
		panic(err)
	}

	// add the json type 'type'
	if fw, err = w.CreateFormField("type"); err != nil {
		panic(err)
	}
	if _, err = fw.Write([]byte("application/json")); err != nil {
		panic(err)
	}

	// add the batch file 'filename'
	f, err := os.Open(batchFile)
	if err != nil {
		return nil, err
	}
	if fw, err = w.CreateFormFile("filename", batchFile); err != nil {
		panic(err)
	}

	if _, err = io.Copy(fw, f); err != nil {
		panic(err)
	}

	err = w.Close() // close multipart writer
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(httpMethod, url, &b) // write the byte buffer
	req.Header.Set("Authorization", "Passcode "+passcode)
	//req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err) //return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("<-- Response:", string(body))
	//fmt.Println("response Status:", resp.Status)

	// handle errors
	if resp.StatusCode != 200 {
		return nil, handleError(resp, body)
	}

	err = json.Unmarshal([]byte(body), &responseType)
	if err != nil {
		return nil, &BeanstreamApiException{resp.StatusCode, 0, 0, err.Error(), "Error parsing Json response", nil}
	}

	//fmt.Printf("responseType: %T : %v\n", responseType, responseType)
	return responseType, nil
}

func handleError(resp *http.Response, body []byte) error {
	// parse json body
	ct := resp.Header.Get("Content-Type")
	fmt.Println("Error content-type: ", ct)
	if ct == "application/json; charset=utf-8" {
		errResp := errorResponse{}
		b := strings.Replace(string(body), "\"reference\":null,", "\"reference\":\"\",", -1)
		fmt.Println("Error message: ", b)
		err := json.Unmarshal([]byte(b), &errResp)
		if err != nil {
			return &BeanstreamApiException{resp.StatusCode, 0, 0, err.Error(), "Error parsing Json error message", nil}
		}
		return &BeanstreamApiException{resp.StatusCode, errResp.Code, errResp.Category, errResp.Message, errResp.Reference, errResp.Details}
	} else {
		return &BeanstreamApiException{resp.StatusCode, 0, 0, "", "Non-json error message. Content Type(" + ct + ")", nil}
	}
}

func Process(httpMethod string, url string, merchId string, apiKey string, responseType interface{}) (interface{}, error) {

	//fmt.Println("--> Request ")
	//fmt.Println("Url: ", url)
	passcode := GenerateAuthCode(merchId, apiKey)
	req, err := http.NewRequest(httpMethod, url, nil)
	req.Header.Set("Authorization", "Passcode "+passcode)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("<-- Response:", string(body))
	//fmt.Println("response Status:", resp.Status)

	// handle errors
	if resp.StatusCode != 200 {
		return nil, handleError(resp, body)
	}

	err = json.Unmarshal([]byte(body), &responseType)
	if err != nil {
		return nil, &BeanstreamApiException{resp.StatusCode, 0, 0, err.Error(), "Error parsing Json response", nil}
	}

	//fmt.Printf("responseType: %T : %v\n", responseType, responseType)
	return responseType, nil
}

func GenerateAuthCode(merchId string, apiKey string) string {
	return base64.StdEncoding.EncodeToString([]byte(string(merchId + ":" + apiKey)))
}

// JSON:
// {
//	"code":200,
//	"category":2,
//	"message":"Transaction cannot be adjusted",
//	"reference":nil
//	"details":[{"field":"card_name","message":"Card owner name is missing"}]
// }
type errorResponse struct {
	Code      int           `json:"code,omitempty"`
	Category  int           `json:"category,omitempty"`
	Message   string        `json:"message,omitempty"`
	Reference string        `json:"reference,omitempty"`
	Details   []ErrorDetail `json:"details,omitempty"`
}
