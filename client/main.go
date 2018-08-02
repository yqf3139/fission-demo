package main

import (
	"time"
	"math/rand"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"os"
	"net/http"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"io"
	"io/ioutil"
	"runtime"
)

const (
	//SERVER_URL = "http://cluster.me:31314"
	SERVER_URL = "http://router.fission"
	SECRET     = "secret"
)

func getToken() string {
	rand.Seed(time.Now().UnixNano())
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   fmt.Sprintf("%v", rand.Intn(10000)),
		"admin": false,
	})
	token, _ := t.SignedString([]byte(SECRET))
	return token
}

func main() {
	token := getToken()
	fmt.Println("Using token: ", token)

	mode := "normal"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	switch mode {
	case "normal":
		// report every 1 min +- 20 sec
		go report(40, 20, token)
		// upload every 5 +- 2 min
		go upload(4*60, 3*60, token)
	case "sleepy":
		// report every 20 +- 5 min
		go report(20*60, 5*60, token)
		// do not upload
	case "extreme":
		// report every 2 second
		go report(2, 1, token)
		// upload every 20 second
		go upload(20, 2, token)
	}
	select {}
}

func report(interval, fluctuation int32, token string) {
	url := fmt.Sprintf("%v/api/client/status", SERVER_URL)
	for {
		status := "working"
		if rand.Float32() > 0.9 {
			status = "error"
		}
		body, _ := json.Marshal(map[string]interface{}{
			"status":  status,
			"battery": rand.Intn(100),
		})
		err := doPut(url, token, body)
		fmt.Println("report done", err)

		t := interval + (rand.Int31n(fluctuation)*2 - fluctuation)
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func copy(src string, dst string) bool {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func upload(interval, fluctuation int32, token string) {
	url := fmt.Sprintf("%v/api/images", SERVER_URL)
	for {
		t := interval + (rand.Int31n(fluctuation)*2 - fluctuation)
		time.Sleep(time.Duration(t) * time.Second)

		files, err := ioutil.ReadDir("/images")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(files) == 0 {
			fmt.Println("images folder is empty")
			continue
		}

		filename := files[rand.Intn(len(files))].Name()
		oldpath := "/images/" + filename
		newpath := fmt.Sprintf("./%v", fmt.Sprintf("%v%v", time.Now().Unix(), filename))
		if copy(oldpath, newpath) {
			os.Remove(oldpath)
			err = doUpload(url, token, newpath)
			fmt.Println("upload done", err)
		}
		os.Remove(newpath)
		runtime.GC()
	}
}

func doPut(url, token string, body []byte) (err error) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", response.Status)
	}
	return
}

func doUpload(url, token, file string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("image", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}
