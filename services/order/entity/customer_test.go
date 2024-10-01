package entity

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestEntity_NewCustomer(t *testing.T) {
	t.Parallel()

	customerID := uuid.New().String()

	patterns := []struct {
		name string
		arg  struct {
			id      string
			name    string
			email   string
			street  string
			city    string
			country string
		}
		want struct {
			customer *Customer
			err      error
		}
	}{
		{
			name: "success: id is not empty",
			arg: struct {
				id      string
				name    string
				email   string
				street  string
				city    string
				country string
			}{
				id:      customerID,
				name:    "John Doe",
				email:   "john.doe@example.com",
				street:  "1600 Pennsylvania Avenue NW",
				city:    "Washington",
				country: "USA",
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: &Customer{
					ID:      customerID,
					Name:    "John Doe",
					Email:   "john.doe@example.com",
					Street:  "1600 Pennsylvania Avenue NW",
					City:    "Washington",
					Country: "USA",
				},
				err: nil,
			},
		},
		{
			name: "success: id is empty",
			arg: struct {
				id      string
				name    string
				email   string
				street  string
				city    string
				country string
			}{
				name:    "John Doe",
				email:   "john.doe@example.com",
				street:  "1600 Pennsylvania Avenue NW",
				city:    "Washington",
				country: "USA",
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: &Customer{
					ID:      uuid.New().String(),
					Name:    "John Doe",
					Email:   "john.doe@example.com",
					Street:  "1600 Pennsylvania Avenue NW",
					City:    "Washington",
					Country: "USA",
				},
				err: nil,
			},
		},
		{
			name: "Fail: name is empty",
			arg: struct {
				id      string
				name    string
				email   string
				street  string
				city    string
				country string
			}{
				name:    "",
				email:   "john.doe@example.com",
				street:  "1600 Pennsylvania Avenue NW",
				city:    "Washington",
				country: "USA",
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
				id      string
				name    string
				email   string
				street  string
				city    string
				country string
			}{
				name:    "John Doe",
				email:   "",
				street:  "1600 Pennsylvania Avenue NW",
				city:    "Washington",
				country: "USA",
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
				id      string
				name    string
				email   string
				street  string
				city    string
				country string
			}{
				name:    "John Doe",
				email:   "john.doe@example.com",
				street:  "",
				city:    "Washington",
				country: "USA",
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
				id      string
				name    string
				email   string
				street  string
				city    string
				country string
			}{
				name:    "John Doe",
				email:   "john.doe@example.com",
				street:  "1600 Pennsylvania Avenue NW",
				city:    "",
				country: "USA",
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
			name: "Fail: country is empty",
			arg: struct {
				id      string
				name    string
				email   string
				street  string
				city    string
				country string
			}{
				name:    "John Doe",
				email:   "john.doe@example.com",
				street:  "1600 Pennsylvania Avenue NW",
				city:    "Washington",
				country: "",
			},
			want: struct {
				customer *Customer
				err      error
			}{
				customer: nil,
				err:      errors.New("country is required"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			customer, err := NewCustomer(tt.arg.id, tt.arg.name, tt.arg.email, tt.arg.street, tt.arg.city, tt.arg.country)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewCustomer() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewCustomer() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(customer, tt.want.customer, cmpopts.IgnoreFields(Customer{}, "ID")); len(d) != 0 {
				t.Errorf("NewCustomer() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
