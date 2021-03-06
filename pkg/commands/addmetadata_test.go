/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"reflect"
	"testing"

	"github.com/kubernetes-sigs/kustomize/pkg/constants"
	"github.com/kubernetes-sigs/kustomize/pkg/fs"
)

func TestParseValidateInput(t *testing.T) {
	var testcases = []struct {
		input        string
		valid        bool
		name         string
		expectedData map[string]string
		kind         KindOfAdd
	}{
		{
			input: "otters:cute",
			valid: true,
			name:  "Adds single input",
			expectedData: map[string]string{
				"otters": "cute",
			},
			kind: label,
		},
		{
			input: "owls:great,unicorns:magical",
			valid: true,
			name:  "Adds two items",
			expectedData: map[string]string{
				"owls":     "great",
				"unicorns": "magical",
			},
			kind: label,
		},
		{
			input: "123:45",
			valid: true,
			name:  "Numeric input is allowed",
			expectedData: map[string]string{
				"123": "45",
			},
			kind: annotation,
		},
		{
			input:        " ",
			valid:        false,
			name:         "Empty space input",
			expectedData: nil,
			kind:         annotation,
		},
	}
	var o addMetadataOptions
	for _, tc := range testcases {
		args := []string{tc.input}
		err := o.ValidateAndParse(args, tc.kind)
		if err != nil && tc.valid {
			t.Errorf("for test case %s, unexpected cmd error: %v", tc.name, err)
		}
		if err == nil && !tc.valid {
			t.Errorf("unexpected error: expected invalid format error for test case %v", tc.name)
		}
		//o.metadata should be the same as expectedData
		if tc.valid {
			if !reflect.DeepEqual(o.metadata, tc.expectedData) {
				t.Errorf("unexpected error: for test case %s, unexpected data was added", tc.name)
			}
		} else {
			if len(o.metadata) != 0 {
				t.Errorf("unexpected error: for test case %s, expected no data to be added", tc.name)
			}
		}
	}
}

func TestRunAddAnnotation(t *testing.T) {
	fakeFS := fs.MakeFakeFS()
	fakeFS.WriteFile(constants.KustomizationFileName, []byte(kustomizationContent))
	var o addMetadataOptions
	o.metadata = map[string]string{"owls": "cute", "otters": "adorable"}

	err := o.RunAddAnnotation(fakeFS, annotation)
	if err != nil {
		t.Errorf("unexpected error: could not write to kustomization file")
	}
	// adding the same test input should not work
	err = o.RunAddAnnotation(fakeFS, annotation)
	if err == nil {
		t.Errorf("expected already in kustomization file error")
	}
	// adding new annotations should work
	o.metadata = map[string]string{"new": "annotation"}
	err = o.RunAddAnnotation(fakeFS, annotation)
	if err != nil {
		t.Errorf("unexpected error: could not write to kustomization file")
	}
}

func TestAddAnnotationNoArgs(t *testing.T) {
	fakeFS := fs.MakeFakeFS()
	cmd := newCmdAddAnnotation(fakeFS)
	err := cmd.Execute()
	if err == nil {
		t.Errorf("expected an error but error is %v", err)
	}
	if err != nil && err.Error() != "must specify annotation" {
		t.Errorf("incorrect error: %v", err.Error())
	}
}
func TestAddAnnotationMultipleArgs(t *testing.T) {
	fakeFS := fs.MakeFakeFS()
	fakeFS.WriteFile(constants.KustomizationFileName, []byte(kustomizationContent))
	cmd := newCmdAddAnnotation(fakeFS)
	args := []string{"this:annotation", "has:spaces"}
	err := cmd.RunE(cmd, args)
	if err == nil {
		t.Errorf("expected an error but error is %v", err)
	}
	if err != nil && err.Error() != "annotations must be comma-separated, with no spaces. See help text for example" {
		t.Errorf("incorrect error: %v", err.Error())
	}
}

func TestRunAddLabel(t *testing.T) {
	fakeFS := fs.MakeFakeFS()
	fakeFS.WriteFile(constants.KustomizationFileName, []byte(kustomizationContent))
	var o addMetadataOptions
	o.metadata = map[string]string{"owls": "cute", "otters": "adorable"}

	err := o.RunAddLabel(fakeFS, label)
	if err != nil {
		t.Errorf("unexpected error: could not write to kustomization file")
	}
	// adding the same test input should not work
	err = o.RunAddLabel(fakeFS, label)
	if err == nil {
		t.Errorf("expected already in kustomization file error")
	}
	// adding new labels should work
	o.metadata = map[string]string{"new": "label"}
	err = o.RunAddLabel(fakeFS, label)
	if err != nil {
		t.Errorf("unexpected error: could not write to kustomization file")
	}
}

func TestAddLabelNoArgs(t *testing.T) {
	fakeFS := fs.MakeFakeFS()

	cmd := newCmdAddLabel(fakeFS)
	err := cmd.Execute()
	if err == nil {
		t.Errorf("expected an error but error is: %v", err)
	}
	if err != nil && err.Error() != "must specify label" {
		t.Errorf("incorrect error: %v", err.Error())
	}
}

func TestAddLabelMultipleArgs(t *testing.T) {
	fakeFS := fs.MakeFakeFS()
	fakeFS.WriteFile(constants.KustomizationFileName, []byte(kustomizationContent))
	cmd := newCmdAddLabel(fakeFS)
	args := []string{"this:input", "has:spaces"}
	err := cmd.RunE(cmd, args)
	if err == nil {
		t.Errorf("expected an error but error is: %v", err)
	}
	if err != nil && err.Error() != "labels must be comma-separated, with no spaces. See help text for example" {
		t.Errorf("incorrect error: %v", err.Error())
	}
}
