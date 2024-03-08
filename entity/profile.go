package entity

type UserProfile struct {
	Id          string `db:"id"`
	FullName    string `db:"full_name"`
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
}

type ProfileRegisterRequest struct {
	FullName    string `validate:"required,gte=3,lte=60,alpha"`
	PhoneNumber string `validate:"required,e164,startswith=+62"`
	Password    string `validate:"required,gte=3,lte=64,anyAlphaCapital,anyNumeric,anySpecialChar"`
}

type ProfileRegisterResponse struct {
	Id string
}

type GetProfileRequest struct {
	ProfileId string
}

type GetProfileResponse struct {
	FullName    string `validate:"required,gte=3,lte=60,alpha"`
	PhoneNumber string `validate:"required,e164,startswith=+62"`
}

type LoginRequest struct {
	PhoneNumber string `validate:"required,e164,startswith=+62"`
	Password    string // no need to validate password on login
}

type LoginResponse struct {
	Token string
}

type UpdateProfileRequest struct {
	Id          string
	FullName    string `validate:"required,gte=3,lte=60,alpha"`
	PhoneNumber string `validate:"required,gte=3,lte=64,anyAlphaCapital,anyNumeric,anySpecialChar"`
}
