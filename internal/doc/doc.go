package doc

type Document struct {
	Author string `fake:"{firstname}" json:"author"`
	Title  string `fake:"{sentence:3}" json:"title"`
	Text   string `fake:"{paragraph:2,2,10, }" json:"text"`
}
