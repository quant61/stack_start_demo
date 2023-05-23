package main

import (
	"encoding/binary"
	"fmt"
	"github.com/quant61/stack_start_demo/internal/auxv"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

type ProcMemReader struct {
	io.ReaderAt
}

func NewProcMemReader(cmd *exec.Cmd) (*ProcMemReader, error) {
	mem, err := os.Open(fmt.Sprintf("/proc/%d/mem", cmd.Process.Pid))
	return &ProcMemReader{mem}, err
}

func GetRegs(cmd *exec.Cmd) (regs syscall.PtraceRegs, err error) {
	err = syscall.PtraceGetRegs(cmd.Process.Pid, &regs)
	return regs, err
}

func printMaps(cmd *exec.Cmd) {
	b, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/maps", cmd.Process.Pid))
	if err != nil {
		fmt.Println("can't read maps:", err)
	} else {
		fmt.Printf("maps:\n%s\n", string(b))
	}
}

func startProcess(args startProcessArgs) *exec.Cmd {
	cmd := exec.Command("")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	cmd.Args = []string{"one", "two", "three"}
	cmd.Path = "./bin"

	cmd.Env = []string{}
	err := cmd.Start()
	fmt.Println("start process: error =", err)
	return cmd
}

func waitForTraceStop(pid int) (syscall.WaitStatus, error) {
	var w syscall.WaitStatus
	_, err := syscall.Wait4(pid, &w, 0, nil)
	if err != nil {
		fmt.Println("can't wait:", err)
		return w, err
	}
	if w.Signaled() {
		fmt.Println("signal", w.Signal())
	}
	if w.Stopped() {
		fmt.Println("stop", w.StopSignal())
	}
	return w, err
}

func elfRun(elfFactory func() ([]byte, binary.ByteOrder), ptrSize int){
	b, ord := elfFactory()
	err := ioutil.WriteFile("bin", b, 0755)
	fmt.Println("write file: error=", err)
	if err != nil {
		return
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	cmd := startProcess(startProcessArgs{})
	waitForTraceStop(cmd.Process.Pid)

	reader := &readerHelper{
		ByteOrder: ord,
		PtrParser: parserFactoryByPtrSize[ptrSize](ord),
	}
	printProcessState(cmd, reader)
	fmt.Println("kill: error = ", cmd.Process.Kill())
	st, err := cmd.Process.Wait()
	fmt.Println("wait", st, err)
	//// TODO: replace with a more clean solution
	//time.Sleep(100*time.Millisecond)
}


func printProcessState(cmd *exec.Cmd, reader *readerHelper) {
	//ptrSize := int64(reader.PtrParser.Len())

	fmt.Printf("child.pid=%d\n", cmd.Process.Pid)
	//var regs syscall.PtraceRegs
	//syscall.PtraceGetRegs(cmd.Process.Pid, &regs)
	regs, err := GetRegs(cmd)
	if err != nil {
		fmt.Println("cannot get process regs", err)
		return
	}

	printNonZeroFields(regs)
	printMaps(cmd)

	fmt.Printf("rsp=%x\n", regs.Rsp)

	mem, err := NewProcMemReader(cmd)
	if err != nil {
		fmt.Println("can't open process memory: ", err)
		os.Exit(1)
	}
	dumpStack(mem, int64(regs.Rsp))

	reader.pos = int64(regs.Rsp)
	reader.ReaderAt = mem
	// reading argc, argv
	argc, err := reader.ReadPtr()
	fmt.Println("argc", argc, err)
	for i := int64(0); i < (argc); i++ {
		ptr, err := reader.ReadPtr()
		if err != nil {
			fmt.Printf("error on reading ptr to argv[%d]: %v\n", i, err)
			continue
		}
		s, err := reader.CStringAt(int64(ptr), 4096)
		if err != nil {
			fmt.Printf("error on reading value of argv[%d] at 0x%x: %v\n", i, ptr, err)
			continue
		}
		fmt.Printf("argv[%d] at 0x%x = %q\n", i, ptr, s)
	}
	v, err := reader.ReadPtr()
	// argv[argc] ?= NULL
	fmt.Println("zero?", v, err)

	// env
	for {
		ptr, _ := reader.ReadPtr()
		if ptr == 0 {
			fmt.Println("ENV end")
			break
		}
		s, _ := reader.CStringAt(ptr, 4096)
		fmt.Printf("ENV at 0x%x %s\n", ptr, s)
	}

	// auxv part
	for {
		_k, err := reader.ReadPtr()
		k := auxv.Auxk(_k)
		if err != nil {
			fmt.Println("cannot read auxv key", err)
		}
		v, err := reader.ReadPtr()
		if err != nil {
			fmt.Println("cannot read auxv value", err)
		}

		switch k {
		case auxv.AT_PLATFORM, auxv.AT_BASE_PLATFORM, auxv.AT_EXECFN:
			s, _ := reader.CStringAt(int64(v), 4096)
			fmt.Printf("auxv[%s] at 0x%x = %q\n", k, uintptr(v), s)
		default:
			fmt.Printf("auxv[%s] =\t0x%x\t = %d\n", k, uintptr(v), v)
		}
		if k == 0 {
			fmt.Println("auxv end")
			break
		}
	}
}
