package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/quant61/stack_start_demo/internal"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var procReadProcessMemory = syscall.MustLoadDLL("kernel32.dll").MustFindProc("ReadProcessMemory")

var getThreadContext = windows.MustLoadDLL("kernel32.dll").MustFindProc("GetThreadContext")

func isErrOk(err error) bool {
	return err == nil || err == syscall.Errno(0)
}

func utf16ToString(exeFile []uint16) string {
	for i, v := range exeFile {
		if v == 0 {
			return string(utf16.Decode(exeFile[:i]))
		}
	}
	return ""
}


const (
	CREATE_SUSPENDED = 0x00000004
	DEBUG_PROCESS    = 0x00000001
)

// from https://github.com/duarten/Threadjack/blob/master/WinNT.h
const (
	CONTEXT_CONTROL = 0x1
	CONTEXT_INTEGER = 0x2
	CONTEXT_SEGMENTS = 0x4
	CONTEXT_FLOATING_POINT = 0x8
	CONTEXT_DEBUG_REGISTERS = 0x10
	CONTEXT_XSTATE = 0x20
    CONTEXT_EXCEPTION_ACTIVE = 0x8000000
    CONTEXT_SERVICE_ACTIVE = 0x10000000
    CONTEXT_EXCEPTION_REQUEST = 0x40000000
    CONTEXT_EXCEPTION_REPORTING = 0x80000000
)

type ProcMemReader struct {
	handle windows.Handle
	// TO cache or not to cache
	CachedPages map[int64][4096]byte
}

func (r ProcMemReader) ReadAt(p []byte, off int64) (n int, err error) {
	var _n uintptr
	err = windows.ReadProcessMemory(r.handle, uintptr(off), &p[0], uintptr(len(p)), &_n)
	return int(_n), err
}

func (r ProcMemReader) Close() error {
	return windows.CloseHandle(r.handle)
}

func NewProcMemReader(cmd *exec.Cmd) (*ProcMemReader, error) {
	h, err := windows.OpenProcess(windows.PROCESS_VM_READ, false, uint32(cmd.Process.Pid))
	if err != nil {
		return nil, err
	}

	return &ProcMemReader{handle: h}, nil
}


func startProcessW(args startProcessArgs) *exec.Cmd {
	// TODO: replace it with minimal PE
	p, _ := os.Executable()
	cmd := exec.Command(p)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: DEBUG_PROCESS | CREATE_SUSPENDED,
	}
	cmd.Args = []string{"one", "two", "three"}
	//cmd.Path = "./bin"

	cmd.Env = []string{}
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	return cmd
}

func streamThreads(ctx context.Context) (chan windows.ThreadEntry32, error) {
	ch := make(chan windows.ThreadEntry32)
	h, err := windows.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPTHREAD, 0)
	if !isErrOk(err) {
		fmt.Println(err)
		return nil, err
	}
	var e windows.ThreadEntry32
	e.Size = uint32(unsafe.Sizeof(e))
	err = windows.Thread32First(h, &e)
	if err != nil {
		fmt.Println("Thread32First", err)
		windows.CloseHandle(h)
		return nil, err
	}
	go func() {
		defer close(ch)
		defer windows.CloseHandle(h)
		for {
			select {
			case ch <- e:
				_ = 0
				//fmt.Println("yield", e.ThreadID)
			case <-ctx.Done():
				fmt.Println("cancelled")
				return
			}
			err = windows.Thread32Next(h, &e)
			if err != nil {
				fmt.Println("Thread32Next", h, err)
				return
			}
		}
	}()

	return ch, nil
}

func findFirstThreadOf(pid int) (*windows.ThreadEntry32, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ch, err := streamThreads(ctx)
	if err != nil {
		return nil, err
	}
	_ = cancel
	for t := range ch {
		if t.OwnerProcessID == uint32(pid) {
			cancel()
			return &t, nil
		}
	}
	return nil, nil
}

