package myfunction

type User struct {
	Username string `json:"username" firestore:"username"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"-" firestore:"password"`
	PremiumUser bool   `json:"premiumUser" firestore:"premiumUser"`
	Links []string `json:"links" firestore:"links"`

	FirestoreID string `json:"-"`
}

func (u *User) SetFirestoreID(id string) {
    u.FirestoreID = id
}

func (u *User) GetFirestoreID() string {
	return u.FirestoreID
}