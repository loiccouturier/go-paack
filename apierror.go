package paack

type ApiErrorDesc struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorKey         string `json:"error_key"`
}

type ApiError struct {
	Err        error          `json:"err"`
	StatusCode int            `json:"status"`
	Errors     []ApiErrorDesc `json:"errors"`
}
