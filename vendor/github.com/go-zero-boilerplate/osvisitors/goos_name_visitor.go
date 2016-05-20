package osvisitors

type GoOSNameVisitor struct{ Name string }

func (g *GoOSNameVisitor) VisitWindows() { g.Name = "windows" }
func (g *GoOSNameVisitor) VisitLinux()   { g.Name = "linux" }
func (g *GoOSNameVisitor) VisitDarwin()  { g.Name = "darwin" }
