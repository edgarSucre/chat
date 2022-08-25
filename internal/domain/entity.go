package domain

type (
	UserParam struct {
		UserName string `validate:"required,gt=5"`
		Password string `validate:"required,gt=5"`
	}

	UserResponse struct {
		ID       int64
		UserName string
		Password string
	}

	RoomParam struct {
		Name string `validate:"required,gt=7"`
	}

	RoomResponse struct {
		ID   int64
		Name string
	}
)
