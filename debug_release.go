//+build release

package ezdbg

func DEBUG(args ...interface{}) {}

func DEBUGSkip(skip int, args ...interface{}) {}

func PRETTY(args ...interface{}) {}

func PRETTYSkip(skip int, args ...interface{}) {}

func SPEW(args ...interface{}) {}

func SPEWSkip(skip int, args ...interface{}) {}

func DUMP(args ...interface{}) {}

func DUMPSkip(skip int, args ...interface{}) {}
