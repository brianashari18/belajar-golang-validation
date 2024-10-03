package golang_validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestValidation(t *testing.T) {
	validate := validator.New()

	if validate == nil {
		t.Error("Validate is nil")
	}
}

func TestValidationVariable(t *testing.T) {
	validate := validator.New()
	user := "Brian"

	err := validate.Var(user, "required")
	if err != nil {
		t.Error("user is nil")
	}
}

func TestValidationTwoVariables(t *testing.T) {
	validate := validator.New()
	pw1 := "correct"
	pw2 := "incorrect"

	err := validate.VarWithValue(pw1, pw2, "eqfield")
	if err != nil {
		t.Error("password is wrong")
	}
}

func TestValidationMultipleTag(t *testing.T) {
	validate := validator.New()
	username := "12"

	err := validate.Var(username, "required,numeric")
	if err != nil {
		t.Error("username is not valid")
	}
}

func TestValidationTagParameter(t *testing.T) {
	validate := validator.New()
	username := "12345"

	err := validate.Var(username, "required,numeric,min=5,max=10")
	if err != nil {
		t.Error("username is not valid")
	}
}

func TestValidationStruct(t *testing.T) {
	type User struct {
		Username string `validate:"required,email"`
		Password string `validate:"required,min=5,max=10"`
	}

	validate := validator.New()
	user := User{
		Username: "Brian@example.com",
		Password: "12345",
	}

	err := validate.Struct(user)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationError(t *testing.T) {
	type User struct {
		Username string `validate:"required,email"`
		Password string `validate:"required,min=5,max=10"`
	}

	validate := validator.New()
	user := User{
		Username: "Brian",
		Password: "12",
	}

	err := validate.Struct(user)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				fmt.Println("error on", fieldError.Field(), "on tag", fieldError.Tag(), "with value", fieldError.Value())
			}
		}
	}
}

