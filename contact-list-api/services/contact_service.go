package services

import (
	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
	"github.com/mathzpereira/c214-seminario/contact-list-api/storage"
)

func GetAllContacts() ([]models.Contact, error) {
    return storage.LoadContacts()
}

func AddContact(newContact models.Contact) error {
    contacts, err := storage.LoadContacts()
    if err != nil {
        return err
    }

    newContact.ID = getNextID(contacts)
    contacts = append(contacts, newContact)
    return storage.SaveContacts(contacts)
}

func getNextID(contacts []models.Contact) int {
    maxID := 0
    for _, c := range contacts {
        if c.ID > maxID {
            maxID = c.ID
        }
    }
    return maxID + 1
}
