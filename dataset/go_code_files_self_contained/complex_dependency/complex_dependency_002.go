package main

func lockedSetupOutputMap() {
	newOutputFuncMap = make(map[string]NewOutputFunc, 100)
	newOutputFuncMap["stdout"] = newOutputFuncStdout
	newOutputFuncMap["stderr"] = newOutputFuncStderr
	newOutputFuncMap["file"] = newOutputFuncFile
	newOutputFuncMap["fd"] = newOutputFuncFd
	newOutputFuncMap["discard"] = newOutputFuncDiscard
	outputMap = make(map[string]*outputWrapper, 100)
}
