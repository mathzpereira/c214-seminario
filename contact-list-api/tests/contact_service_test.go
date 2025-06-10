package services_test

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
	"github.com/mathzpereira/c214-seminario/contact-list-api/services"
	"github.com/mathzpereira/c214-seminario/contact-list-api/storage"
)

// --- Testes para GetAllContacts ---

func TestGetAllContacts_Success(t *testing.T) {
	// Arrange
	expectedContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return expectedContacts, nil
	})
	defer patch.Unpatch()

	// Act
	contacts, err := services.GetAllContacts()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedContacts, contacts)
}

func TestGetAllContacts_StorageError(t *testing.T) {
	// Arrange
	expectedError := errors.New("database connection failed")

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, expectedError
	})
	defer patch.Unpatch()

	// Act
	contacts, err := services.GetAllContacts()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, contacts)
	assert.Equal(t, expectedError, err)
}

func TestGetAllContacts_FileNotFound(t *testing.T) {
	// Arrange
	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storage.ErrFileNotFound
	})
	defer patch.Unpatch()

	// Act
	contacts, err := services.GetAllContacts()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, contacts)
	assert.NotNil(t, contacts) // Deve ser uma fatia vazia, não nil
}

// --- Testes para AddContact ---

func TestAddContact_Success(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "11987654321"},
	}
	newContact := models.Contact{Name: "Bob", Email: "bob@test.com", Phone: "21998765432"}
	expectedContactWithID := models.Contact{ID: 2, Name: "Bob", Email: "bob@test.com", Phone: "21998765432"}
	contactsAfterAdd := append(existingContacts, expectedContactWithID)

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	var capturedSavedContacts []models.Contact
	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		capturedSavedContacts = contacts
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	addedContact, err := services.AddContact(newContact)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedContactWithID, addedContact)
	assert.Equal(t, contactsAfterAdd, capturedSavedContacts)
}

func TestAddContact_EmptyName(t *testing.T) {
	// Arrange
	newContact := models.Contact{Name: "", Email: "test@test.com", Phone: "1234567890"}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		assert.Fail(t, "LoadContacts não deveria ser chamado para validação de nome vazio")
		return nil, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado para validação de nome vazio")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.AddContact(newContact)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contact name cannot be empty")
}

func TestAddContact_InvalidEmail(t *testing.T) {
	// Arrange
	newContact := models.Contact{Name: "Test", Email: "invalid-email", Phone: "1234567890"}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		assert.Fail(t, "LoadContacts não deveria ser chamado para validação de email inválido")
		return nil, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado para validação de email inválido")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.AddContact(newContact)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email format")
}

func TestAddContact_InvalidPhone(t *testing.T) {
	// Arrange
	newContact := models.Contact{Name: "Test", Email: "valid@email.com", Phone: "123"}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		assert.Fail(t, "LoadContacts não deveria ser chamado para validação de telefone inválido")
		return nil, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado para validação de telefone inválido")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.AddContact(newContact)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid phone format")
}

func TestAddContact_StorageLoadError(t *testing.T) {
	// Arrange
	newContact := models.Contact{Name: "Bob", Email: "bob@test.com", Phone: "21998765432"}
	storageError := errors.New("failed to load contacts from storage")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storageError
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado se LoadContacts falhar")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.AddContact(newContact)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
}

func TestAddContact_StorageSaveError(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{}
	newContact := models.Contact{Name: "Bob", Email: "bob@test.com", Phone: "21998765432"}
	expectedContactWithID := models.Contact{ID: 1, Name: "Bob", Email: "bob@test.com", Phone: "21998765432"}
	contactsToSave := append(existingContacts, expectedContactWithID)
	storageError := errors.New("failed to save contacts to storage")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Equal(t, contactsToSave, contacts)
		return storageError
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.AddContact(newContact)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
}

// --- Testes para GetContactByID ---

func TestGetContactByID_Success(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 3, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
		{ID: 4, Name: "Marcos Vinícius", Email: "marcos.vinicius@gmail.com", Phone: ""},
	}
	expectedContact := mockContacts[1] // Carlos Eduardo (ID: 3)

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	result, err := services.GetContactByID(3)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedContact, result)
}

func TestGetContactByID_NotFound(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 3, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	result, err := services.GetContactByID(99) // ID que não existe

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contact not found")
	assert.Equal(t, models.Contact{}, result)
}

func TestGetContactByID_StorageError(t *testing.T) {
	// Arrange
	expectedError := errors.New("storage read error")

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, expectedError
	})
	defer patch.Unpatch()

	// Act
	result, err := services.GetContactByID(1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, models.Contact{}, result)
}

// --- Testes para UpdateContactById ---

