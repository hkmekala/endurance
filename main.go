package main

//Importing Standard Packages
import (
	"encoding/binary"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"

	"codeelite.com/controller"
)

func main() {
	fmt.Print("Started Server at 8080")
	//http.Handle("/", http.FileServer(http.Dir("node_modules")))
	//router := httprouter.New()
	http.HandleFunc("/", showIndex)
	http.HandleFunc("/executecode", executecode)
	//resourcesRouter()
	//router.GET("/", showIndex)
	//router.GET("/", http.FileServer(http.Dir(".")))
	//router.ServeFiles("./node_modules", http.Dir("./node_modules"))
	//router.POST("/executecode", executecode)
	http.ListenAndServe(":8080", nil)

}

/*func FileServe(rw http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("views"))
}*/
const maxInt32 = 1<<(32-1) - 1

func writeLen(b []byte, l int) []byte {
	if 0 > l || l > maxInt32 {
		panic("writeLen: invalid length")
	}
	var lb [4]byte
	binary.BigEndian.PutUint32(lb[:], uint32(l))
	return append(b, lb[:]...)
}

func readLen(b []byte) ([]byte, int) {
	if len(b) < 4 {
		panic("readLen: invalid length")
	}
	l := binary.BigEndian.Uint32(b)
	if l > maxInt32 {
		panic("readLen: invalid length")
	}
	return b[4:], int(l)
}
func Encode(s []string) []byte {
	var b []byte
	b = writeLen(b, len(s))
	for _, ss := range s {
		b = writeLen(b, len(ss))
		b = append(b, ss...)
	}
	return b
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type OutCode struct {
	Code   string
	Output string
}

func executecode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code_str := r.Form["scode"]
	code := []byte(code_str[0])
	input_str := r.Form["scode_input"]
	input := []byte(input_str[0])
	fmt.Println(input_str)
	fmt.Println(code_str)
	randfolder := RandStringBytes(12)
	path := "./controller/vol/" + randfolder

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	_, copyHostCfile := exec.Command("/bin/bash", "-c", "cp ./controller/vol/main.c "+path).Output()
	if copyHostCfile != nil {
		panic(copyHostCfile)
	}
	_, copyHostShellfile := exec.Command("/bin/bash", "-c", "cp ./controller/vol/compile.sh "+path).Output()
	if copyHostShellfile != nil {
		panic(copyHostShellfile)
	}
	_, copyHostInputfile := exec.Command("/bin/bash", "-c", "cp ./controller/vol/input.txt "+path).Output()
	if copyHostInputfile != nil {
		panic(copyHostInputfile)
	}
	fmt.Println("Created Host ENV")
	err := ioutil.WriteFile("./controller/vol/"+randfolder+"/main.c", code, 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./controller/vol/"+randfolder+"/input.txt", input, 0777)
	if err != nil {
		panic(err)
	}
	runner.Runcode(path, randfolder)

	//time.Sleep(time.Second * 10)
	output, err := ioutil.ReadFile("./controller/vol/" + randfolder + "/data.txt")
	if err != nil {
		panic(err)
	}
	output_errors, err := ioutil.ReadFile("./controller/vol/" + randfolder + "/errors.txt")
	if err != nil {
		panic(err)
	}
	templ_output := string(output) + string(output_errors)
	//fmt.Fprintf(w, "%s", string(code))
	//fmt.Fprintf(w, "%s", ps.ByName("code"))
	templ, err := template.ParseFiles("views/index.html")

	if err != nil {
		panic(err)
	}
	templ_output_obj := OutCode{
		Code:   code_str[0],
		Output: templ_output,
	}
	err = templ.Execute(w, templ_output_obj)
}
func showIndex(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("views/index.html")
	//static_html, err := ioutil.ReadFile("views/index.html")

	if err != nil {
		panic(err)
	}

	err = templ.Execute(w, &OutCode{})
	if err != nil {
		panic(err)
	}
	//fmt.Fprintf(w, "%s", static_html)

}

/*func resourcesRouter() {
	searchDir := "./node_modules"

	fileList := []string{}
	_ = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	//router := httprouter.New()

	for _, file := range fileList[4:] {

		fmt.Println(file)

		//router.ServeFiles("./"+file, http.Dir(searchDir))
	}
	fmt.Print("------------------------------------------------------------Statin Glb")
	arr, err := filepath.Glob("node_modules/*")
	if err != nil {
		panic(err)
	}
	for _, data := range arr {
		fmt.Println(data)
		router.ServeFiles("./"+data, http.Dir(searchDir))
	}

}*/
