package loglerpart

import (
	"github.com/logiqai/logiqctl/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)


// Mon Jul  5 17:58:59 PDT 2021
// * Created
// * Support pattern-signature extraction and organization
//

// MyError dummy for custom erro:
type MyError struct{}

func (m *MyError) Error() string {
	return "boom"
}

//
// PsListType list current PS tag
//
type PsListType struct {
	DateTime         string            `json:"datetime"`
	PsNodeCnt        int               `json:"pnode_cnt"`
	PsNodeCntComment string            `json:"pnode_cnt_comment"`
	PsListComment    string            `json:"ps_list_comment"`
	PsList           map[string]string `json:"ps_list"`
}

//type PsListType struct {
//	DateTime         string
//	PsNodeCnt        int
//	PsNodeCntComment string
//	PsListComment    string
//	PsList           map[string]string
//}

//var PsList[string]string
//var NewMap[string]

var (

	EnablePsFlag = 1
	Get_ps_from_file_flag = 1 // from ps_base_rule.json (file) or rules.go
	MaxPsCount = 8000
	cfgstatepath string = "/src/bitbucket.org/logiqcloud/logler/cfgstate/"
	once sync.Once
	PsCount = make (map[string]int)
	PsByteCount = make (map[string]int)
	PsByteSqCount = make (map[string]int)
	PsExlog = make (map[string]string)

	PsListAll PsListType
	Ps2Pnode = make(map[string]string)

	LogLineCount = 0
	MaxLogLineCount = 50000
	StopCh = make(chan int)
	PsmodCmd string = ""
	LogiqctlVersion string = "<un-tracked>"
)

// how many log message with a PsId

func InitPaths() {

	gopath, exists := os.LookupEnv("GOPATH")

	if exists {
		// Print the value of the environment variable
		cfgstatepath = gopath + cfgstatepath
		// testlogpath = gopath + testlogpath

	} else {
		fmt.Println("err> GOPATH env not exists")
		os.Exit(1)

	}

}

func DumpsPsInfo() {

	DumpCurrentPsList("ps_list")
	DumpCurrentPsStat("ps_stat")
}

func Init(vv string) {
// initialize all variables here
// PS stat tracking

	LogiqctlVersion = vv

	fin:="ps_base_rule"
	jsonFile, err := os.Open(cfgstatepath + fin + ".json")
	if err != nil {
		jsonFile, err = os.Open("./"+fin + ".json")
		if err!=nil {
			//* fmt.Println("Set PSExt to local-mode")
			Get_ps_from_file_flag=0
		} else {
			fmt.Println("Set localFiles mode")
			jsonFile.Close()
		}
	} else {
		jsonFile.Close()
	}
	// fmt.Println("get_ps_from_file_flag=", get_ps_from_file_flag)

	fin="ps_list"
	jsonFile, err = os.Open(cfgstatepath + fin + ".json")
	if err != nil {
		jsonFile, err = os.Open("./"+fin + ".json")
		if err!=nil {
			// create empty ps_list.json
			f, err1 := os.Create("./"+fin+".json")
			MyChkError(err1, "Create ./"+fin+".json failed", 1)
			//  "{
			//  	"datetime": "2020-06-26 00:23:50.811438416",
			//  	"pnode_cnt": 1,
			//  	"pnode_cnt_comment": "current pnode count.  auto generate.",
			//  	"ps_list_comment": "pattern signature (PS) collection.  auto generate.",
			//  	"ps_list": {
			//  	}
			//  }
			fmt.Fprintf(f, "{\n")
			currentTime := time.Now()
			fmt.Fprintf(f, "    \"datetime\": \"%s\",\n", currentTime.Format("2006-01-02 15:04:05.000000"))
			fmt.Fprintf(f, "    \"pnode_cnt\": 1, \n")
			fmt.Fprintf(f, "    \"pnode_cnt_comment\": \"current pnode count.  auto generated.\",\n")
			fmt.Fprintf(f, "    \"ps_list_comment\": \"pattern signature (PS) collection.  Auto generated.\",\n")
			fmt.Fprintf(f, "    \"ps_list\": {\n")
			fmt.Fprintf(f, "    }\n")
			fmt.Fprintf(f, "}\n")
			f.Close()
		}
	} else {
		jsonFile.Close()
	}
	GetPsListFile(fin)

	// to reinit again
	PsCount = make (map[string]int)
	PsByteCount = make (map[string]int)
	PsByteSqCount = make (map[string]int)
	PsExlog = make (map[string]string)
}

// MyChkError check error and exit
//
func MyChkError(err error,
	mesg string,
	ExitFlag int) {

	if err != nil {
		fmt.Println("ERROR> err chk, ", mesg)
		fmt.Println("       ", err)
		if ExitFlag == 1 {
			os.Exit(1)
		}
	}
}

// MyJSONChkError check error and exit, not working
//
func MyJSONChkError(err error,
	mesg string,
	ExitFlag int) {

	j := ""
	i := 1

	if jsonError, ok := err.(*json.SyntaxError); ok {
		line, character, lcErr := lineAndCharacter(j, int(jsonError.Offset))
		fmt.Fprintf(os.Stderr, "test %d failed with error: Cannot parse JSON schema due to a syntax error at line %d, character %d: %v\n", i+1, line, character, jsonError.Error())
		if lcErr != nil {
			fmt.Fprintf(os.Stderr, "Couldn't find the line and character position of the error due to error %v\n", lcErr)
		}
	}
	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		line, character, lcErr := lineAndCharacter(j, int(jsonError.Offset))
		fmt.Fprintf(os.Stderr, "test %d failed with error: The JSON type '%v' cannot be converted into the Go '%v' type on struct '%s', field '%v'. See input file line %d, character %d\n", i+1, jsonError.Value, jsonError.Type.Name(), jsonError.Struct, jsonError.Field, line, character)
		if lcErr != nil {
			fmt.Fprintf(os.Stderr, "test %d failed with error: Couldn't find the line and character position of the error due to error %v\n", i+1, lcErr)
		}
	}

	if err != nil {
		fmt.Println("ERROR> err chk, ", mesg)
		fmt.Println("       ", err)
		if ExitFlag == 1 {
			os.Exit(1)
		}
	}
}

