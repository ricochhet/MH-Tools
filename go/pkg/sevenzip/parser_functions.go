package sevenzip

func Execute(fnName string, fnArgs []string) {
	switch fnName {
	case "extract":
		builtinExtractFn(fnArgs)
	case "copy":
		builtinCopyFn(fnArgs)
	case "delete":
		builtinDeleteFn(fnArgs)
	default:
	}
}
