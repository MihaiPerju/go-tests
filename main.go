package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/machinebox/graphql"
)

type Partner struct {
	Id         string `json:"id"`
	PositionID string `json:"positionID"`
}

type Customer struct {
	AccountID     string `json:"accountID"`
	AccountNumber string `json:"accountNumber"`
}

type Token struct {
	ExpiresAt string `json:"expiresAt"`
}

type Principal struct {
	IsPartner  bool     `json:"isPartner"`
	IsCustomer bool     `json:"isCustomer"`
	Partner    Partner  `json:"partner"`
	Customer   Customer `json:"customer"`
}

type WhoAmI struct {
	Token     Token     `json:"token"`
	Principal Principal `json:"principal"`
}

type Data struct {
	WhoAmI WhoAmI `json:"whoAmI"`
}

func queryGraphQL(Token string) (*Data, error) {
	client := graphql.NewClient("https://clubhouse.dev.uw.systems/graphql")

	req := graphql.NewRequest(`
		query WhoAmI {
			whoAmI {
				token {
					expiresAt
				}
				principal {
					...on StaffPrincipal{
						id
						email
						organisation
					}
					...on UserPrincipal{
						isPartner
						isCustomer
						customer{
							accountID
							accountNumber
						}
						partner{
							id
							positionID
						}
					}
				}
			}
		}
	`)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept", "application/json")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	respData := Data{}

	if err := client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}

	return &respData, nil
}

func customerTokenTest() error {

	token := "eyJhbGciOiJSUzUxMiIsImtpZCI6IlNIQTI1NjowWEtNVkpBQWtPWVlUVkZWUmg0dzRnUUhQVHlBaFAyRGxSVmQ5aDlRMzdRIiwidHlwIjoiSldUIn0.eyJhdWQiOiJhMTAyODJlMy03MWI1LTQxN2ItYjkwMi1kMzMwZmE0YmMyNzAiLCJleHAiOjE2NzY1MzkwNTUsImlhdCI6MTY3NjUzODk5NSwiaXNzIjoiY3VzdG9tZXItYXV0aC1wcm92aWRlciIsImFjY291bnQubnVtYmVyIjoiMDAwMDAwMCIsImFjY291bnQuaWQiOiI1ZWUzY2Y5Mi02NGQ3LTVlOGItOGY2OC02NDg3MWI5NTUzNjEiLCJwZXJzb24uaWQiOiIyMjUyYTM3NS1lN2MwLTU2YzEtYWZhYS0xMDc4MGNhYjdiYjciLCJzY29wZXMiOlsiY3VzdG9tZXIuMDAwMDAwMC5yZWFkIiwiY3VzdG9tZXIuMDAwMDAwMC53cml0ZSIsImZ1bGwtYWNjZXNzIl0sInRva2VuLmlkIjoiLV94VlR6SGNsTHFBRkJlM29rT1VVVGQ4ZlJONjlKbFRScVJJajhrUk9IS2htdjc3aFlEX2tRPT0ifQ.PENld68e3vQO6R_WiRkNmOSuOFrVRIUZyr5bbgU78wNp-R75yunOXbym8hNqyLKAQKWODWzhPWDRHbDS6zkwnk1jqzirCSIJ0wBxfy04DxkIL8PYHbTConO2G1rPu4gETr6zOakz9M6js_fhhfy9PdRYfCtxLoNhIu9i5p7ivOIdmnz-8CloBA9GsI94bp3AoViWmsT_d-QqLFz66wfTwKE0KMmLm7iGHjvOt82Sv0AQpe-drSarL5Ri-t1iNG6Fds2VGkw4XDgVyCi7pUjIOHKJgsKr9N3E22dNauOjSO5DSa2_eY98BB-1irIx8CTIU-mDwq1WBwwQ1kA_TLkalA"
	resp, err := queryGraphQL(token)
	if err != nil {
		return err
	}

	// check token field exists
	if resp.WhoAmI.Token.ExpiresAt == "" {
		return errors.New("Token field not returned from GraphQL query")
	}

	// check the customer fields are as expected
	if !resp.WhoAmI.Principal.IsCustomer {
		return errors.New("isCustomer field should be true")
	}
	if resp.WhoAmI.Principal.Customer.AccountID == "" {
		return errors.New("Customer AccountId not found")
	}
	if resp.WhoAmI.Principal.Customer.AccountNumber == "" {
		return errors.New("Customer AccountNumber not found")
	}

	// check the partner fields are nul
	if resp.WhoAmI.Principal.IsPartner {
		return errors.New("isCustomer field should be false")
	}
	if resp.WhoAmI.Principal.Partner.Id != "" {
		return errors.New("Partner Id found, expected empty string")
	}
	if resp.WhoAmI.Principal.Partner.PositionID != "" {
		return errors.New("Partner PositionID found, expected empty string")
	}

	return nil
}

func principalTokenTest() error {

	token := ""
	resp, err := queryGraphQL(token)
	if err != nil {
		return err
	}

	// check token field exists
	if resp.WhoAmI.Token.ExpiresAt == "" {
		return errors.New("Token field not returned from GraphQL query")
	}

	// check the partner fields are as expected
	if !resp.WhoAmI.Principal.IsPartner {
		return errors.New("isCustomer field should be true")
	}
	if resp.WhoAmI.Principal.Partner.Id == "" {
		return errors.New("Partner Id not found")
	}
	if resp.WhoAmI.Principal.Partner.PositionID == "" {
		return errors.New("Partner PositionID not found")
	}

	// check the customer fields are nul
	if resp.WhoAmI.Principal.IsCustomer {
		return errors.New("isCustomer field should be true")
	}
	if resp.WhoAmI.Principal.Customer.AccountID != "" {
		return errors.New("Customer AccountId found, expected empty string")
	}
	if resp.WhoAmI.Principal.Customer.AccountNumber != "" {
		return errors.New("Customer AccountNumber found, expected empty string")
	}

	return nil
}

func main() {

	if err := customerTokenTest(); err != nil {
		fmt.Println("Customer token test fail", err)
	}
	if err := principalTokenTest(); err != nil {
		fmt.Println("Principal token test fail", err)
	}

	fmt.Println("Done.")

}
