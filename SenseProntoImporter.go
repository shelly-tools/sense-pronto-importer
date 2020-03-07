/**
 * @package   SenseProntoImporter
 * @copyright Thorsten Eurich
 * @license   GNU Affero General Public License (https://www.gnu.org/licenses/agpl-3.0.de.html)
 *
 * @todo lots of documentation
 *
 * Simple tool to replace IR pronto codes from a Shelly Sense.
 */

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func toJSON(m interface{}) string {
	js, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return strings.ReplaceAll(string(js), ",", ", ")
}

func listIR(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")

	sense := "http://" + ip + "/ur/list"
	req, _ := http.NewRequest("GET", sense, nil)
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Fprintf(w, string(body))
}

func addIR(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	pronto := r.URL.Query().Get("pronto")

	sense := "http://" + ip + "/ir/add?id=" + id + "&name=" + url.QueryEscape(name) + "&pronto=" + pronto
	req, _ := http.NewRequest("GET", sense, nil)

	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Fprintf(w, string(body))
}

func tplIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}
	items := struct {
		Shelly string
	}{
		Shelly: "IP",
	}
	t.Execute(w, items)
}

func main() {
	http.HandleFunc("/", tplIndex)
	http.HandleFunc("/list", listIR)
	http.HandleFunc("/add", addIR)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
