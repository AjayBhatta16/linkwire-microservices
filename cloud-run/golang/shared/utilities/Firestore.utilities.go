package myfunction

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func GetFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	
	return client, err
}

func GetItemsByFieldValue[TData FirestoreRecordBase](
	collectionName string, 
	fieldName string, 
	fieldValue string) ([]TData, error) {
		
	ctx := context.Background()

	client, err := GetFirestoreClient(ctx)

	if err != nil {
		log.Println("GetItemsByFieldValue - Error creating Firestore client:", err)
		return nil, err
	}

	defer client.Close()
	
	var items []TData

	log.Printf("Querying collection '%s' for documents where '%s' == '%s'", collectionName, fieldName, fieldValue)

	iter := client.Collection(collectionName).Where(fieldName, "==", fieldValue).Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()

	if err != nil {
		log.Println("GetItemsByFieldValue - Error fetching documents:", err)
		return nil, err
	}

	for _, doc := range docs {
		var data TData
		err := doc.DataTo(&data)
		if err != nil {
			log.Println("GetItemsByFieldValue - Error converting document data:", err)
			continue
		}
		data.SetFirestoreID(doc.Ref.ID)
		items = append(items, data)
	}

	return items, nil
}

func CreateItem[TData FirestoreRecordBase](
	collectionName string, 
	data TData) error {

	ctx := context.Background()

	client, err := GetFirestoreClient(ctx)

	if err != nil {
		log.Println("CreateItem - Error creating Firestore client:", err)
		return err
	}

	defer client.Close()

	_, _, err2 := client.Collection(collectionName).Add(ctx, data)

	if err2 != nil {
		log.Println("CreateItem - Error adding document:", err2)
		return err2
	}

	return nil
}

func UpdateItem[TData FirestoreRecordBase](
	collectionName string, 
	firestoreID string,
	data TData) error {

	ctx := context.Background()

	client, err := GetFirestoreClient(ctx)

	if err != nil {
		log.Println("UpdateItem - Error creating Firestore client:", err)
		return err
	}

	defer client.Close()

	_, err2 := client.Collection(collectionName).Doc(firestoreID).Set(ctx, data)

	if err2 != nil {
		log.Println("UpdateItem - Error updating document:", err2)
		return err2
	}

	return nil
}