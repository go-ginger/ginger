package validate

var Validators map[string]IValidator

func AddNewValidator(validator IValidator) {
	validator.Initialize()
	Validators[*validator.GetKey()] = validator
}

func init() {
	Validators = make(map[string]IValidator)
	AddNewValidator(new(Required))

	for _, validator := range Validators {
		validator.Initialize()
	}
}
