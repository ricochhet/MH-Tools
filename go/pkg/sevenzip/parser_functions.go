package sevenzip

func Execute(fnName string, fnArgs []string) {
	switch fnName {
	case "extract":
		BuiltinExtractFn(fnArgs)
	case "copy":
		BuiltinCopyFn(fnArgs)
	case "delete":
		BuiltinDeleteFn(fnArgs)
	default:
	}
}
