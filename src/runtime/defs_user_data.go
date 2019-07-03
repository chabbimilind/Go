// created by cgo -cdefs and then converted to Go
// cgo -cdefs defs_linux.go defs1_linux.go

package runtime

type PMUEvent struct {
    Cat       uint32 // don't use type, which is a keyword reserved in golang
    Code      uint64
    Period    uint64
    PreciseIP uint8
}