func TestValidationCrossField(t *testing.T) {
	type User struct {
		Username        string `validate:"required,email"`
		Password        string `validate:"required,min=5,max=10"`
		ConfirmPassword string `validate:"required,eqfield=Password"`
	}

	validate := validator.New()
	user := User{
		Username:        "Brian@example.com",
		Password:        "12345",
		ConfirmPassword: "1",
	}

	err := validate.Struct(user)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNestedStruct(t *testing.T) {
	type Address struct {
		Country string `validate:"required"`
		Street  string `validate:"required"`
	}

	type User struct {
		Username string  `validate:"required"`
		Password string  `validate:"required"`
		Address  Address `validate:"required"`
	}

	validate := validator.New()
	user := User{
		Username: "",
		Password: "",
		Address: Address{
			Country: "",
			Street:  "",
		},
	}

	err := validate.Struct(user)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationCollection(t *testing.T) {
	type Address struct {
		Country string `validate:"required"`
		Street  string `validate:"required"`
	}

	type User struct {
		Username  string    `validate:"required"`
		Password  string    `validate:"required"`
		Addresses []Address `validate:"required,dive"`
	}

	validate := validator.New()
	user := User{
		Username: "",
		Password: "",
		Addresses: []Address{
			Address{
				Country: "",
				Street:  "",
			},
			Address{
				Country: "",
				Street:  "",
			},
		},
	}

	err := validate.Struct(user)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationBasicCollection(t *testing.T) {
	type Address struct {
		Country string `validate:"required"`
		Street  string `validate:"required"`
	}

	type User struct {
		Username  string    `validate:"required"`
		Password  string    `validate:"required"`
		Addresses []Address `validate:"required,dive"`
		Hobbies   []string  `validate:"required,dive,required,min=3"`
	}

	validate := validator.New()
	user := User{
		Username: "",
		Password: "",
		Addresses: []Address{
			Address{
				Country: "",
				Street:  "",
			},
			Address{
				Country: "",
				Street:  "",
			},
		},
		Hobbies: []string{
			"Gaming",
			"Coding",
			"",
			"X",
		},
	}

	err := validate.Struct(user)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationMap(t *testing.T) {
	type Address struct {
		Country string `validate:"required"`
		Street  string `validate:"required"`
	}

	type School struct {
		Name string `validate:"required"`
	}

	type User struct {
		Username  string            `validate:"required"`
		Password  string            `validate:"required"`
		Addresses []Address         `validate:"required,dive"`
		Hobbies   []string          `validate:"required,dive,required,min=3"`
		Schools   map[string]School `validate:"required,dive,keys,required,min=3,endkeys,required"`
	}

	validate := validator.New()
	user := User{
		Username: "",
		Password: "",
		Addresses: []Address{
			Address{
				Country: "",
				Street:  "",
			},
			Address{
				Country: "",
				Street:  "",
			},
		},
		Hobbies: []string{
			"Gaming",
			"Coding",
			"",
			"X",
		},
		Schools: map[string]School{
			"SD": {
				Name: "SD",
			},
			"SMP": {
				Name: "SMP",
			},
			"SMA": {
				Name: "",
			},
		},
	}

	err := validate.Struct(user)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationBasicMap(t *testing.T) {
	type Address struct {
		Country string `validate:"required"`
		Street  string `validate:"required"`
	}

	type School struct {
		Name string `validate:"required"`
	}

	type User struct {
		Username  string            `validate:"required"`
		Password  string            `validate:"required"`
		Addresses []Address         `validate:"required,dive"`
		Hobbies   []string          `validate:"required,dive,required,min=3"`
		Schools   map[string]School `validate:"required,dive,keys,required,min=3,endkeys,required"`
		Wallets   map[string]int    `validate:"required,dive,keys,required,min=3,endkeys,required,gt=0"`
	}

	validate := validator.New()
	user := User{
		Username: "",
		Password: "",
		Addresses: []Address{
			Address{
				Country: "",
				Street:  "",
			},
			Address{
				Country: "",
				Street:  "",
			},
		},
		Hobbies: []string{
			"Gaming",
			"Coding",
			"",
			"X",
		},
		Schools: map[string]School{
			"SD": {
				Name: "SD",
			},
			"SMP": {
				Name: "SMP",
			},
			"SMA": {
				Name: "",
			},
		},
		Wallets: map[string]int{
			"BCA":     0,
			"Mandiri": 1000000,
			"":        10,
		},
	}

	err := validate.Struct(user)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlias(t *testing.T) {
	validate := validator.New()
	validate.RegisterAlias("varchar", "required,max=255")

	type Seller struct {
		Id   string `validate:"varchar"`
		Name string `validate:"varchar"`
	}

	seller := Seller{
		Id:   "",
		Name: "",
	}

	err := validate.Struct(seller)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func MustValidUsername(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		if value != strings.ToUpper(value) {
			return false
		}
		if len(value) < 5 {
			return false
		}
	}
	return true
}

func TestCustomValidation(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("username", MustValidUsername)
	if err != nil {
		return
	}

	type RegisterUser struct {
		Name     string `validate:"username"`
		Password string `validate:"required"`
	}

	user := RegisterUser{
		Name:     "BRIAN",
		Password: "",
	}

	err = validate.Struct(user)
	if err != nil {
		fmt.Println(err.Error())
	}
}

var regexNumber = regexp.MustCompile(`^[0-9]+$`)

func MustValidPin(field validator.FieldLevel) bool {
	length, err := strconv.Atoi(field.Param())
	if err != nil {
		panic(err)
	}

	value := field.Field().String()
	if !regexNumber.MatchString(value) {
		return false
	}

	return len(value) == length
}

func TestValidationParameter(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("pin", MustValidPin)
	if err != nil {
		return
	}

	type Login struct {
		Name     string `validate:"required"`
		Password string `validate:"required,pin=6"`
	}

	login := Login{
		Name:     "brian",
		Password: "1",
	}

	err = validate.Struct(login)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestOrRule(t *testing.T) {
	validate := validator.New()
	type Login struct {
		Name     string `validate:"required,email|numeric"`
		Password string `validate:"required"`
	}

	login := Login{
		Name:     "brian",
		Password: "1",
	}

	err := validate.Struct(login)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func MustEqualIgnoreCase(field validator.FieldLevel) bool {
	value, _, _, ok := field.GetStructFieldOK2()
	if !ok {
		panic("fields not ok")
	}

	firstValue := strings.ToUpper(field.Field().String())
	secondValue := strings.ToUpper(value.String())

	return firstValue == secondValue
}

func TestCustomValidationCrossField(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("must_equal_ignore_case", MustEqualIgnoreCase)
	if err != nil {
		panic(err)
	}

	type Login struct {
		Name     string `validate:"required,must_equal_ignore_case=Email|must_equal_ignore_case=Phone"`
		Password string `validate:"required"`
		Email    string `validate:"required,email|numeric"`
		Phone    string `validate:"required,numeric"`
	}

	login := Login{
		Name:     "brian@gmail.com",
		Password: "1",
		Email:    "brian@gmail.com",
		Phone:    "1234",
	}

	err = validate.Struct(login)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type RegisterRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required,numeric"`
}

func MustValidRegisterSuccess(level validator.StructLevel) {
	registerRequest := level.Current().Interface().(RegisterRequest)

	if registerRequest.Username == registerRequest.Email || registerRequest.Username == registerRequest.Phone {

	} else {
		level.ReportError(registerRequest.Username, "Username", "Username", "valid_username", "")

	}
}

func TestStructLevelValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(MustValidRegisterSuccess, RegisterRequest{})

	request := RegisterRequest{
		Username: "brian",
		Password: "12345",
		Email:    "brian@gmail.com",
		Phone:    "012345",
	}

	err := validate.Struct(request)
	if err != nil {
		fmt.Println(err.Error())
	}
}
