package debugger

/*
 * go pprof.
 * author: CC
 * email : 151503324@qq.com
 * date  : 2021.08.06
 */

import (
	"os"
	"runtime/pprof"
	"sync"
)

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

type Pprof struct {
	sync.Mutex
	profile *string
	fp      *os.File
}

/************************************************************************/
// export functions.
/************************************************************************/

// open pprof.
func (this *Pprof) Open(fname string) {
	// this.profile = flag.String("cpuprofile", "", "write cpu profile")
	// flag.Parse()
	// f, err := os.Create(*this.profile)
	f, err := os.Create(fname)
	if err != nil {
		LogError("os.Create:", err)
		return
	}
	this.fp = f
	pprof.StartCPUProfile(this.fp)
}

// close pprof.
func (this *Pprof) Close() {
	// defer os.Exit(1)
	pprof.StopCPUProfile()
	this.fp.Close()
	LogInfo("pprof.go - Closed")
}

/************************************************************************/
// moudule functions.
/************************************************************************/

/************************************************************************/
// unit tests.
/************************************************************************/
