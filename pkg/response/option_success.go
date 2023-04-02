package response

type (
	withMsg  struct{ msg string }
	withData struct{ data any }
	withMeta struct{ meta any }
)

func (w withMsg) Set(app *appSuccess) { app.Message = w.msg }

func (w withData) Set(app *appSuccess) { app.Data = w.data }

func (w withMeta) Set(app *appSuccess) { app.Meta = w.meta }

// WithMsg option to add given msg to success response as `message` field.
func WithMsg(msg string) SuccessOption {
	return withMsg{msg}
}

// WithData option to add the given data to success response as `data` field.
func WithData(data any) SuccessOption {
	return withData{data}
}

// WithMeta option to add the given meta to success response as `meta` field.
func WithMeta(meta any) SuccessOption {
	return withMeta{meta}
}