// from https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-context
type ThreadContext struct {
	PHome [6]uint64
	ContextFlags uint32
	MxCsr uint32
	SegCs uint16
	SegDs uint16
	SegEs uint16
	SegFs uint16
	SegGs uint16
	SegSs uint16
	EFlags uint32
	Dr0    uint64
	Dr1    uint64
	Dr2    uint64
	Dr3    uint64
	Dr6    uint64
	Dr7    uint64
	Rax    uint64
	Rcx    uint64
	Rdx    uint64
	Rbx    uint64
	Rsp    uint64
	Rbp    uint64
	Rsi    uint64
	Rdi    uint64
	R8     uint64
	R9     uint64
	R10    uint64
	R11    uint64
	R12    uint64
	R13    uint64
	R14    uint64
	R15    uint64
	Rip    uint64
	// dummy space
	_notImplemented [1024]byte
}

func PrintMaps(cmd *exec.Cmd) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(cmd.Process.Pid))
	if err != nil{
		fmt.Println(err)
		return
	}
	var buf windows.MemoryBasicInformation
	var pos uintptr
	for {
		err = windows.VirtualQueryEx(h, pos, &buf, uintptr(unsafe.Sizeof(buf)))
		if err != nil {
			fmt.Println(err)
			return
		}
		// `text/template`, maybe?

		typStr := map[uint32]string{
			0x1000000: "MEM_IMAGE",
			0x40000: "MEM_MAPPED",
			0x20000: "MEM_PRIVATE",
			0:       "MEM_NULL",
		}[buf.Type]
		if typStr == "" {
			typStr = fmt.Sprintf("typ(0x%x)", buf.Type)
		}
		stateStr := map[uint32]string{
			0x1000: "MEM_COMMIT",
			0x10000: "MEM_FREE",
			0x2000: "MEM_RESERVE",
			0: "MEM_NULLSTATE",
		}[buf.State]
		if stateStr == "" {
			stateStr = fmt.Sprintf("typ(0x%x)", buf.State)
		}
		fmt.Printf("baseAddr=0x%x\tsize=0x%x\tprot=0x%x\t%s\t%s\n",
			buf.BaseAddress, buf.RegionSize, buf.AllocationProtect, typStr, stateStr)
		if pos != buf.BaseAddress {
			panic("Pos != buf.BaseAddress")
		}
		pos += buf.RegionSize
	}

}


func run() {
	cmd := startProcessW(startProcessArgs{})
	mem, err := NewProcMemReader(cmd)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("new pid:", cmd.Process.Pid)
	_ = mem

	th, err := findFirstThreadOf(cmd.Process.Pid)
	if th == nil {
		fmt.Println(err)
		return
	}
	h, err := windows.OpenThread(windows.THREAD_GET_CONTEXT, false, th.ThreadID)
	fmt.Println(h, err)


	var ctx ThreadContext
	ctx.ContextFlags = CONTEXT_CONTROL|CONTEXT_INTEGER|CONTEXT_SEGMENTS|CONTEXT_FLOATING_POINT|CONTEXT_DEBUG_REGISTERS
	r1, r2, err := getThreadContext.Call(uintptr(h), uintptr(unsafe.Pointer(&ctx)))
	fmt.Println(r1, r2, err)

	printNonZeroFields(ctx)
	PrintMaps(cmd)

	// TODO: unhardcode order and binary size
	order := binary.LittleEndian
	reader := &internal.ReaderHelper{
		ByteOrder: order,
		PtrSize: 8,
	}
	reader.Pos = int64(ctx.Rsp)
	reader.ReaderAt = mem
	dumpStack(mem, int64(ctx.Rsp))

	a, err := reader.ReadPtr()
	if err != nil {
		fmt.Println("can't read mem", err)
		return
	}

	fmt.Printf("stack at 0x%x = %d\n", reader.Pos-int64(reader.PtrSize), a)
}
