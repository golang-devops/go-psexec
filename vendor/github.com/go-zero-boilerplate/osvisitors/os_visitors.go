package osvisitors

var (
	WindowsOs OsType = &windows{}
	LinuxOs   OsType = &linux{}
	DarwinOs  OsType = &darwin{}

	AllList []OsType = []OsType{
		WindowsOs,
		LinuxOs,
		DarwinOs,
	}
)

type OsType interface {
	Accept(v OsTypeVisitor)
}

type OsTypeVisitor interface {
	VisitWindows()
	VisitLinux()
	VisitDarwin()
}

type windows struct{}
type linux struct{}
type darwin struct{}

func (*windows) Accept(v OsTypeVisitor) { v.VisitWindows() }
func (*linux) Accept(v OsTypeVisitor)   { v.VisitLinux() }
func (*darwin) Accept(v OsTypeVisitor)  { v.VisitDarwin() }
