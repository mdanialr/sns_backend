package response

import "github.com/go-playground/validator/v10"

type (
	withErr       struct{ err error }
	withErrMsg    struct{ msg string }
	withErrCode   struct{ code string }
	withErrDetail struct{ detail any }
	withErrValid  struct{ valid validator.ValidationErrors }
)

func (w withErr) Set(app *appError) { app.Message = w.err.Error() }

func (w withErrMsg) Set(app *appError) { app.Message = w.msg }

func (w withErrCode) Set(app *appError) { app.Code = w.code }

func (w withErrDetail) Set(app *appError) { app.Detail = w.detail }

func (w withErrValid) Set(app *appError) { app.Detail = NewValidationErrors(w.valid) }

// WithErr option to add given error message to error response as `message`
// field.
func WithErr(err error) AppErrorOption {
	return withErr{err}
}

// WithErrMsg option to add given message to error response as `message` field.
func WithErrMsg(msg string) AppErrorOption {
	return withErrMsg{msg}
}

// WithErrCode option to add given code to error response as `code` field.
func WithErrCode(code string) AppErrorOption {
	return withErrCode{code}
}

// WithErrDetail option to add given detail to error response as `detail`
// field.
func WithErrDetail(detail any) AppErrorOption {
	return withErrDetail{detail}
}

// WithErrValidation option to add given valid to error response as `detail`
// field.
func WithErrValidation(valid validator.ValidationErrors) AppErrorOption {
	return withErrValid{valid}
}
