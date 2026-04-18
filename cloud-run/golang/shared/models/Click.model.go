package myfunction

type Click struct {
	FirestoreRecordBase

	ClickID   string `json:"clickID" firestore:"clickID"`
	IP        string `json:"ip" firestore:"ip"`
	LinkID    string `json:"linkID" firestore:"linkID"`
	Timestamp int64  `json:"timestamp" firestore:"timestamp"`
	UserAgent string `json:"userAgent" firestore:"userAgent"`
	OS        string `json:"os" firestore:"os"`
	Client    string `json:"client" firestore:"client"`
	Device    string `json:"device" firestore:"device"`
	Location  string `json:"location" firestore:"location"`
	ISP       string `json:"isp" firestore:"isp"`
	Mobile    bool   `json:"mobile" firestore:"mobile"`
	Proxy     bool   `json:"proxy" firestore:"proxy"`
	Hosting   bool   `json:"hosting" firestore:"hosting"`
	ASN       string `json:"asn" firestore:"asn"`

	FirestoreID string
}