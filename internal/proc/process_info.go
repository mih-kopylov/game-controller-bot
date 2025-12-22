package proc

type ProcessInfo struct {
	Pid       int32
	Name      string
	User      string
	Exe       string
	ParentPid int32
	Terminal  string
}
