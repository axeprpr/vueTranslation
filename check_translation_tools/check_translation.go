package main
import (
    "fmt"
    "flag"
    "log"
    "os"
    "crypto/md5"
    "encoding/hex"
    "net/http"
    "path/filepath"
    "strings"
    "io/ioutil"
    "time"
    "net/url"
    "regexp"
    "strconv"
	"github.com/json-iterator/go"
)

// 百度翻译
var appID = 20180821000197188
var password = "VC_KuXBgnmJghlMsm5Aw"
var Url = "http://api.fanyi.baidu.com/api/trans/vip/translate"

func get_current_path() string {
    path, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
    }
    return path
}

func file_walk(searchPath string, searchFile string) []string {
    fileList := []string{}
    err := filepath.Walk(searchPath, func(path string, f os.FileInfo, err error) error {
        if f.IsDir() {
            return nil
        }
        if strings.HasSuffix(path, "."+searchFile){
            fileList = append(fileList, path)
        }
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    return fileList
}

func a_in_b(a string, b []string) bool {
    for _, i := range b {
        if i == a {
            return false
        }
    }
    return true
}

type TranslateModel struct{
    Q string
    From string
    To string
    Appid int
    Salt int
    Sign string
}
func NewTranslateModeler(q, from, to string) TranslateModel{
    tran := TranslateModel{
        Q: q,
        From: from,
        To: to,
    }
    tran.Appid = appID
    tran.Salt = time.Now().Second()
    content := strconv.Itoa(appID) + q + strconv.Itoa(tran.Salt) + password
    sign := SumString(content)//计算sign值
    tran.Sign = sign
    return tran
}

func (tran TranslateModel) ToValues() url.Values{
    values := url.Values{
        "q": {tran.Q},
        "from": {tran.From},
        "to": {tran.To},
        "appid":{strconv.Itoa(tran.Appid)},
        "salt": {strconv.Itoa(tran.Salt)},
        "sign": {tran.Sign},
    }
    return values
}

func SumString(content string) string{
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(content))
    bys := md5Ctx.Sum(nil)
    value := hex.EncodeToString(bys)
    return value
}

func trans_en2zh(content string) string{
    tran := NewTranslateModeler(content, "en", "zh")
    values := tran.ToValues()
    resp, err := http.PostForm(Url, values)
    if err != nil{
        fmt.Println(err)
        return "wrong translation"
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil{
        fmt.Println(err)
        return "wrong translation"
    } 
    res := jsoniter.Get(body, "trans_result", 0,"dst").ToString()
    reg := regexp.MustCompile(`‘|’|“|”`)
    res = reg.ReplaceAllString(res, "\"")
    return res
}

func main() {
    currentPath := get_current_path()
    path := flag.String("p", currentPath, "Path to search. Default is pwd.")
    searchFile := flag.String("f", "vue", "File format to search.")
    matchFile := flag.String("m", "zh_en.js", "File to match with json string in it.")
    showFilePath := flag.String("s", "false", "show file path.")
    translate := flag.String("t", "false", "auto translate.")
    flag.Parse()
    
    enList := []string{}
    pathList := []string{}
    fileList := file_walk(*path, *searchFile)
    for _, file := range fileList {
        content, err := ioutil.ReadFile(file)
        if err != nil {
            fmt.Println(err)
        }
        contentStr := string(content)
        r, _ := regexp.Compile(`_\(([^()]|(\?R))*\)`)
        out := r.FindAllStringSubmatch(contentStr, -1)
        for _, i := range out {
            enstr := strings.Replace(i[0], "_(", "", -1)
            enstr = strings.Replace(enstr,")","",-1)
            r, _  := regexp.Compile(`,[ ]*{.*}`)
            enstr = r.ReplaceAllString(enstr, ``)
            ok := a_in_b(enstr, enList)
            if ok {
                enList = append(enList, enstr)
                pathList = append(pathList, file)
            }
        }
    }
    
    if *matchFile != "zh_en.js" {
        matchFileContent, err := ioutil.ReadFile(*matchFile)
        matchFileContentStr := string(matchFileContent)
        if err != nil {
            fmt.Println(err)
        }
        for index, i := range enList{
            if !strings.Contains(matchFileContentStr, i) {
                if *showFilePath == "true" {
                    res := "\"在此输入翻译内容\""
                    if *translate == "true"{
                        res = trans_en2zh(i)
                    }
                    fmt.Println(pathList[index])
                    fmt.Println(i+": " + res + ",")
                }else{
                    res := "\"在此输入翻译内容\""
                    if *translate == "true"{
                        res = trans_en2zh(i)
                    }
                    fmt.Println(i+": " + res + ",")
                }
            }
        }
    }
}
