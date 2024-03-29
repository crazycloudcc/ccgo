package errcode

import "github.com/crazycloudcc/ccgo/datastructs"

/*
 * [file desc]
 * author: CC
 * email : 151503324@qq.com
 * date  : 2017.06.17
 */

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

const (
	Success int32 = 0
	Unknown int32 = 1000 + iota
)

var errfmt *datastructs.Hash

/************************************************************************/
// export functions.
/************************************************************************/

/************************************************************************/
// moudule functions.
/************************************************************************/

func init() {
	errfmt = datastructs.NewHash(1)
	errfmt.Add(Success, "success.")
	errfmt.Add(Unknown, "error unknown.")
}

/************************************************************************/
// unit tests.
/************************************************************************/