//
// MyOkCheck takes ok bool and print mesg and exit depending on input
//
func MyOkCheck(ok bool,
	mesg string,
	ExitFlag int) {

	if !ok {
		fmt.Println("ERROR> ok chk, ", mesg)
		if ExitFlag == 1 {
			os.Exit(1)
		}
	}

}

type Pair struct {
	Key string
	Value int
}
type PairList []Pair
func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func DumpCurrentPsStat(fout string) {
	patternFile := fmt.Sprintf("./%s.out",fout)
	file, err := os.Create(patternFile)
	MyChkError(err, fmt.Sprintf("Cannot open %s", patternFile), 1)

	keys := make([]string, 0, len(PsCount))
	totalsum := 0.0
	totalbytesum := 0.0
	for k, _:= range PsCount {
		keys = append(keys, k)
		totalsum += float64(PsCount[k])
		totalbytesum += float64(PsByteCount[k])
	}

	p := make(PairList, len(PsByteCount))

	i := 0
	for k, v := range PsByteCount {
		p[i] = Pair{k, v}
		i+=1
	}
	sort.Sort(sort.Reverse(p))

	sort.Strings(keys)
	//fmt.Println("keys=", keys)
	fmt.Fprintf(file, "Version (logiqctl): %s\n", LogiqctlVersion)
	fmt.Fprintf(file, "Version (psmod): %s\n\n", GetPsmodVersion())
	fmt.Fprintf(file, "Total Logs: %d \n", LogLineCount)
	fmt.Fprintf(file, "Total Pattern-signature Count: %d \n\n", len(PsCount))
	fmt.Fprintf(file, "PsId  LogByteCnt  LogByteCntPrcnt  LogCnt  LogCntPrcnt  Avg(Byte/Line)  ~SD(Byte/Line)\n")
	fmt.Fprintf(file, " Log_Pattern_Signature --< >--\n")
	fmt.Fprintf(file, " First_Log_Example --< >--\n")
	fmt.Fprintf(file, "===================================================================\n")
	for ii,_ := range p {

			k := p[ii].Key

			// fmt.Println("kstr=", kstr, "   k=", k)

			sd := -1.0
		    if PsCount[k]>2 {
				vari := float64(PsByteSqCount[k]-(PsByteCount[k]*PsByteCount[k])/PsCount[k])/float64(PsCount[k]-1)
				sd = vari/2.0
				// fmt.Println("vari=", vari, "   sd=", sd)
			}
			fmt.Fprintf(file, "\n%s  %d  %5.2f%%  %d  %5.2f%%  %7.0f  %7.0f\n",
				k,
				PsByteCount[k], (float64(PsByteCount[k])/totalbytesum)*100.0,
				PsCount[k], (float64(PsCount[k])/totalsum)*100.0,
				float64(PsByteCount[k]/PsCount[k]), sd,
			)
			//fmt.Fprintf(file, " LogPatSign: --<%s>--\n",   strings.TrimSpace(ReplaceAllStringPsList(PsListAll.PsList[k])))
		fmt.Fprintf(file, " LogPatSign: --<%s>--\n",   strings.TrimSpace(PsListAll.PsList[k]))
		fmt.Fprintf(file, " LogExample: --<%s>--\n",   strings.TrimSpace(PsExlog[k]))
	}
	//fmt.Fprintf(file, "Total Logs: %d\n", int(totalsum))
	//fmt.Printf("Pattern signatures generated at %s\n", patternFile)

}

