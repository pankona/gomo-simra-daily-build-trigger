package trigger

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func getEnvVar(varName string) (result string) {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == varName {
			return pair[1]
		}
	}
	return ""
}

func createHttpClient(r *http.Request) *http.Client {
	return urlfetch.Client(appengine.NewContext(r))
}

func init() {
	http.HandleFunc("/trigger", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	circleCIBuild(w, r)
	dockerBuild(w, r)
}

func circleCIBuild(w http.ResponseWriter, r *http.Request) {
	circleci_url := "https://circleci.com/api/v1/project/pankona/gomo-simra/tree/master"
	query := url.Values{"circle-token": {getEnvVar("CIRCLECI_API_KEY")}}

	req, _ := http.NewRequest(
		"POST",
		circleci_url+"?"+query.Encode(),
		nil,
	)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := createHttpClient(r)
	resp, err := client.Do(req)

	if err == nil {
		fmt.Fprint(w, "err == nil")
		fmt.Fprint(w, resp)
	} else {
		fmt.Fprint(w, "err != nil")
		fmt.Fprint(w, err)
	}
}

func dockerBuild(w http.ResponseWriter, r *http.Request) {
	//curl -H "Content-Type: application/json" --data '{"docker_tag": "master"}' -X POST https://registry.hub.docker.com/u/pankona/gomo-simra/trigger/8e5d67eb-c9b3-4b1e-bde1-1cbaad06d47b/
}
