package parsers

import "google.golang.org/protobuf/proto"

/*
 * ProtoBuf解析器
 * author: CC
 * email : 151503324@qq.com
 * date  : 2021.08.06
 */

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

/************************************************************************/
// export functions.
/************************************************************************/

// Unmarshal TODO
func Unmarshal(buf []byte, pb proto.Message) error {
	return proto.Unmarshal(buf, pb)
}

// Marshal TODO
func Marshal(pb proto.Message) ([]byte, error) {
	return proto.Marshal(pb)
}

/************************************************************************/
// moudule functions.
/************************************************************************/

/************************************************************************/
// unit tests.
/************************************************************************/
