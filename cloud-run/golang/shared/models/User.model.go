package myfunction

type User struct {
	FirestoreRecordBase

	Username string `json:"username" firestore:"username"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"password" firestore:"password"`
	PremiumUser bool   `json:"premiumUser" firestore:"premiumUser"`
	Links string[] `json:"links" firestore:"links"`

	FirestoreID string

	GetFirestoreID() string {
		return c.FirestoreID
	}

	SetFirestoreID(id string) {
		c.FirestoreID = id
	}
}