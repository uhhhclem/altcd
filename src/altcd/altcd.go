package altcd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UINode is used to marshal an Altamont object to the JSON
// representation of a node in the UI's hierarchy.
type UINode struct {
	DisplayName string      `json:"displayName"`
	Data        interface{} `json:"data"`
	ChildUrl    string      `json:"childUrl"`
	TemplateUrl string      `json:"templateUrl"`
}

var serverData = map[string][]UINode{
	"/cd/": []UINode{
		{DisplayName: "packages", ChildUrl: "/cd/packages"},
		{DisplayName: "troubleshooters", ChildUrl: "/cd/troubleshooters"},
		{DisplayName: "others", ChildUrl: "/cd/others"},
	},
	"/cd/packages": []UINode{
		{
			DisplayName: "play provider",
			TemplateUrl: "package.ng",
		},
		{
			DisplayName: "repairs provider",
			TemplateUrl: "package.ng",
		},
	},
	"/cd/troubleshooters": []UINode{
		{
			DisplayName: "basic troubleshooters",
			ChildUrl:    "/cd/troubleshooters/basic",
		},
		{DisplayName: "authenticators"},
		{DisplayName: "fallback troubleshooters"},
	},
	"/cd/troubleshooters/basic": []UINode{
		{
			DisplayName: "where is the store app",
			TemplateUrl: "troubleshooter.ng",
		},
		{
			DisplayName: "rma returns",
			TemplateUrl: "troubleshooter.ng",
		},
	},
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	url := r.URL.Path
	data, ok := serverData[url]
	if !ok {
		err = fmt.Errorf("No data found for URL %q", url)
		return
	}
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Write(b)
}

func RunServer(staticDir string) {
	for url := range serverData {
		http.HandleFunc(url, urlHandler)
	}
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))
	http.ListenAndServe(":8080", nil)
}
