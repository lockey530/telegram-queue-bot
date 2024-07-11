package dbaccess

import (
	"fmt"
	"log"
)

func CheckIfAdmin(handle string) (bool, error) {
	var count int

	err := db.Get(&count, `
		SELECT 
			COUNT(*)
		FROM
			admins
		WHERE
			admin_handle = $1;
	`, handle)

	if err != nil {
		log.Println("failed to retrieve admin list.")
		return false, err
	}
	return count > 0, nil
}

func SeeAdminList(retriever string) ([]string, error) {
	var adminList []string

	err := db.Select(&adminList, `
		SELECT 
			(admin_handle)
		FROM
			admins;
	`)

	if err != nil {
		log.Println("failed to retrieve admin list.")
		return nil, err
	} else {
		log.Printf("admin list retrieved by %s\n", retriever)
		return adminList, nil
	}
}

func AddAdmin(toAdd string, requestor string) error {
	_, err := db.Exec(`
		INSERT INTO admins
			(admin_handle)
		VALUES
			$1;
	`, toAdd)

	if err != nil {
		log.Printf("failed to add @%s by @%s\n", toAdd, requestor)
		return err
	} else {
		log.Printf("admin @%s successfully added by @%s\n", toAdd, requestor)
		return nil
	}
}

// Kick a person with a Telegram handle.
// Returns a chatID for further communications.
func KickPerson(telegramHandle string) (chatID int64, err error) {
	// this will need to be adjusted after implementing the waiting feature.
	var chat []int64
	err = db.Select(&chat, `
		DELETE FROM queue 
		WHERE user_handle = $1
		RETURNING chat_id;
	`, telegramHandle[1:])
	if err != nil {
		return 0, fmt.Errorf("failed to kick @%v. %v", telegramHandle, err)
	}
	chatID = chat[0]

	return chatID, nil
}
