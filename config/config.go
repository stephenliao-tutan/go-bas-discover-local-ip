package config

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"os"
	"path"
	"sync"
)

const (
	BASDiscover_HomeDir      = ".basdis"
	BASDiscover_CFG_FileName = "basdis.json"
)

type BasDiscover struct {
	MgtHttpPort    int      `json:"mgthttpport"`
	KeyPath        string   `json:"keypath"`
	CmdListenPort  string   `json:"cmdlistenport"`

}

var (
	bascfgInst     *BasDiscover
	bascfgInstLock sync.Mutex
)

func (bc *BasDiscover) InitCfg() *BasDiscover {
	bc.MgtHttpPort = 50818
	bc.KeyPath = "/keystore"
	bc.CmdListenPort = "127.0.0.1:59527"

	return bc
}

func (bc *BasDiscover) Load() *BasDiscover {
	if !tools.FileExists(GetBASDisCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetBASDisCFGFile())
	if err != nil {
		log.Println("load file failed", err)
		return nil
	}

	//bc1:=&BASDConfig{}

	err = json.Unmarshal(jbytes, bc)
	if err != nil {
		log.Println("load configuration unmarshal failed", err)
		return nil
	}

	return bc

}

func newBasDisCfg() *BasDiscover {

	bc := &BasDiscover{}

	bc.InitCfg()

	return bc
}

func GetBasDisCfg() *BasDiscover {
	if bascfgInst == nil {
		bascfgInstLock.Lock()
		defer bascfgInstLock.Unlock()
		if bascfgInst == nil {
			bascfgInst = newBasDisCfg()
		}
	}

	return bascfgInst
}

func PreLoad() *BasDiscover {
	bc := &BasDiscover{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *BasDiscover {
	bc := &BasDiscover{}

	bc.InitCfg()

	bcontent, err := tools.OpenAndReadAll(file)
	if err != nil {
		log.Fatal("Load Config file failed")
		return nil
	}

	err = json.Unmarshal(bcontent, bc)
	if err != nil {
		log.Fatal("Load Config From json failed")
		return nil
	}

	bascfgInstLock.Lock()
	defer bascfgInstLock.Unlock()
	bascfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *BasDiscover) *BasDiscover) *BasDiscover {
	bascfgInstLock.Lock()
	defer bascfgInstLock.Unlock()

	lbc := newBasDisCfg().Load()

	if lbc != nil {
		bascfgInst = lbc
	} else {
		lbc = newBasDisCfg()
	}

	bascfgInst = initfromcmd(lbc)

	return bascfgInst
}

func GetBASDisHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, BASDiscover_HomeDir)
}

func GetBASDisCFGFile() string {
	return path.Join(GetBASDisHomeDir(), BASDiscover_CFG_FileName)
}


func (bc *BasDiscover) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetBASDisHomeDir()) {
		os.MkdirAll(GetBASDisHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetBASDisCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func IsInitialized() bool {
	if tools.FileExists(GetBASDisCFGFile()) {
		return true
	}

	return false
}
