package mdstandards

import rego.v1

# Fail closed when the standard is not loaded. In the normal path the binary
# rejects this before evaluation; this rule keeps the bundle honest when it is
# run directly with opa eval.
deny contains finding if {
	not data.standard.dependable_statuses
	finding := {
		"rule": "policy.load_error",
		"document": "",
		"target": "",
		"severity": "error",
		"message": "standard.json has no dependable_statuses",
	}
}

# Rule 2: the referenced document's status must be dependable. Exact match,
# case-sensitive. v1 is intentionally stricter than PDS-STD-004 8.5.2 because
# reference notes are not parsed; see the design doc.
deny contains finding if {
	some doc in input.documents
	is_spec_scope(doc)
	some target in doc.depends_on
	target_doc := index[target]
	status := object.get(target_doc, "status", "")
	not dependable_for(target_doc, status)
	finding := {
		"rule": "refs.status_not_dependable",
		"document": doc.path,
		"target": target,
		"severity": "error",
		"message": sprintf("%s depends on %s with status %q, which is not dependable", [doc.path, target, status]),
	}
}