// DumpCurrentPsList Open JSON and read base rules
//
func DumpCurrentPsList(fout string) {

	currentTime := time.Now()
	// fmt.Println("Time with NanoSeconds: ",
	//      currentTime.Format("2006-01-02 15:04:05.000000000"))
	PsListAll.DateTime = currentTime.Format("2006-01-02 15:04:05.000000000")

	b, err := json.MarshalIndent(&PsListAll, "", "  ")
	MyChkError(err, "DumpCurrentPsList DumpCurrentPsList MarshalIndent", 1)

	b = EscOff(b)

	file, err := os.Create("./" + fout + ".out")
	MyChkError(err, "DumpCurrentPsList Cannot open ./"+fout+".out", 1)

	fmt.Fprintf(file, "%+v\n", string(b))

	//fmt.Printf("%+v\n", string(b))
	//os.Stdout.Write(b)

}


//
// GetPsListFile get PsList from file
//
func GetPsListFile(fin string) {

	//var m map[string]json.RawMessage
	//var PsList map[string]string

	// Open our jsonFile
	jsonFile, err := os.Open(cfgstatepath + fin + ".json")
	if err!=nil {
		// fmt.Println("Not open " + cfgstatepath + fin + ".json")
		jsonFile, err = os.Open("./" + fin + ".json")
		MyChkError(err, "Cannot open "+cfgstatepath+fin+".json", 1)
	}

	//* fmt.Println("Successfully Opened ", cfgstatepath+fin+".json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	MyChkError(err, "JSON reading error", 1)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	//err = json.Unmarshal(byteValue, &m)
	err = json.Unmarshal(byteValue, &PsListAll)
	MyChkError(err, "JSON Unmarshal error", 1)

	//PsListRaw, ok := m["ps_list"]
	//MyOkCheck(ok, "expect to find 'ps_list'", 1)
	//err = json.Unmarshal(PsListRaw, &PsList)
	//MyChkError(err, "JSON PsList Unmarshal error", 1)

	for pnode, ps := range PsListAll.PsList {
		//fmt.Printf("%s -> %s\n", pnode, ps)
		Ps2Pnode[ps] = pnode
	}

}

