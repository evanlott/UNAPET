package pest

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

// must incorporate naming convention into this
// pull last assignment name, increment it
func (_ UploadFunctions) callCreateAssignment(req *http.Request) (bool, string) {

	form := processForm(req)

	err := createAssignment(form.courseName, form.assignmentDisplayName, form.assignmentName, form.runtime, form.numTestCases, form.compilerOptions, form.startDate, form.endDate)

	if err != nil {
		return false, err.Error()
	}

	assignmentFolder := "data/" + form.assignmentName
	// create a folder for this assignment on disk
	err = os.Mkdir(assignmentFolder, 0755)

	if err != nil {
		return false, "Error creating a directory for this assignment on the server."
	}

	// set max memory to hold uploaded file to.. what should this be?
	req.ParseMultipartForm(32 << 20)

	// accept as many test cases as were uploaded
	for i := 0; ; i++ {

		testCase, testCaseHandler, err0 := req.FormFile("testCase" + strconv.Itoa(i))
		desiredOutput, outputHandler, err1 := req.FormFile("desiredOutput" + strconv.Itoa(i))

		// check if last test case, output pair
		// force them to upload at least 1 pair
		if (err0 != nil || err1 != nil) && i > 1 {
			break
		} else {
			os.RemoveAll(assignmentFolder)
			return false, "Upload failed. File could not be received. You must upload at least one (test case, output) pair."
		}

		defer testCase.Close()
		defer desiredOutput.Close()

		testCaseFileName := testCaseHandler.Filename
		outputFileName := outputHandler.Filename

		ext1 := filepath.Ext(testCaseFileName)
		ext2 := filepath.Ext(outputFileName)

		if ext1 != "txt" || ext2 != "txt" {
			os.RemoveAll(assignmentFolder)
			return false, "You may only upload .txt files for test cases and desired outputs."
		}

		// build the save path
		testCaseSavePath := "some path"
		outputSavePath := "some path"

		err0 = saveFile(testCaseSavePath, testCase)
		err1 = saveFile(outputSavePath, desiredOutput)

		if err0 != nil || err1 != nil {
			os.RemoveAll(assignmentFolder)
			return false, "Could not save test case or desired output to server."
		}
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
