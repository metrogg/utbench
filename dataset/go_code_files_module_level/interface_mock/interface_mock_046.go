func NewPetstoreAPI(spec *loads.Document) *PetstoreAPI {
	return &PetstoreAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		UrlformConsumer:     runtime.DiscardConsumer,
		XMLConsumer:         runtime.XMLConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		XMLProducer:         runtime.XMLProducer(),
		PetAddPetHandler: pet.AddPetHandlerFunc(func(params pet.AddPetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation PetAddPet has not yet been implemented")
		}),
		UserCreateUserHandler: user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation UserCreateUser has not yet been implemented")
		}),
		UserCreateUsersWithArrayInputHandler: user.CreateUsersWithArrayInputHandlerFunc(func(params user.CreateUsersWithArrayInputParams) middleware.Responder {
			return middleware.NotImplemented("operation UserCreateUsersWithArrayInput has not yet been implemented")
		}),
		UserCreateUsersWithListInputHandler: user.CreateUsersWithListInputHandlerFunc(func(params user.CreateUsersWithListInputParams) middleware.Responder {
			return middleware.NotImplemented("operation UserCreateUsersWithListInput has not yet been implemented")
		}),
		StoreDeleteOrderHandler: store.DeleteOrderHandlerFunc(func(params store.DeleteOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation StoreDeleteOrder has not yet been implemented")
		}),
		PetDeletePetHandler: pet.DeletePetHandlerFunc(func(params pet.DeletePetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation PetDeletePet has not yet been implemented")
		}),
		UserDeleteUserHandler: user.DeleteUserHandlerFunc(func(params user.DeleteUserParams) middleware.Responder {
			return middleware.NotImplemented("operation UserDeleteUser has not yet been implemented")
		}),
		PetFindPetsByStatusHandler: pet.FindPetsByStatusHandlerFunc(func(params pet.FindPetsByStatusParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation PetFindPetsByStatus has not yet been implemented")
		}),
		PetFindPetsByTagsHandler: pet.FindPetsByTagsHandlerFunc(func(params pet.FindPetsByTagsParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation PetFindPetsByTags has not yet been implemented")
		}),
		StoreGetOrderByIDHandler: store.GetOrderByIDHandlerFunc(func(params store.GetOrderByIDParams) middleware.Responder {
			return middleware.NotImplemented("operation StoreGetOrderByID has not yet been implemented")
		}),
		PetGetPetByIDHandler: pet.GetPetByIDHandlerFunc(func(params pet.GetPetByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation PetGetPetByID has not yet been implemented")
		}),
		UserGetUserByNameHandler: user.GetUserByNameHandlerFunc(func(params user.GetUserByNameParams) middleware.Responder {
			return middleware.NotImplemented("operation UserGetUserByName has not yet been implemented")
		}),
		UserLoginUserHandler: user.LoginUserHandlerFunc(func(params user.LoginUserParams) middleware.Responder {
			return middleware.NotImplemented("operation UserLoginUser has not yet been implemented")
		}),
		UserLogoutUserHandler: user.LogoutUserHandlerFunc(func(params user.LogoutUserParams) middleware.Responder {
			return middleware.NotImplemented("operation UserLogoutUser has not yet been implemented")
		}),
		StorePlaceOrderHandler: store.PlaceOrderHandlerFunc(func(params store.PlaceOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation StorePlaceOrder has not yet been implemented")
		}),
		PetUpdatePetHandler: pet.UpdatePetHandlerFunc(func(params pet.UpdatePetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation PetUpdatePet has not yet been implemented")
		}),
		PetUpdatePetWithFormHandler: pet.UpdatePetWithFormHandlerFunc(func(params pet.UpdatePetWithFormParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation PetUpdatePetWithForm has not yet been implemented")
		}),
		UserUpdateUserHandler: user.UpdateUserHandlerFunc(func(params user.UpdateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation UserUpdateUser has not yet been implemented")
		}),

		// Applies when the "api_key" header is set
		APIKeyAuth: func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (api_key) api_key from header param [api_key] has not yet been implemented")
		},
		PetstoreAuthAuth: func(token string, scopes []string) (interface{}, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (petstore_auth) has not yet been implemented")
		},

		// default authorizer is authorized meaning no requests are blocked
		APIAuthorizer: security.Authorized(),
	}
}
