package validation

import (
	"github.com/sirupsen/logrus"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/webhook/resourcesemantics"
)

var types = map[schema.GroupVersionKind]resourcesemantics.GenericCRD{
	// v1
	v1.SchemeGroupVersion.WithKind("Task"):     &v1.Task{},
	v1.SchemeGroupVersion.WithKind("Pipeline"): &v1.Pipeline{},
}

// Validator is a container for mutation
type Validator struct {
	Logger *logrus.Entry
}

// NewValidator returns an initialised instance of Validator
func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{Logger: logger}
}

// pipelineValidators is an interface used to group functions mutating pods
type pipelineValidator interface {
	Validate(v1.Pipeline) (validation, error)
	Name() string
}

type validation struct {
	Valid  bool
	Reason string
}

// ValidatePod returns true if a pod is valid
func (v *Validator) ValidatePipeline(pipeline v1.Pipeline) (validation, error) {
	var pipelineName string
	if pipeline.Name != "" {
		pipelineName = pipeline.Name
	} else {
		if pipeline.ObjectMeta.GenerateName != "" {
			pipelineName = pipeline.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", pipelineName)
	log.Print("delete me")

	// list of all validations to be applied to the pod
	validations := []pipelineValidator{
		nameValidator{v.Logger},
	}

	// apply all validations
	for _, v := range validations {
		var err error
		vp, err := v.Validate(pipeline)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !vp.Valid {
			return validation{Valid: false, Reason: vp.Reason}, err
		}
	}

	return validation{Valid: true, Reason: "valid pipeline"}, nil
}
