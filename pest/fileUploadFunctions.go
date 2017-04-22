package pest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type UploadFunctions struct{}

func (_ UploadFunctions) sourceCodeUpload(req *http.Request) (bool, string) {

	// set max memory to hold uploaded file to.. what should this be?
	req.ParseMultipartForm(32 << 20)

	uploadedFile, handler, err := req.FormFile("sourceFile")

	if err != nil {
		return false, "Upload failed. File could not be received. Max file size is: ?????"
	}

	defer uploadedFile.Close()

	fileName := handler.Filename

	// check if it is a .cpp file
	ext := filepath.Ext(fileName)

	if ext != "cpp" {
		return false, "You may only upload .cpp files."
	}

	// build the save path depending on the class, assignment, student name, and sub number -_-
	savePath := "some path"

	err = saveFile(savePath, uploadedFile)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

// How to accept multiple test case files at once.. Just have a max number, 5?
// create a folder for this assignment on disk
func (_ UploadFunctions) callCreateAssignment(form Request) (bool, string) {

	err := createAssignment(form.courseName, form.assignmentDisplayName, form.assignmentName, form.runtime, form.numTestCases, form.compilerOptions, form.startDate, form.endDate)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func saveFile(savePath string, inputFile io.Reader) error {

	saveFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return errors.New("Error. Couldn't save file to server. Disk full or access denied.")
	}

	defer saveFile.Close()

	// copy the uploaded file from memory to the new location
	io.Copy(saveFile, inputFile)

	return nil
}
