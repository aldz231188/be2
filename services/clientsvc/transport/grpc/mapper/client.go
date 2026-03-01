package mapper

import (
	clientv1 "be2/contracts/gen/client/v1"
	// "be2/services/clientsvc/app"
	"be2/services/clientsvc/domain"
	// "be2/services/domain/shared/date"
	"strings"
	// "time"
	"github.com/google/uuid"
)

type CreateClientRequest struct {
	P *clientv1.CreateClientRequest
}

// type CreateClientResponse struct {
// 	P *clientv1.CreateClientResponse
// }

// type CreateClientInput struct {
// 	P *app.CreateClientInput
// }

// type CreateClientRequest struct {
// 	ClientName    string               `json:"client_name"`
// 	ClientSurname string               `json:"client_surname"`
// 	Birthday      date.DateOnly        `json:"birthday"`
// 	Gender        string               `json:"gender"`
// 	Address       CreateAddressRequest `json:"address"`
// }

// type ClientResponse struct {
// 	ID               string          `json:"id"`
// 	ClientName       string          `json:"client_name"`
// 	ClientSurname    string          `json:"client_surname"`
// 	Birthday         date.DateOnly   `json:"birthday"`
// 	Gender           string          `json:"gender"`
// 	RegistrationDate date.DateOnly   `json:"registration_date"`
// 	Address          AddressResponse `json:"address"`
// }

func (r CreateClientRequest) ToDomainClient() (domain.Client, error) {
	errs := domain.NewValidationErrors()

	userid := strings.TrimSpace(r.P.GetUserid())
	if userid == "" {
		errs.Add("user_id", "is required")
	}
	name := strings.TrimSpace(r.P.GetName())
	if name == "" {
		errs.Add("client_name", "is required")
	}

	surname := strings.TrimSpace(r.P.GetSurname())
	if surname == "" {
		errs.Add("client_surname", "is required")
	}

	// if r.Birthday.Time.IsZero() {
	// 	errs.Add("birthday", "is required")
	// } else if r.Birthday.Time.After(time.Now().UTC()) {
	// 	errs.Add("birthday", "cannot be in the future")
	// }

	// gender, err := parseGender(r.Gender)
	// if err != nil {
	// 	if validationErrs, ok := err.(*domain.ValidationErrors); ok {
	// 		errs.Merge(validationErrs)
	// 	} else {
	// 		return app.CreateClientInput{}, err
	// 	}
	// }

	// address, err := r.Address.ToDomainAddress()
	// if err != nil {
	// 	if validationErrs, ok := err.(*domain.ValidationErrors); ok {
	// 		errs.Merge(validationErrs)
	// 	} else {
	// 		return app.CreateClientInput{}, err
	// 	}
	// }

	if errs.HasErrors() {
		return domain.Client{}, errs
	}

	return domain.Client{
		UserID:        uuid.MustParse(userid),
		ClientName:    name,
		ClientSurname: surname,
		// Birthday:      r.Birthday,
		// Gender:        gender,
		// Address:       address,
	}, nil
}

func FromDomainClient(c domain.Client) clientv1.CreateClientResponse {
	return clientv1.CreateClientResponse{
		Clientid: c.ID.String(),
	}

}

func parseGender(s string) (domain.Gender, error) {
	value := strings.TrimSpace(s)
	switch {
	case strings.EqualFold(value, "male"):
		return domain.MALE, nil
	case strings.EqualFold(value, "female"):
		return domain.FEMALE, nil
	case value == "":
		errs := domain.NewValidationErrors()
		errs.Add("gender", "is required")
		return domain.MALE, errs
	default:
		errs := domain.NewValidationErrors()
		errs.Add("gender", "must be either male or female")
		return domain.MALE, errs
	}
}
