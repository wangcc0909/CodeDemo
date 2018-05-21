package mock

type Retrieve struct {
	Contents string
}

func (r *Retrieve) Post(url string, form map[string]string) string {
	r.Contents = form["contents"]
	return "ok"
}

func (r *Retrieve) Get(url string) string {

	return r.Contents
}