func TestUpdateContactById_Success(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "1111111111"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Phone: "2222222222"},
	}
	updatedContactInput := models.Contact{Name: "Alicia Updated", Email: "alicia.updated@example.com", Phone: "11987654321"}
	expectedContactAfterUpdate := models.Contact{ID: 1, Name: "Alicia Updated", Email: "alicia.updated@example.com", Phone: "11987654321"}

	contactsAfterUpdate := []models.Contact{
		expectedContactAfterUpdate,
		existingContacts[1],
	}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	var capturedSavedContacts []models.Contact
	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		capturedSavedContacts = contacts
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	contact, err := services.UpdateContactById(1, updatedContactInput)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedContactAfterUpdate, contact)
	assert.Equal(t, contactsAfterUpdate, capturedSavedContacts)
}

func TestUpdateContactById_NotFound(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
	}
	updatedContactInput := models.Contact{Name: "NonExistent", Email: "no@example.com", Phone: "1234567890"}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado se o contato não for encontrado")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.UpdateContactById(99, updatedContactInput)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contact not found for update")
}

func TestUpdateContactById_InvalidName(t *testing.T) {
	// Arrange
	updatedContactInput := models.Contact{Name: "", Email: "valid@email.com", Phone: "1234567890"}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		assert.Fail(t, "LoadContacts não deveria ser chamado para validação")
		return nil, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado para validação")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.UpdateContactById(1, updatedContactInput)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contact name cannot be empty for update")
}

func TestUpdateContactById_InvalidEmail(t *testing.T) {
	// Arrange
	updatedContactInput := models.Contact{Name: "Valid Name", Email: "invalid-email", Phone: "1234567890"}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		assert.Fail(t, "LoadContacts não deveria ser chamado para validação")
		return nil, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado para validação")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.UpdateContactById(1, updatedContactInput)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email format for update")
}

func TestUpdateContactById_InvalidPhone(t *testing.T) {
	// Arrange
	updatedContactInput := models.Contact{Name: "Valid Name", Email: "valid@email.com", Phone: "123"}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		assert.Fail(t, "LoadContacts não deveria ser chamado para validação")
		return nil, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado para validação")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.UpdateContactById(1, updatedContactInput)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid phone format for update")
}

func TestUpdateContactById_StorageLoadError(t *testing.T) {
	// Arrange
	updatedContactInput := models.Contact{Name: "Alice", Email: "alice@example.com", Phone: "1111111111"}
	storageError := errors.New("storage load failed")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storageError
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado se LoadContacts falhar")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.UpdateContactById(1, updatedContactInput)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
}

func TestUpdateContactById_StorageSaveError(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{{ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "1111111111"}}
	updatedContactInput := models.Contact{Name: "Alice Updated", Email: "alice.updated@example.com", Phone: "11987654321"}
	expectedContactAfterUpdate := models.Contact{ID: 1, Name: "Alice Updated", Email: "alice.updated@example.com", Phone: "11987654321"}
	contactsAfterUpdate := []models.Contact{expectedContactAfterUpdate}
	storageError := errors.New("storage save failed")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Equal(t, contactsAfterUpdate, contacts)
		return storageError
	})
	defer patchSave.Unpatch()

	// Act
	_, err := services.UpdateContactById(1, updatedContactInput)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
}

// --- Testes para DeleteContactById ---

func TestDeleteContactById_Success(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "111111111"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Phone: "222222222"},
		{ID: 3, Name: "Marcos Vinícius", Email: "marcos@example.com", Phone: "333333333"},
	}
	contactsAfterDelete := []models.Contact{existingContacts[0], existingContacts[2]}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	var capturedSavedContacts []models.Contact
	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		capturedSavedContacts = contacts
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	err := services.DeleteContactById(2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, contactsAfterDelete, capturedSavedContacts)
}

func TestDeleteContactById_NotFound(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
	}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado se o contato não for encontrado")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	err := services.DeleteContactById(99)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contact not found for deletion")
}

func TestDeleteContactById_StorageLoadError(t *testing.T) {
	// Arrange
	storageError := errors.New("storage load error")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storageError
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Fail(t, "SaveContacts não deveria ser chamado se LoadContacts falhar")
		return nil
	})
	defer patchSave.Unpatch()

	// Act
	err := services.DeleteContactById(1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
}

func TestDeleteContactById_StorageSaveError(t *testing.T) {
	// Arrange
	existingContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
	}
	contactsAfterDelete := []models.Contact{}
	storageError := errors.New("storage save error")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return existingContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		assert.Equal(t, contactsAfterDelete, contacts)
		return storageError
	})
	defer patchSave.Unpatch()

	// Act
	err := services.DeleteContactById(1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
}

// --- Testes para GetContactsSummary ---

