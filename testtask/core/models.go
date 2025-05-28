package core

type Person struct {
	ID          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" db:"people_id"`
	Name        string  `json:"name" example:"Dmitry" db:"name"`
	Surname     string  `json:"surname" example:"Ushakov" db:"surname"`
	Patronymic  *string `json:"patronymic,omitempty" example:"Vasilevich" db:"patronymic"`
	Age         int     `json:"age" example:"40" db:"age"`
	Gender      string  `json:"gender" example:"male" db:"gender"`
	Nationality string  `json:"nationality" example:"RU" db:"nationality"`
}

type APIAgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type APIGenderResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type APINationResponse struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

type PersonRequest struct {
	Name       string  `json:"name" example:"Dmitry"`
	Surname    string  `json:"surname" example:"Ushakov"`
	Patronymic *string `json:"patronymic,omitempty" example:"Vasilevich"`
}

type PersonFilters struct {
	Age         string `json:"age" example:"30"`
	Gender      string `json:"gender" example:"male"`
	Nationality string `json:"nationality" example:"RU"`
	Limit       string `json:"limit" example:"10"`
	Offset      string `json:"offset" example:"0"`
}
