package services

import (
	"strings"

	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
	"github.com/mathzpereira/c214-seminario/contact-list-api/storage"
)

type ContactSummary struct {
	Total           int      `json:"total"`
	WithEmail       int      `json:"with_email"`
	WithPhone       int      `json:"with_phone"`
	LastContactName string   `json:"last_contact_name,omitempty"`
	DuplicatedNames []string `json:"duplicated_names,omitempty"`
}

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

func GetContactsSummary() (ContactSummary, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return ContactSummary{}, err
	}

	summary := ContactSummary{
		Total:     len(contacts),
		WithEmail: 0,
		WithPhone: 0,
	}

	nameMap := make(map[string]int)

	for _, c := range contacts {
		if strings.TrimSpace(c.Email) != "" {
			summary.WithEmail++
		}
		if strings.TrimSpace(c.Phone) != "" {
			summary.WithPhone++
		}
		nameMap[strings.ToLower(c.Name)]++
	}

	if summary.Total > 0 {
		summary.LastContactName = contacts[len(contacts)-1].Name
	}

	for name, count := range nameMap {
		if count > 1 {
			summary.DuplicatedNames = append(summary.DuplicatedNames, name)
		}
	}

	return summary, nil
}

func SearchContactsByName(name string) ([]models.Contact, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return nil, err
	}

	var results []models.Contact
	for _, contact := range contacts {
		if strings.HasPrefix(strings.ToLower(contact.Name), strings.ToLower(name)) {
			results = append(results, contact)
		}
	}

	return results, nil
}

func GetEmailProviders() (map[string]int, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return nil, err
	}

	providers := make(map[string]int)

	for _, c := range contacts {
		if strings.TrimSpace(c.Email) == "" {
			continue
		}

		parts := strings.Split(c.Email, "@")
		if len(parts) != 2 {
			continue
		}

		domain := strings.ToLower(parts[1])
		providers[domain]++
	}

	return providers, nil
}