package validation

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// nameValidator is a container for validating the name of pods
type nameValidator struct {
	Logger logrus.FieldLogger
}

// nameValidator implements the pipelineValidator interface
var _ pipelineValidator = (*nameValidator)(nil)

// Name returns the name of nameValidator
func (n nameValidator) Name() string {
	return "name_validator"
}

// Validate inspects the name of a given pipeline and returns validation.
// The returned validation is only valid if the pipeline name does not contain some
// bad string.
func (n nameValidator) Validate(pipeline v1.Pipeline) (validation, error) {
	badString := "offensive"

	if strings.Contains(pipeline.Name, badString) {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintf("pipeline name contains %q", badString),
		}
		return v, nil
	}

	return validation{Valid: true, Reason: "valid name"}, nil
}
