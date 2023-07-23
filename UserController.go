package main

type UserController struct {
	userService *UserService
}

func NewUserController(userService *UserService) UserController {
	return UserController{userService: userService}
}

func (uc *UserController) createUser(wAndR WriterAndRequest) {
	userId, err := uc.userService.createUser()

	if err != nil {
		wAndR.w.WriteHeader(500)
		return
	}

	wAndR.w.WriteHeader(200)
	wAndR.w.Write([]byte(userId))
}
