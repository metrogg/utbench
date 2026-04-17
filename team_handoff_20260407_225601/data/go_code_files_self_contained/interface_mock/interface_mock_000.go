package main

func init() {
	AddInitEventHandler(func(typeOf Event, value interface{}) (responseOf EventResponse) {
		if typeOf == REVEL_BEFORE_MODULES_LOADED {
			RegisterServerEngine(GO_NATIVE_SERVER_ENGINE, func() ServerEngine { return &GoHttpServer{} })
		}
		return
	})
}
