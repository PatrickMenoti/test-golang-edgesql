package testedgesql

import (
	"fmt"
	"log"
	"path"

	python "github.com/go-python/cpy3"
)

func runPythonCommand(command string, args ...string) error {
	py := python.PyImport_ImportModule("subprocess")
	if py == nil {
		return fmt.Errorf("failed to import subprocess module")
	}

	pyRun := py.GetAttrString("run")
	if pyRun == nil {
		return fmt.Errorf("failed to get subprocess.run function")
	}

	pyArgs := python.PyList_New(len(args))
	for i, arg := range args {
		python.PyList_SetItem(pyArgs, i, python.PyUnicode_FromString(arg))
	}

	pyResult := pyRun.CallFunctionObjArgs(pyArgs)
	if pyResult == nil {
		return fmt.Errorf("failed to execute command: %s", command)
	}

	return nil
}

func main() {
	// Initialize Python
	if !python.Py_IsInitialized() {
		python.Py_Initialize()
		defer python.Py_Finalize()
	}

	// Install dependencies
	fmt.Println("Installing dependencies...")
	deps := []string{"pip", "install", "mysql-connector-python", "psycopg2"}
	err := runPythonCommand("pip", deps...)
	if err != nil {
		log.Fatalf("Failed to install dependencies: %v", err)
	}

	// Set environment variables
	pyOS := python.PyImport_ImportModule("os")
	if pyOS == nil {
		log.Fatal("Failed to import os module")
	}
	pyEnviron := pyOS.GetAttrString("environ")
	if pyEnviron == nil {
		log.Fatal("Failed to access os.environ")
	}
	pyEnviron.SetAttrString("AZION_TOKEN", python.PyUnicode_FromString("YOUR_TOKEN_HERE"))

	// Run the script
	fmt.Println("Running the script...")
	pyResult, err := python.PyRun_AnyFile(path.Join("edgesql-shell", "edgesql-shell.py"))
	if err != nil {
		log.Fatalf("Error running Python script: %v", err)
	}

	fmt.Printf("Python script executed with return code: %d\n", pyResult)
}
