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

func GetContactByID(id int) (models.Contact, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return models.Contact{}, err
	}

	for _, contact := range contacts {
		if contact.ID == id {
			return contact, nil
		}
	}

	return models.Contact{}, err
}

func UpdateContactById(id int, updatedContact models.Contact) (models.Contact, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return models.Contact{}, err
	}

	var found bool
	var updatedList []models.Contact

	for _, contact := range contacts {
		if contact.ID == id {
			updatedContact.ID = id
			updatedList = append(updatedList, updatedContact)
			found = true
		} else {
			updatedList = append(updatedList, contact)
		}
	}

	if !found {
		return models.Contact{}, err

	}

	if err := storage.SaveContacts(updatedList); err != nil {
		return models.Contact{}, err
	}

	return updatedContact, nil
}

func DeleteContactById(id int) error {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return err
	}

	var updatedContacts []models.Contact
	found := false

	for _, contact := range contacts {
		if contact.ID != id {
			updatedContacts = append(updatedContacts, contact)
		} else {
			found = true
		}
	}

	if !found {
		return err
	}

	if err := storage.SaveContacts(updatedContacts); err != nil {
		return err
	}

	return nil
}
