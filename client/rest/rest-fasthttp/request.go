package fasthttp

type requestOptions struct {
	headers     map[string]string
	queryParams map[string]string
}

type RequestOption func(*requestOptions)

func WithRequestHeader(key, value string) RequestOption {
	return func(r *requestOptions) {
		if r.headers == nil {
			r.headers = make(map[string]string)
		}
		r.headers[key] = value
	}
}

func WithQueryParam(key, value string) RequestOption {
	return func(r *requestOptions) {
		if r.queryParams == nil {
			r.queryParams = make(map[string]string)
		}
		r.queryParams[key] = value
	}
}
