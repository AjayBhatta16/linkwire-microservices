package myfunction

type Link struct {
	FirestoreRecordBase

	TrackingID      string  `json:"trackingID" firestore:"trackingID"`
	DisplayID       string  `json:"displayID" firestore:"displayID"`
	RedirectURL     string  `json:"redirectURL" firestore:"redirectURL"`
	Note            string  `json:"note" firestore:"note"`
	SiteTitle       string  `json:"siteTitle" firestore:"siteTitle"`
	SiteDescription string  `json:"siteDescription" firestore:"siteDescription"`
	SiteBannerURL   string  `json:"siteBannerURL" firestore:"siteBannerURL"`
	UseLogin        bool    `json:"useLogin" firestore:"useLogin"`
	LoginPageBrand  string  `json:"loginPageBrand" firestore:"loginPageBrand"`
	CreatedBy       string  `json:"createdBy" firestore:"createdBy"`
	Clicks          []Click `json:"clicks" firestore:"clicks"`

	FirestoreID string `json:"-"`
}