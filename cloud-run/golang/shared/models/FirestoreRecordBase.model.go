package myfunction

type FirestoreRecordBase interface {
	GetFirestoreID() string
	SetFirestoreID(string)
}