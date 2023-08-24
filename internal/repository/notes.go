package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	database "github.com/devigor/donna-markdown-service/internal/config/db"
	"github.com/devigor/donna-markdown-service/internal/contracts"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Select() ([]*contracts.NotesBody, error) {
	db, err := database.OpenConn()
	if err != nil {
		log.Fatalln("Error to connect database\n%r", err)
	}

	rows, error := db.Query(context.Background(),
		"SELECT id, content, created_at, updated_at FROM donna_notes")
	defer db.Close(context.Background())
	defer rows.Close()

	if error != nil {
		log.Fatalln(error)
	}

	var results []*contracts.NotesBody

	for rows.Next() {
		var noteStruct contracts.NotesBody
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&noteStruct.Id, &noteStruct.Content, &createdAt, &updatedAt); err != nil {
			log.Fatalln(err)
		}

		noteStruct.CreatedAt = timestamppb.New(createdAt) // Converta para timestamppb.Timestamp
		noteStruct.UpdatedAt = timestamppb.New(updatedAt)

		results = append(results, &noteStruct)
		fmt.Println(results)

	}

	return results, nil
}

func Create(content string) error {
	data := &contracts.NotesBody{
		Id:      uuid.New().String(),
		Content: content,
	}

	db, err := database.OpenConn()
	if err != nil {
		log.Fatalln("Error to connect database\n%r", err)
		return err
	}

	_, error := db.Exec(context.Background(),
		"INSERT INTO donna_notes (id, content) VALUES ($1, $2)", data.Id, data.Content)

	defer db.Close(context.Background())
	return error
}