func TestGetContactsSummary_Success(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "1111111111"},
		{ID: 2, Name: "Bob", Email: "", Phone: "2222222222"},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", Phone: ""},
		{ID: 4, Name: "Alice", Email: "alice2@example.com", Phone: "4444444444"},
	}
	expectedSummary := services.ContactSummary{
		Total:           4,
		WithEmail:       3,
		WithPhone:       3,
		LastContactName: "Alice",
		DuplicatedNames: []string{"alice"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	summary, err := services.GetContactsSummary()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedSummary.Total, summary.Total)
	assert.Equal(t, expectedSummary.WithEmail, summary.WithEmail)
	assert.Equal(t, expectedSummary.WithPhone, summary.WithPhone)
	assert.Equal(t, expectedSummary.LastContactName, summary.LastContactName)

	actualDuplicatedNames := make([]string, len(summary.DuplicatedNames))
	copy(actualDuplicatedNames, summary.DuplicatedNames)
	expectedDuplicatedNames := make([]string, len(expectedSummary.DuplicatedNames))
	copy(expectedDuplicatedNames, expectedSummary.DuplicatedNames)

	for i := range actualDuplicatedNames {
		actualDuplicatedNames[i] = strings.ToLower(actualDuplicatedNames[i])
	}
	for i := range expectedDuplicatedNames {
		expectedDuplicatedNames[i] = strings.ToLower(expectedDuplicatedNames[i])
	}

	sort.Strings(actualDuplicatedNames)
	sort.Strings(expectedDuplicatedNames)
	assert.Equal(t, expectedDuplicatedNames, actualDuplicatedNames)
}

func TestGetContactsSummary_EmptyContacts(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{}
	expectedSummary := services.ContactSummary{
		Total:           0,
		WithEmail:       0,
		WithPhone:       0,
		LastContactName: "",
		DuplicatedNames: nil,
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	summary, err := services.GetContactsSummary()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, summary.DuplicatedNames)
	assert.Equal(t, expectedSummary, summary)
}

func TestGetContactsSummary_StorageError(t *testing.T) {
	// Arrange
	storageError := errors.New("summary storage error")

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storageError
	})
	defer patch.Unpatch()

	// Act
	summary, err := services.GetContactsSummary()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
	assert.Equal(t, services.ContactSummary{}, summary)
}

func TestGetContactsSummary_FileNotFound(t *testing.T) {
	// Arrange
	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storage.ErrFileNotFound
	})
	defer patch.Unpatch()

	// Act
	summary, err := services.GetContactsSummary()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, summary.Total)
	assert.Empty(t, summary.LastContactName)
	assert.Empty(t, summary.DuplicatedNames)
	assert.Equal(t, services.ContactSummary{}, summary)
}

// --- Testes para SearchContactsByName ---

func TestSearchContactsByName_Success(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
		{ID: 3, Name: "Alicia", Email: "alicia@example.com"},
	}
	expectedResults := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 3, Name: "Alicia", Email: "alicia@example.com"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	results, err := services.SearchContactsByName("ali")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, results)
}

func TestSearchContactsByName_NoMatch(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	results, err := services.SearchContactsByName("xyz")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestSearchContactsByName_EmptyName(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 2, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	results, err := services.SearchContactsByName("") // Busca por nome vazio

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockContacts, results)
}

func TestSearchContactsByName_StorageError(t *testing.T) {
	// Arrange
	storageError := errors.New("search storage error")

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storageError
	})
	defer patch.Unpatch()

	// Act
	_, err := services.SearchContactsByName("test")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
}

// --- Testes para GetEmailProviders ---

func TestGetEmailProviders_Success(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@gmail.com"},
		{ID: 3, Name: "Charlie", Email: "charlie@outlook.com"},
		{ID: 4, Name: "David", Email: "david@example.com"},
		{ID: 5, Name: "Eve", Email: "invalid-email"},
		{ID: 6, Name: "Frank", Email: ""},
	}
	expectedProviders := map[string]int{
		"example.com": 2,
		"gmail.com":   1,
		"outlook.com": 1,
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	providers, err := services.GetEmailProviders()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProviders, providers)
}

func TestGetEmailProviders_EmptyContacts(t *testing.T) {
	// Arrange
	mockContacts := []models.Contact{}
	expectedProviders := map[string]int{}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Act
	providers, err := services.GetEmailProviders()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, providers)
	assert.Equal(t, expectedProviders, providers)
}

func TestGetEmailProviders_StorageError(t *testing.T) {
	// Arrange
	storageError := errors.New("providers storage error")

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, storageError
	})
	defer patch.Unpatch()

	// Act
	providers, err := services.GetEmailProviders()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, storageError, err)
	assert.Nil(t, providers)
}
