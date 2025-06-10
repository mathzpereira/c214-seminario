package services

import (
	"errors"
	"regexp" // Importe se estiver usando regex para validação
	"sort"
	"strings" // Importe se estiver usando strings.ToLower ou similar

	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
	"github.com/mathzpereira/c214-seminario/contact-list-api/storage"
)

// ContactSummary é uma struct auxiliar para o GetContactsSummary
type ContactSummary struct {
	Total           int
	WithEmail       int
	WithPhone       int
	LastContactName string
	DuplicatedNames []string
}

// NewContactService foi removido na versão monkey.Patch, pois as funções são chamadas diretamente.
// Se você ainda tem uma struct ContactService que encapsula o storage, remova-a para a versão monkey.Patch
// ou adapte para que ela chame as funções globais do storage.

// A validação de contato pode ser uma função auxiliar
func validateContact(contact models.Contact) error {
	if contact.Name == "" {
		return errors.New("contact name cannot be empty")
	}
	if contact.Email != "" {
		if match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, contact.Email); !match {
			return errors.New("invalid email format")
		}
	}
	if contact.Phone != "" {
		if match, _ := regexp.MatchString(`^[0-9]{10,15}$`, contact.Phone); !match { // Exemplo de regex: 10 a 15 dígitos
			return errors.New("invalid phone format")
		}
	}
	return nil
}

// GetAllContacts retorna todos os contatos.
func GetAllContacts() ([]models.Contact, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		// Se o arquivo não foi encontrado, retorne uma lista vazia e nenhum erro.
		if err == storage.ErrFileNotFound {
			return []models.Contact{}, nil
		}
		// Para qualquer outro erro do storage, retorne o erro.
		return nil, err
	}
	return contacts, nil
}

// AddContact adiciona um novo contato.
func AddContact(newContact models.Contact) (models.Contact, error) {
	// Validação inicial antes de interagir com o storage
	if err := validateContact(newContact); err != nil {
		return models.Contact{}, err
	}

	contacts, err := storage.LoadContacts()
	if err != nil && err != storage.ErrFileNotFound {
		return models.Contact{}, err
	}
	if err == storage.ErrFileNotFound {
		contacts = []models.Contact{} // Inicializa como fatia vazia se o arquivo não existe
	}

	// Gerar o próximo ID
	var nextID int
	if len(contacts) > 0 {
		nextID = contacts[len(contacts)-1].ID + 1
	} else {
		nextID = 1
	}
	newContact.ID = nextID

	contacts = append(contacts, newContact)

	err = storage.SaveContacts(contacts)
	if err != nil {
		return models.Contact{}, err
	}

	return newContact, nil // Retorna o contato com o ID gerado
}

// GetContactByID retorna um contato pelo seu ID.
func GetContactByID(id int) (models.Contact, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return models.Contact{}, err // Retorna erro do storage
	}

	for _, contact := range contacts {
		if contact.ID == id {
			return contact, nil
		}
	}
	return models.Contact{}, errors.New("contact not found")
}

// UpdateContactById atualiza um contato existente pelo ID.
func UpdateContactById(id int, updatedContact models.Contact) (models.Contact, error) {
	// Validação inicial
	if updatedContact.Name == "" {
		return models.Contact{}, errors.New("contact name cannot be empty for update")
	}
	if updatedContact.Email != "" {
		if match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, updatedContact.Email); !match {
			return models.Contact{}, errors.New("invalid email format for update")
		}
	}
	if updatedContact.Phone != "" {
		if match, _ := regexp.MatchString(`^[0-9]{10,15}$`, updatedContact.Phone); !match {
			return models.Contact{}, errors.New("invalid phone format for update")
		}
	}

	contacts, err := storage.LoadContacts()
	if err != nil {
		return models.Contact{}, err
	}

	found := false
	for i, contact := range contacts {
		if contact.ID == id {
			contacts[i].Name = updatedContact.Name
			contacts[i].Email = updatedContact.Email
			contacts[i].Phone = updatedContact.Phone
			found = true
			updatedContact.ID = id // Garante que o ID do retorno esteja correto
			break
		}
	}

	if !found {
		return models.Contact{}, errors.New("contact not found for update")
	}

	err = storage.SaveContacts(contacts)
	if err != nil {
		return models.Contact{}, err
	}
	return updatedContact, nil
}

// DeleteContactById deleta um contato pelo ID.
func DeleteContactById(id int) error {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return err
	}

	foundIndex := -1
	for i, contact := range contacts {
		if contact.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return errors.New("contact not found for deletion")
	}

	// Remove o contato da slice
	contacts = append(contacts[:foundIndex], contacts[foundIndex+1:]...)

	err = storage.SaveContacts(contacts)
	if err != nil {
		return err
	}
	return nil
}

// GetContactsSummary retorna estatísticas sobre os contatos.
func GetContactsSummary() (ContactSummary, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		if err == storage.ErrFileNotFound {
			return ContactSummary{}, nil // Retorna summary vazia se arquivo não existe
		}
		return ContactSummary{}, err
	}

	summary := ContactSummary{
		Total: len(contacts),
	}

	nameCounts := make(map[string]int)
	for _, contact := range contacts {
		if contact.Email != "" {
			summary.WithEmail++
		}
		if contact.Phone != "" {
			summary.WithPhone++
		}
		nameCounts[strings.ToLower(contact.Name)]++
	}

	for name, count := range nameCounts {
		if count > 1 {
			summary.DuplicatedNames = append(summary.DuplicatedNames, name)
		}
	}

	// Ordenar nomes duplicados para consistência nos testes
	sort.Strings(summary.DuplicatedNames)

	if len(contacts) > 0 {
		summary.LastContactName = contacts[len(contacts)-1].Name
	}

	return summary, nil
}

// SearchContactsByName busca contatos por nome.
func SearchContactsByName(name string) ([]models.Contact, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return nil, err
	}

	var results []models.Contact
	searchLower := strings.ToLower(name)
	for _, contact := range contacts {
		if strings.Contains(strings.ToLower(contact.Name), searchLower) {
			results = append(results, contact)
		}
	}
	return results, nil
}

// GetEmailProviders conta a ocorrência de cada provedor de e-mail.
func GetEmailProviders() (map[string]int, error) {
	contacts, err := storage.LoadContacts()
	if err != nil {
		return nil, err
	}

	providers := make(map[string]int)
	emailRegex := regexp.MustCompile(`@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`)

	for _, contact := range contacts {
		matches := emailRegex.FindStringSubmatch(contact.Email)
		if len(matches) > 1 {
			provider := strings.ToLower(matches[1])
			providers[provider]++
		}
	}
	return providers, nil
}
