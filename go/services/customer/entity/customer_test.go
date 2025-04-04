package entity

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tusmasoma/go-microservice-k8s/proto/customer"
)

func TestEntity_NewCustomer(t *testing.T) {
	t.Parallel()
	customerID := uuid.New().String()
	patterns := []struct {
		name string
		arg  struct {
			id          string
			name        string
			email       string
			street      string
			city        string
			countryCode customer.CountryCode
		}
		want struct {
			customer *Customer
			err      error
		}
	}{
		{
			name: "Success",
			arg: struct {
				id          string
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				id:          customerID,
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: &Customer{
					Customer: customer.Customer{
						Id:          customerID,
						Name:        "John Doe",
						Email:       "john.doe@example.com",
						Street:      "1600 Pennsylvania Avenue NW",
						City:        "Washington",
						CountryCode: customer.CountryCode_JPN,
					},
				},
				err: nil,
			},
		},
		{
			name: "Fail: id is empty",
			arg: struct {
				id          string
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				id:          "",
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("id is required"),
			},
		},
		{
			name: "Fail: name is empty",
			arg: struct {
				id          string
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				id:          customerID,
				name:        "",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("name is required"),
			},
		},
		{
			name: "Fail: email is empty",
			arg: struct {
				id          string
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				id:          customerID,
				name:        "John Doe",
				email:       "",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("email is required"),
			},
		},
		{
			name: "Fail: street is empty",
			arg: struct {
				id          string
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				id:          customerID,
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("street is required"),
			},
		},
		{
			name: "Fail: city is empty",
			arg: struct {
				id          string
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				id:          customerID,
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("city is required"),
			},
		},
		{
			name: "Fail: countryCode is empty",
			arg: struct {
				id          string
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				id:          customerID,
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_COUNTRY_CODE_UNSPECIFIED,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("countryCode is required"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			customer, err := NewCustomer(tt.arg.id, tt.arg.name, tt.arg.email, tt.arg.street, tt.arg.city, tt.arg.countryCode)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewCustomer() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewCustomer() error = %v, wantErr %v", err, tt.want.err)
			}
			if tt.want.err != nil {
				return
			}
			assert.Equal(t, customer.GetId(), tt.want.customer.GetId())
			assert.Equal(t, customer.GetName(), tt.want.customer.GetName())
			assert.Equal(t, customer.GetEmail(), tt.want.customer.GetEmail())
			assert.Equal(t, customer.GetStreet(), tt.want.customer.GetStreet())
			assert.Equal(t, customer.GetCity(), tt.want.customer.GetCity())
			assert.Equal(t, customer.GetCountryCode(), tt.want.customer.GetCountryCode())
		})
	}
}

func TestEntity_CreateCustomer(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name string
		arg  struct {
			name        string
			email       string
			street      string
			city        string
			countryCode customer.CountryCode
		}
		want struct {
			customer *Customer
			err      error
		}
	}{
		{
			name: "Success",
			arg: struct {
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: &Customer{
					Customer: customer.Customer{
						Name:        "John Doe",
						Email:       "john.doe@example.com",
						Street:      "1600 Pennsylvania Avenue NW",
						City:        "Washington",
						CountryCode: customer.CountryCode_JPN,
					},
				},
				err: nil,
			},
		},
		{
			name: "Fail: name is empty",
			arg: struct {
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				name:        "",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("name is required"),
			},
		},
		{
			name: "Fail: email is empty",
			arg: struct {
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				name:        "John Doe",
				email:       "",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("email is required"),
			},
		},
		{
			name: "Fail: street is empty",
			arg: struct {
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "",
				city:        "Washington",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("street is required"),
			},
		},
		{
			name: "Fail: city is empty",
			arg: struct {
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "",
				countryCode: customer.CountryCode_JPN,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("city is required"),
			},
		},
		{
			name: "Fail: countryCode is empty",
			arg: struct {
				name        string
				email       string
				street      string
				city        string
				countryCode customer.CountryCode
			}{
				name:        "John Doe",
				email:       "john.doe@example.com",
				street:      "1600 Pennsylvania Avenue NW",
				city:        "Washington",
				countryCode: customer.CountryCode_COUNTRY_CODE_UNSPECIFIED,
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("countryCode is required"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			customer, err := CreateCustomer(tt.arg.name, tt.arg.email, tt.arg.street, tt.arg.city, tt.arg.countryCode)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("CreateCustomer() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateCustomer() error = %v, wantErr %v", err, tt.want.err)
			}
			if tt.want.err != nil {
				return
			}
			assert.Equal(t, customer.GetName(), tt.want.customer.GetName())
			assert.Equal(t, customer.GetEmail(), tt.want.customer.GetEmail())
			assert.Equal(t, customer.GetStreet(), tt.want.customer.GetStreet())
			assert.Equal(t, customer.GetCity(), tt.want.customer.GetCity())
			assert.Equal(t, customer.GetCountryCode(), tt.want.customer.GetCountryCode())
		})
	}
}

func Test_Customer_Proto(t *testing.T) {
	user := Customer{
		Customer: customer.Customer{
			Id:          "1",
			Name:        "John Doe",
			Email:       "john.doe@example.com",
			Street:      "1600 Pennsylvania Avenue NW",
			City:        "Washington",
			CountryCode: customer.CountryCode_JPN,
		},
	}
	expect := &customer.Customer{
		Id:          "1",
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Street:      "1600 Pennsylvania Avenue NW",
		City:        "Washington",
		CountryCode: customer.CountryCode_JPN,
	}
	assert.Equal(t, expect, user.Proto())
}

func Test_Customers_Proto(t *testing.T) {
	users := Customers{
		&Customer{
			customer.Customer{
				Id:          "1",
				Name:        "John Doe",
				Email:       "john.doe@example.com",
				Street:      "1600 Pennsylvania Avenue NW",
				City:        "Washington",
				CountryCode: customer.CountryCode_JPN,
			},
		},
		&Customer{
			customer.Customer{
				Id:          "2",
				Name:        "Ken Doe",
				Email:       "ken.doe@example.com",
				Street:      "1600 Pennsylvania Avenue NW",
				City:        "Washington",
				CountryCode: customer.CountryCode_JPN,
			},
		},
	}
	expect := []*customer.Customer{
		{
			Id:          "1",
			Name:        "John Doe",
			Email:       "john.doe@example.com",
			Street:      "1600 Pennsylvania Avenue NW",
			City:        "Washington",
			CountryCode: customer.CountryCode_JPN,
		},
		{
			Id:          "2",
			Name:        "Ken Doe",
			Email:       "ken.doe@example.com",
			Street:      "1600 Pennsylvania Avenue NW",
			City:        "Washington",
			CountryCode: customer.CountryCode_JPN,
		},
	}
	assert.Equal(t, expect, users.Proto())
}
