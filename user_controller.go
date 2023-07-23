package main

type userController struct {
	userService *userService
}

func newUserController(userService *userService) userController {
	return userController{userService: userService}
}

func (uc *userController) createUser(wAndR writerAndRequest) {
	userId, err := uc.userService.createUser()

	if err != nil {
		wAndR.w.WriteHeader(500)
		return
	}

	wAndR.w.WriteHeader(200)
	wAndR.w.Write([]byte(userId))
}
