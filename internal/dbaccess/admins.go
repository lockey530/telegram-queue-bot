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

// Kick a person with a Telegram handle.
// Returns a chatID for further communications.
func RemoveFirstInQueue(toRemove string) (chatID int64, err error) {
	// this will need to be adjusted after implementing the waiting feature.
	var chat []int64
	err = db.Select(&chat, `
		DELETE FROM queue 
		WHERE user_handle = $1
		RETURNING chat_id;
	`, toRemove)
	if err != nil {
		return 0, fmt.Errorf("failed to kick @%v. %v", toRemove, err)
	}
	chatID = chat[0]

	return chatID, nil
}

// Kick a person with a Telegram handle.
// Returns a chatID for further communications.
func KickPerson(toKick string) (chatID int64, err error) {
	// this will need to be adjusted after implementing the waiting feature.
	var chat []int64
	err = db.Select(&chat, `
		DELETE FROM queue 
		WHERE user_handle = $1
		RETURNING chat_id;
	`, toKick[1:])
	if err != nil {
		return 0, fmt.Errorf("failed to kick @%v. %v", toKick, err)
	}
	chatID = chat[0]

	return chatID, nil
}

func SeeAdminList(retriever string) ([]string, error) {
	var adminList []string

	err := db.Select(&adminList, `
		SELECT 
			(admin_handle)
		FROM
			admins;
	`)

	log.Println(adminList)

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
			($1);
	`, toAdd)

	if err != nil {
		log.Printf("failed to add @%s by @%s\n", toAdd, requestor)
		return err
	} else {
		log.Printf("admin @%s successfully added by @%s\n", toAdd, requestor)
		return nil
	}
}

func RemoveAdmin(toRemove string, requestor string) (string, error) {
	res, err := db.Exec(`
		DELETE FROM
			admins
		WHERE
			admin_handle = $1 AND removable = true;
	`, toRemove)

	if err != nil {
		log.Printf("failed to remove @%s by @%s\n", toRemove, requestor)
		return "", err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("failed to remove @%s by @%s. Error: %s\n", toRemove, requestor, err.Error())
		return "", err
	} else if affected == 0 {
		log.Printf("failed to remove @%s by @%s as @%s is either a protected admin, or does not exist.\n", toRemove, requestor, toRemove)

		issue := fmt.Sprintf("could not remove @%s because @%s does not exist or cannot be removed.\n", toRemove, toRemove)
		return issue, fmt.Errorf(issue)
	} else {
		log.Printf("admin @%s successfully removed by @%s\n", toRemove, requestor)
		return "", nil
	}
}
