package common

func TryError(callBack func(), errorFunc func(error)) {
	func() {
		defer func()  {
			if re := recover(); re != nil {
				errorFunc(re.(error))
			}
		}()
		callBack()
	}()
}