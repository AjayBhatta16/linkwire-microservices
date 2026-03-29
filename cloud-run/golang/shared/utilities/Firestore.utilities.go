package myfunction

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

func GetFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	projectID := os.Getenv("GCP_PROJECT_ID")

	client, err := Firestore.NewClient(ctx, projectID)
	
	return client, err
}

func GetItemsByFieldValue[TData FirestoreRecordBase](
	collectionName string, 
	fieldName string, 
	fieldValue string) ([]TData, error) 
{
	ctx := context.Background()

	client, err := GetFirestoreClient(ctx)

	if err != nil {
		log.PrintLn("GetItemsByFieldValue - Error creating Firestore client:", err)
		return nil, err
	}

	defer client.Close()
	
	var items []TData

	iter := client.Collection(collectionName).Where(fieldName, "==", fieldValue).Documents(ctx)
	defer iter.Stop()

	docs, err := iter.GetAll()

	if err != nil {
		log.PrintLn("GetItemsByFieldValue - Error fetching documents:", err)
		return nil, err
	}

	for _, doc := range docs {
		var data TData
		err := doc.DataTo(&data)
		if err != nil {
			log.PrintLn("GetItemsByFieldValue - Error converting document data:", err)
			continue
		}
		data.SetFirestoreID(doc.Ref.ID)
		items = append(items, data)
	}

	return items, nil
}

func CreateItem[TData FirestoreRecordBase](
	collectionName string, 
	data TData) error 
{
	ctx := context.Background()

	client, err := GetFirestoreClient(ctx)

	if err != nil {
		log.PrintLn("CreateItem - Error creating Firestore client:", err)
		return err
	}

	defer client.Close()

	_, _, err := client.Collection(collectionName).Add(ctx, data)

	if err != nil {
		log.PrintLn("CreateItem - Error adding document:", err)
		return err
	}

	return nil
}

func UpdateItem[TData FirestoreRecordBase](
	collectionName string, 
	firestoreID string,
	data TData) error
{
	ctx := context.Background()

	client, err := GetFirestoreClient(ctx)

	if err != nil {
		log.PrintLn("UpdateItem - Error creating Firestore client:", err)
		return err
	}

	defer client.Close()

	_, err := client.Collection(collectionName).Doc(firestoreID).Set(ctx, data)

	if err != nil {
		log.PrintLn("UpdateItem - Error updating document:", err)
		return err
	}

	return
}