package main

// Error can ship a cmd output as well as the start interface. Useful for understanding why a system command (exec.Command) failed
type Error struct {
	err       error
	CmdOutput string
}

//
// Used to create an error of a command line command did not work.
//
func makeErr(err error, out string) *Error {

	if err != nil {
		return &Error{
			err:       err,
			CmdOutput: out,
		}
	}

	return nil
}

/* End File */
