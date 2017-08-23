package core

type Result map[string]interface{}

func MakeResult(ok bool, err string) (r Result) {
	r = make(Result)
	r.SetErr(err)
	r.SetOk(ok)
	return
}

func NewResult() (r Result) {
	r = make(Result)
	r["Err"], r["Ok"] = "", false
	return
}

func (r Result) SetOk(ok bool) {
	r["Ok"] = ok
}

func (r Result) SetErr(err string) {
	r["Err"] = err
}

func (r Result) Set(key string, value interface{}) {
	//not use reflect
	if key == "" {
		return
	}
	r[key] = value
}
