package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"syscall"
	"template-renderer/test"
	"testing"
)

const yamlData1 = `
a: 1
b:
  c: 2
`

const yamlData2 = `
c: 3
`

const jsonData = `{"a": 1, "b": {"c": 2}}`

const TempPath = "../tmp/"

func TestSimple(t *testing.T) {
	defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), yamlData1})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	file, err := os.ReadFile(testDir(t) + "/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A=1\nB=2", string(file))
}

func TestExtension(t *testing.T) {
	defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template2", "--secrets", yamlData1})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	file, err := os.ReadFile(testDir(t) + "/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A=1", string(file))
}

func TestJSON(t *testing.T) {
	defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template", jsonData})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	file, err := os.ReadFile(testDir(t) + "/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A=1\nB=2", string(file))
}

func TestMultipleData(t *testing.T) {
	defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template3", yamlData1, yamlData2})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	file, err := os.ReadFile(testDir(t) + "/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A=1\nC=3", string(file))
}

func TestObject(t *testing.T) {
	defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template4", yamlData1})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	file, err := os.ReadFile(testDir(t) + "/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A={\"a\":1,\"b\":{\"c\":2}}\nC={\"c\":2}", string(file))
}

func TestMissingData(t *testing.T) {
	defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template3", yamlData1})

	err := rootCmd.Execute()

	test.AssertNotEqual(t, nil, err)
	test.AssertEqual(t, "template: test.txt.template3:2:5: executing \"test.txt.template3\" at <.c>: map has no entry for key \"c\"", err.Error())
}

func TestRuntimeData(t *testing.T) {
	defer os.RemoveAll(TempPath)

	os.MkdirAll(testDir(t), os.ModePerm)
	githubOutput, _ := os.CreateTemp(testDir(t), "")
	defer githubOutput.Close()

	var buffer bytes.Buffer
	os.Setenv("GITHUB_OUTPUT", githubOutput.Name())
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template5", "--output-runtime-placeholder-files", "--runtime", yamlData1})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	githubOutputData, _ := io.ReadAll(githubOutput)
	test.AssertEqual(t, "runtime-placeholder-files=../tmp/TestRuntimeData/test.txt\n", string(githubOutputData))

	file, err := os.ReadFile(testDir(t) + "/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A=1\nB=2", string(file))
}

func TestFilePermissions(t *testing.T) {
	defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template6", "--copy-permissions", "--secrets", yamlData1})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	file, err := os.ReadFile(testDir(t) + "/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A=1", string(file))

	var fileInfo syscall.Stat_t
	if err := syscall.Stat(testDir(t)+"/test.txt", &fileInfo); err != nil {
		assert.Fail(t, err.Error())
	}
	test.AssertEqual(t, os.FileMode(0755).String(), os.FileMode(fileInfo.Mode).String())
}

func TestDirectoryPermissions(t *testing.T) {
	//defer os.RemoveAll(TempPath)

	var buffer bytes.Buffer
	rootCmd := NewRootCmd()
	rootCmd.SetOut(&buffer)
	rootCmd.SetArgs([]string{"-i", "../test/data/templates", "-o", testDir(t), "-t", ".template7", "--copy-permissions", "--secrets", yamlData1})

	err := rootCmd.Execute()

	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "", buffer.String())

	file, err := os.ReadFile(testDir(t) + "/dir/test.txt")
	test.AssertEqual(t, nil, err)
	test.AssertEqual(t, "A=1", string(file))

	var fileInfo syscall.Stat_t
	if err := syscall.Stat(testDir(t)+"/dir/test.txt", &fileInfo); err != nil {
		assert.Fail(t, err.Error())
	}
	test.AssertEqual(t, os.FileMode(0755).String(), os.FileMode(fileInfo.Mode).String())
}

func testDir(t *testing.T) string {
	return TempPath + t.Name()
}
