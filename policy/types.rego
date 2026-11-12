package mdstandards

import rego.v1

# Shared scope helpers for legacy single-type mode and multi-type routing.
# Legacy mode keeps the pre-v1.2 corpus-level guards; multi mode routes by
# each document's canonical type and its rule_set.

is_multi if {
	object.get(input, "mode", "") == "multi"
}

document_types := object.get(data.standard, "document_types", {})
legacy_type := object.get(input, "document_type", "spec")

profile_for(doc) := object.get(document_types, object.get(doc, "type", ""), {}) if {
	is_multi
}

profile_for(_) := object.get(document_types, "idea", {}) if {
	not is_multi
	legacy_type == "idea"
}

# Legacy spec deliberately returns a defined empty object. `object.get` on an
# undefined argument does not fall through to its default, so every legacy
# branch must define profile_for explicitly.
profile_for(_) := {} if {
	not is_multi
	legacy_type == "spec"
}

is_spec_scope(_) if {
	not is_multi
	legacy_type == "spec"
}

is_spec_scope(doc) if {
	is_multi
	object.get(profile_for(doc), "rule_set", "") == "spec"
}

is_idea_scope(_) if {
	not is_multi
	legacy_type == "idea"
}

is_idea_scope(doc) if {
	is_multi
	object.get(profile_for(doc), "rule_set", "") == "idea"
}

# v1 uses the top-level status set for every target type. Per-type overrides
# are deferred to Phase 2.
dependable_for(_, status) if {
	some s in object.get(data.standard, "dependable_statuses", [])
	s == status
}

# Fail closed when the bundle is evaluated directly (e.g. opa eval) without
# the CLI's pre-flight validation.
deny contains finding if {
	is_multi
	count(document_types) == 0
	finding := {
		"rule": "policy.load_error",
		"document": "",
		"target": "document_types",
		"severity": "error",
		"message": "multi-type mode requires document_types",
	}
}

deny contains finding if {
	is_multi
	some doc in input.documents
	object.get(profile_for(doc), "rule_set", "") == ""
	finding := {
		"rule": "policy.load_error",
		"document": doc.path,
		"target": "rule_set",
		"severity": "error",
		"message": sprintf("%s has no rule_set for type %q", [doc.path, object.get(doc, "type", "")]),
	}
}
