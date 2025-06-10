package service

import (
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/mathzpereira/c214-seminario/contact-list-api/models"
	"github.com/mathzpereira/c214-seminario/contact-list-api/services"
	"github.com/mathzpereira/c214-seminario/contact-list-api/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetContactByID_Success_ExpectedValidContact(t *testing.T) {
	// Fixture
	expectedContact := models.Contact{
		ID:    4,
		Name:  "Carlos Eduardo",
		Email: "carlos.eduardo@gmail.com",
		Phone: "551199998877",
	}

	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 3, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
		{ID: 4, Name: "Marcos Vinícius", Email: "marcos.vinicius@gmail.com", Phone: ""},
		{ID: 5, Name: "Juliana Souza", Email: "", Phone: "551197654321"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Exercise
	result, err := services.GetContactByID(3)

	// Assert
	assert.Equal(t, result, expectedContact)
	assert.NoError(t, err, "Expected no error when contact is found")
}

func TestGetContactByID_NotFound_ExpectedEmpty(t *testing.T) {
	// Fixture

	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 3, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
		{ID: 4, Name: "Marcos Vinícius", Email: "marcos.vinicius@gmail.com", Phone: ""},
		{ID: 5, Name: "Juliana Souza", Email: "", Phone: "551197654321"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Exercise
	result, err := services.GetContactByID(2)

	// Assert
	assert.Equal(t, models.Contact{}, result)
	assert.NoError(t, err)
}

func TestUpdateContactById_Success_ExpectedUpdatedContact(t *testing.T) {
	// Fixture
	updatedContact := models.Contact{
		Name:  "Carlos Eduardo Atualizado",
		Email: "carlos.eduardo.novo@gmail.com",
		Phone: "551199887766",
	}

	expectedContact := models.Contact{
		ID:    3,
		Name:  "Carlos Eduardo Atualizado",
		Email: "carlos.eduardo.novo@gmail.com",
		Phone: "551199887766",
	}

	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 3, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
		{ID: 4, Name: "Marcos Vinícius", Email: "marcos.vinicius@gmail.com", Phone: ""},
		{ID: 5, Name: "Juliana Souza", Email: "", Phone: "551197654321"},
	}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		return nil
	})
	defer patchSave.Unpatch()

	// Exercise
	result, err := services.UpdateContactById(3, updatedContact)

	// Assert
	assert.Equal(t, expectedContact, result)
	assert.NoError(t, err)
}

func TestUpdateContactById_NotFound_ExpectedEmpty(t *testing.T) {
	// Fixture
	updatedContact := models.Contact{
		Name:  "Contato Inexistente",
		Email: "inexistente@gmail.com",
		Phone: "551199999999",
	}

	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 3, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
		{ID: 4, Name: "Marcos Vinícius", Email: "marcos.vinicius@gmail.com", Phone: ""},
		{ID: 5, Name: "Juliana Souza", Email: "", Phone: "551197654321"},
	}

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		return nil
	})
	defer patchSave.Unpatch()

	// Exercise
	result, err := services.UpdateContactById(2, updatedContact)

	// Assert
	assert.Equal(t, models.Contact{}, result)
	assert.NoError(t, err)
}

func TestUpdateContactById_LoadError_ExpectedError(t *testing.T) {
	// Fixture
	updatedContact := models.Contact{
		Name:  "Teste Erro",
		Email: "teste@gmail.com",
		Phone: "551199999999",
	}

	expectedError := errors.New("failed to load contacts")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, expectedError
	})
	defer patchLoad.Unpatch()

	// Exercise
	result, err := services.UpdateContactById(1, updatedContact)

	// Assert
	assert.Equal(t, models.Contact{}, result)
	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func TestUpdateContactById_SaveError_ExpectedError(t *testing.T) {
	// Fixture
	updatedContact := models.Contact{
		Name:  "Teste Erro Save",
		Email: "teste.save@gmail.com",
		Phone: "551199888777",
	}

	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 3, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: "551199998877"},
		{ID: 4, Name: "Marcos Vinícius", Email: "marcos.vinicius@gmail.com", Phone: ""},
		{ID: 5, Name: "Juliana Souza", Email: "", Phone: "551197654321"},
	}

	expectedError := errors.New("failed to save contacts")

	patchLoad := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patchLoad.Unpatch()

	patchSave := monkey.Patch(storage.SaveContacts, func(contacts []models.Contact) error {
		return expectedError
	})
	defer patchSave.Unpatch()

	// Exercise
	result, err := services.UpdateContactById(3, updatedContact)

	// Assert
	assert.Equal(t, models.Contact{}, result)
	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func TestGetContactsSummary_Success_ExpectedCompleteStatistics(t *testing.T) {
	// Fixture
	mockContacts := []models.Contact{
		{ID: 1, Name: "Fernanda Lima", Email: "fernanda.lima@yahoo.com", Phone: "551198765432"},
		{ID: 2, Name: "Carlos Eduardo", Email: "carlos.eduardo@gmail.com", Phone: ""},
		{ID: 3, Name: "Marcos Vinícius", Email: "", Phone: "551197654321"},
		{ID: 4, Name: "Juliana Souza", Email: "juliana.souza@gmail.com", Phone: "551199887766"},
		{ID: 5, Name: "fernanda lima", Email: "fernanda.outro@gmail.com", Phone: "551188776655"},
	}

	expectedSummary := services.ContactSummary{
		Total:           5,
		WithEmail:       4,
		WithPhone:       4,
		LastContactName: "fernanda lima",
		DuplicatedNames: []string{"fernanda lima"},
	}

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return mockContacts, nil
	})
	defer patch.Unpatch()

	// Exercise
	result, err := services.GetContactsSummary()

	// Assert
	assert.Equal(t, result, expectedSummary)
	assert.NoError(t, err)
}

func TestGetContactsSummary_EmptyList_ExpectedZeroStatistics(t *testing.T) {
	// Fixture
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

	// Exercise
	result, err := services.GetContactsSummary()

	// Assert
	assert.Equal(t, result, expectedSummary)
	assert.NoError(t, err)
}

func TestGetContactsSummary_LoadError_ExpectedError(t *testing.T) {
	// Fixture
	expectedError := errors.New("failed to load contacts from storage")

	patch := monkey.Patch(storage.LoadContacts, func() ([]models.Contact, error) {
		return nil, expectedError
	})
	defer patch.Unpatch()

	// Exercise
	result, err := services.GetContactsSummary()

	// Assert
	assert.Equal(t, services.ContactSummary{}, result)
	assert.Error(t, err)
}
