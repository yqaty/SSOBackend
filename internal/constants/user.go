package constants

type Gender int

const (
	GenderMale   Gender = 1
	GenderFemale Gender = 2
	GenderOther  Gender = 3
)

type LoginType string

const (
	LoginTypeEmail         LoginType = "email"
	LoginTypePhoneSMS      LoginType = "sms"
	LoginTypePhonePassword LoginType = "phone"
	LoginTypeLarkOauth     LoginType = "lark"
)