// PsCheckAndReturnTag takes a LogPs returns pnode tag.
// if PsItem is not in the list, create new entry and update list
func PsCheckAndReturnTag(PsItem string, mesg string) string {

	//h := sha256.New()
	//h.Write([]byte(PsItem))
	//PsHash := h.Sum(nil)

	//fmt.Println("PsItem=", PsItem)
	//fmt.Printf("   PsHash=%x\n", PsHash)

	pnode, ok := Ps2Pnode[PsItem]
	if ok {
		PsCount[pnode] += 1
		PsByteCount[pnode] += len(mesg)
		PsByteSqCount[pnode] += len(mesg)*len(mesg)

		return (pnode)
	} else {
		// Postgres Table Update here
		// Retreive or update from Postgres Global
		// Temp below

		PsNodeCntStr := ""
		if PsListAll.PsNodeCnt>=MaxPsCount {
			PsNodeCntStr = "MaxPs"
		} else {
			PsNodeCntStr = "p" + strconv.Itoa((PsListAll).PsNodeCnt)
			Ps2Pnode[PsItem] = PsNodeCntStr
			PsListAll.PsList[PsNodeCntStr] = PsItem

			PsListAll.PsNodeCnt += 1
			PsExlog[PsNodeCntStr] = mesg

		}

		PsCount[PsNodeCntStr] += 1
		PsByteCount[PsNodeCntStr] = len(mesg)
		PsByteSqCount[PsNodeCntStr] = len(mesg)*len(mesg)

		return (PsNodeCntStr)
	}

}

// for debugging, not working
//func lineAndCharacter
func lineAndCharacter(input string, offset int) (line int, character int, err error) {
        lf := rune(0x0A)

        if offset > len(input) || offset < 0 {
                return 0, 0, fmt.Errorf("Couldn't find offset %d within the input.", offset)
        }

        // Humans tend to count from 1.
        line = 1

        for i, b := range input {
                if b == lf {
                        line++
                        character = 0
                }
                character++
                if i == offset {
                        break
                }
        }

        return line, character, nil
}


func IncLogLineCount() {

	LogLineCount+=1

	if (LogLineCount>=MaxLogLineCount) {
		 StopCh <- 1
	}
}

func SetupCloseHandler() {

// MaxLogLine reached

	go func() {
		<-StopCh
		fmt.Println("\r- Max log lines: ", MaxLogLineCount, " reached!  exit")
		DumpCurrentPsStat("ps_stat")
		os.Exit(0)
	}()

	c := make(chan os.Signal)
// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// capture any exit signas
	signal.Notify(c, os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		DumpCurrentPsStat("ps_stat")
		os.Exit(0)
	}()
}

//EscOff convert escape char to char
func EscOff(b []byte) []byte {

	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	return b
}

func ProcessLogCmd(pp string) string {

	// skip quotes
	//	regexp.MustCompile("^\\'")

	//fmt.Println("in ProcessLogCmd: ", pp)
	out, err := exec.Command(PsmodCmd, pp).Output()
	if err != nil {
		utils.HandleError2(err, fmt.Sprintf("psmod=<%s> executable returns error ", PsmodCmd))
	}
	// fmt.Printf(" pp=<%s>\n", pp)
	// fmt.Printf("out=<%s>\n", out)
	return (string(out))
}

func GetPsmodVersion() string {

	//fmt.Println("in ProcessLogCmd: ", pp)
	out, err := exec.Command(PsmodCmd, "Gget", "Vver").Output()
	if err != nil {
		fmt.Println("error seen")
		os.Exit(1)
		// log.Fatal(err)
	}
	return (string(out))
}

func CheckPsmod() {

	_, err := exec.Command("ls", "./psmod").Output()
	if (err == nil) {
		PsmodCmd = "./psmod"
		return
	}

	_, err = exec.Command("which", "psmod").Output()
	if err == nil {
		PsmodCmd = "psmod"
		return
	}

	// for windows system
	_, err = exec.Command("cmd", "/c", "dir", "psmod.exe").CombinedOutput()
	if (err == nil) {
		PsmodCmd = "psmod.exe"
		return
	}

	fmt.Println("Enter 'psmod' location-name: <path/name>:")
	_, err = fmt.Scanln(&PsmodCmd)
	if (err != nil) {
		utils.HandleError(err)
	}
	return
	/*
	_, err = exec.Command("ls", PsmodCmd).Output()
	if err!=nil {
		errmsg := errors.New(fmt.Sprintf("Pattern-signature generation requires PSMOD add-on executable.  File '%s' not found.", PsmodCmd))
		utils.HandleError(errmsg)
	}
	*/

}






