package internal

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"strings"
)

func Parse(yaml string) (obj runtime.Object, kind *schema.GroupVersionKind, err error) {

	decoder := scheme.Codecs.UniversalDeserializer()
	var resource runtime.Object
	var groupVersionKind *schema.GroupVersionKind
	for _, resourceYAML := range strings.Split(string(yaml), "---") {

		// skip empty documents, `Decode` will fail on them
		if len(resourceYAML) == 0 {
			continue
		}

		// - obj is the API object (e.g., Deployment)
		// - groupVersionKind is a generic object that allows
		//   detecting the API type we are dealing with, for
		//   accurate type casting later.
		obj, groupVersionKindInfo, err := decoder.Decode(
			[]byte(resourceYAML),
			nil,
			nil)
		if err != nil {
			log.Print(err)
			continue
		}

		resource = obj
		groupVersionKind = groupVersionKindInfo
	}

	return resource, groupVersionKind, err
}
