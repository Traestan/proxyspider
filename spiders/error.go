package spiders

//Error - структура ошибки
type Error struct {
	msg  string
	code int
}

//NewOVCServiceError создает новую ошибку
func ProxySpiderError(msg string) *Error {
	return &Error{
		msg: msg,
	}
}

//Error наслудует интерфейс error
func (e *Error) Error() string {
	return e.msg
}
