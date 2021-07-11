package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/hojunin/hjcoin/blockchain"
)


const( 
	port string = ":4000"
	templateDir string = "explorer/templates/" 
)
var templates *template.Template
type homeData struct{
	PageTitle string
	Blocks []*blockchain.Block
}

// * /home을 제어하는 핸들러
func home (rw http.ResponseWriter, r *http.Request)  {
	data := homeData{"Home", blockchain.GetBlockchain().AllBlock()}
	templates.ExecuteTemplate(rw, "home", data)
}

// * /add를 제어하는 헨들러
func add(rw http.ResponseWriter, r *http.Request)  {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}



func Start()  {
	// * 미리 작성해둔 템플릿을 로드합니다. 
	templates = template.Must(template.ParseGlob(templateDir+"pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir+"partials/*.gohtml"))
	// * 핸들러를 호출합니다.
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	// * 서버를 가동합니다.
	log.Fatal(http.ListenAndServe(port,nil))
}