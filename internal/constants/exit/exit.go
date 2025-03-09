// Package exit defines preferable exit codes for programs.
// It reflects the definitions found in sysexits.h
package exit

type ExitCode = int

const (
	// Usage is returned when the command was used incorrectly,
	// e.g., with the wrong number of arguments, a bad flag, a bad syntax in a parameter, or whatever.
	Usage ExitCode = 64

	// DataErr is returned when the input data was incorrect in some way.
	// This should only be used for user's data and not system files.
	DataErr ExitCode = 65

	// NoInput is returned when an input file (not a system file) did not exist or was not readable.
	// This could also include errors like "No message" to a mailer (if it cared to catch it).
	NoInput ExitCode = 66

	// NoUser is returned when the user specified did not exist.
	// This might be used for mail addresses or remote logins.
	NoUser ExitCode = 67

	// NoHost is returned when the host specified did not exist.
	// This is used in mail addresses or network requests.
	NoHost ExitCode = 68

	// Unavailable is returned when a service is unavailable.
	// This can occur if a sup port program or file does not exist.
	// This can also be used as a catchall message when something you wanted to do doesn't work, but you don't know why.
	Unavailable ExitCode = 69

	// Software is returned when an internal software error has been detected.
	// This should be limited to non-operating system related errors as possible.
	Software ExitCode = 70

	// OSErr is returned when an operating system error has been detected.
	// This is intended to be used for such things as "cannot fork", "cannot create pipe", or the like.
	// It includes things like getuid returning a user that does not exist in the passwd file.
	OSErr ExitCode = 71

	// OSFile is returned when some system file does not exist, cannot be opened, or has some sort of error.
	OSFile ExitCode = 72

	// CantCreate is returned when a (user specified) output file cannot be created.
	CantCreate ExitCode = 73

	// IOErr is returned when an error occurred while doing I/O on some file.
	IOErr ExitCode = 74

	// TempFail is returned when temporary failure, indicating something that is not really an error.
	// In sendmail, this means that a mailer (e.g.) could not create a connection,
	// and the request should be reattempted later.
	TempFail ExitCode = 75

	// Protocol is returned when the remote system returned something that was "not possible" during a protocol exchange.
	Protocol ExitCode = 76

	// NoPerm is returned when you did not have sufficient permission to perform the operation.
	// This is not intended for file sys­tem problems, which should use [NoInput]
	// or [CantCreate], but rather for higher level permissions.
	NoPerm ExitCode = 77

	// Config is returned when Something was found in an unconfigured or misconfigured state.
	Config ExitCode = 78
)
