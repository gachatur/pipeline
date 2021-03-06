/*
Copyright 2019 The Tekton Authors

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

package workspace

import (
	"context"
	"fmt"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/list"
)

// ValidateBindings will return an error if the bound workspaces in wb don't satisfy the declared
// workspaces in w.
func ValidateBindings(w []v1beta1.WorkspaceDeclaration, wb []v1beta1.WorkspaceBinding) error {
	// This will also be validated at webhook time but in case the webhook isn't invoked for some
	// reason we'll invoke the same validation here.
	for _, b := range wb {
		if err := b.Validate(context.Background()); err != nil {
			return fmt.Errorf("binding %q is invalid: %v", b.Name, err)
		}
	}

	declNames := make([]string, len(w))
	for i := range w {
		declNames[i] = w[i].Name
	}
	bindNames := make([]string, len(wb))
	for i := range wb {
		bindNames[i] = wb[i].Name
	}
	if err := list.IsSame(declNames, bindNames); err != nil {
		return fmt.Errorf("bound workspaces did not match declared workspaces: %v", err)
	}
	return nil
}

// ValidateOnlyOnePVCIsUsed checks that a list of WorkspaceBinding uses only one
// persistent volume claim.
//
// This is only useful to validate that WorkspaceBindings in TaskRuns are compatible
// with affinity rules enforced by the AffinityAssistant.
func ValidateOnlyOnePVCIsUsed(wb []v1beta1.WorkspaceBinding) error {
	workspaceVolumes := make(map[string]bool)
	for _, w := range wb {
		if w.PersistentVolumeClaim != nil {
			workspaceVolumes[w.PersistentVolumeClaim.ClaimName] = true
		}
		if w.VolumeClaimTemplate != nil {
			workspaceVolumes[w.Name] = true
		}
	}

	if len(workspaceVolumes) > 1 {
		return fmt.Errorf("more than one PersistentVolumeClaim is bound")
	}
	return nil
}
