package dbaccess

import (
	"fmt"
	"log"

	"github.com/josh1248/nusc-queue-bot/internal/types"
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

func RemoveFirstInQueue() (userHandle int64, err error) {
	var chat []struct {
		ChatID     int64  `db:"chat_id"`
		UserHandle string `db:"user_handle"`
	}
	err = db.Select(&chat, `
		DELETE FROM queue 
		WHERE user_handle = (
			SELECT 
				user_handle
			FROM 
				queue
			ORDER BY 
				joined_at
			LIMIT 1
		)
		RETURNING chat_id, user_handle;
	`)
	if err != nil {
		return -1, fmt.Errorf("failed to remove people from queue. Error: %v", err)
	} else if len(chat) == 0 {
		return -1, fmt.Errorf("no people present in the queue")
	}

	handleRemoved := chat[0].UserHandle
	log.Printf("successfully removed first person (%s) in queue.\n", handleRemoved)

	return chat[0].ChatID, nil
}

// position should be 1-indexed.
func GetPositionInQueue(position int) (userHandle string, chatID int64, err error) {
	var chat []types.QueueUser

	err = db.Select(&chat, `
		SELECT 
			*
		FROM 
			queue 
		ORDER BY 
			joined_at;
	`)

	if err != nil {
		return "", -1, err
	} else if len(chat) < position {
		return "", -1, fmt.Errorf("your queue only has %v people\n", len(chat))
	}

	return chat[position-1].UserHandle, chat[position-1].ChatID, nil
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
	`, toKick)
	if err != nil {
		return 0, fmt.Errorf("failed to kick @%v. %v", toKick, err)
	}

	if len(chat) != 1 {
		return -1, fmt.Errorf("@%s is not in the queue", toKick)
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
