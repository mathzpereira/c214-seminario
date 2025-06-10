package storage

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Join(filepath.Dir(b), "..", "data")
	dataFile   = filepath.Join(basePath, "contacts.json")
)

func LoadContacts() ([]models.Contact, error) {
	var contacts []models.Contact
	file, err := os.OpenFile(dataFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return contacts, err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	if len(byteValue) == 0 {
		return contacts, nil
	}

	err = json.Unmarshal(byteValue, &contacts)
	return contacts, err
}

func SaveContacts(contacts []models.Contact) error {
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}
