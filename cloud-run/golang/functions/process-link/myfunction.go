package myfunction

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
    "github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	subscriber, err := utilities.NewSubscriber(context.Background())

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating subscriber: %v", err), http.StatusInternalServerError)
		return
	}

	defer subscriber.Close()

	req, _, err := utilities.Receive[Request](subscriber, r)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error receiving message: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received message - linkID: %s", req.LinkID)

	data, err2 := utilities.GetItemsByFieldValue[models.Link, *models.Link](constants.LINK_CONTAINER_NAME, "trackingID", req.LinkID)

	if err2 != nil {
		http.Error(w, fmt.Sprintf("Error fetching link data: %v", err2), http.StatusInternalServerError)
		return
	}

	link := data[0]

	meta, err3 := FetchPageMeta(link.RedirectURL)

	if err3 != nil {
		http.Error(w, fmt.Sprintf("Error fetching page metadata: %v", err3), http.StatusInternalServerError)
		return
	}

	link.SiteTitle = meta.Title
	link.SiteDescription = meta.Description
	link.SiteBannerURL = meta.OGImage

	err4 := utilities.UpdateItem[models.Link](constants.LINK_CONTAINER_NAME, link.FirestoreID, link)

	if err4 != nil {
		http.Error(w, fmt.Sprintf("Error updating link data: %v", err4), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func FetchPageMeta(url string) (*PageMeta, error) {
	client := &http.Client{Timeout: 10 * time.Second}
 
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; pagemeta-bot/1.0)")
 
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error fetching URL: %w", err)
	}

	defer resp.Body.Close()
 
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}
 
	doc, err := html.Parse(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error parsing HTML: %w", err)
	}
 
	meta := &PageMeta{}

	extractMeta(doc, meta)

	return meta, nil
}

func extractMeta(n *html.Node, meta *PageMeta) {
	if n.Type == html.ElementNode {
		switch strings.ToLower(n.Data) {
 
			case "title":
				if meta.Title == "" && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					meta.Title = strings.TrimSpace(n.FirstChild.Data)
				}
	
			case "meta":
				name := attrVal(n, "name")
				property := attrVal(n, "property")
				content := attrVal(n, "content")
	
				switch {
					case strings.EqualFold(name, "description") && meta.Description == "":
						meta.Description = content
					case strings.EqualFold(property, "og:description") && meta.Description == "":
						meta.Description = content
					case strings.EqualFold(property, "og:image") && meta.OGImage == "":
						meta.OGImage = content
					case strings.EqualFold(property, "og:title") && meta.OGTitle == "":
						meta.OGTitle = content
					case strings.EqualFold(name, "twitter:image") && meta.OGImage == "":
						meta.OGImage = content
					case strings.EqualFold(name, "twitter:title") && meta.OGTitle == "":
						meta.OGTitle = content
					case strings.EqualFold(name, "twitter:description") && meta.Description == "":
						meta.Description = content
				}
		}
	}
 
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		extractMeta(child, meta)
	}
}

func attrVal(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, name) {
			return a.Val
		}
	}
	return ""
}