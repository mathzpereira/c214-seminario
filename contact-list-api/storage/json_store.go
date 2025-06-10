package storage

import (
	"encoding/json"
	"errors" // Mantenha esta importação
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
)

// Variável para erro de arquivo não encontrado
var ErrFileNotFound = errors.New("file not found")

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Join(filepath.Dir(b), "..", "data")
	dataFile   = filepath.Join(basePath, "contacts.json")
)

func LoadContacts() ([]models.Contact, error) {
	var contacts []models.Contact
	// Tenta abrir o arquivo. Se não existir, os.IsNotExist(err) será true.
	file, err := os.OpenFile(dataFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		// Se o erro não for 'arquivo não encontrado' (por exemplo, permissão), retorne-o
		if os.IsNotExist(err) {
			return contacts, ErrFileNotFound // Retorne o seu erro específico
		}
		return contacts, err // Outro erro de abertura
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file) // 'err' agora é importante aqui
	if err != nil {
		return contacts, err
	}

	if len(byteValue) == 0 {
		return contacts, nil
	}

	err = json.Unmarshal(byteValue, &contacts)
	return contacts, err
}

func SaveContacts(contacts []models.Contact) error {
	data, err := json.MarshalIndent(contacts, "", "  ") // 2 espaços para indentação
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}
