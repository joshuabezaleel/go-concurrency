package tickersingleclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type Response struct {
	Confirmed Confirmed `json:"confirmed"`
}

type Confirmed struct {
	Value int `json:"value"`
}

var allow throttleFunc

func main() {
	allow = throttle(10) // allow up to 10 requests / second

	http.HandleFunc("/confirmed", getConfirmedCOVID)
	log.Fatal(http.ListenAndServe(":8093", nil))
}

func getConfirmedCOVID(w http.ResponseWriter, r *http.Request) {
	if !allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many request"))
		return
	}

	httpClient := http.Client{}
	request, err := http.NewRequest("GET", "https://covid19.mathdro.id/api", nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// throttling section here

type throttleFunc func() bool

// throttle limits the request to at most N requests / second.
// this function returns a func() bool which indicates
// subsequent request can be safely made or not.
func throttle(n int) throttleFunc {
	var counter uint32
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			atomic.StoreUint32(&counter, 0)
		}
	}()
	return func() bool {
		c := atomic.AddUint32(&counter, 1)
		return int(c) <= n
	}
}
