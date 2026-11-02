package mdstandards

import rego.v1

# Rule 1: every depends_on target must exist in the governed corpus.
deny contains finding if {
	object.get(input, "document_type", "spec") == "spec"
	some doc in input.documents
	some target in doc.depends_on
	not index[target]
	finding := {
		"rule": "refs.missing_target",
		"document": doc.path,
		"target": target,
		"severity": "error",
		"message": sprintf("%s depends on %s, which is not a governed document", [doc.path, target]),
	}
}

index := {doc.path: doc |
	some doc in input.documents
}
