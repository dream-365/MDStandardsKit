package mdstandards

import rego.v1

# Idea rules only run when the repository declared type: idea. The CLI
# rejects a missing profile before evaluation; this guard keeps the bundle
# honest when it is run directly with opa eval.
is_idea if {
	object.get(input, "document_type", "spec") == "idea"
}

profile := data.standard.document_types.idea

# The extractor fills doc.missing for idea documents, so Rego does not need
# to guess whether an omitted JSON field means "missing" or "empty".
deny contains finding if {
	is_idea
	some doc in input.documents
	some field in doc.missing
	finding := {
		"rule": "idea.missing_field",
		"document": doc.path,
		"target": field,
		"severity": "error",
		"message": sprintf("%s is missing required field %q", [doc.path, field]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	object.get(doc, "title", "") == ""
	finding := {
		"rule": "idea.empty_title",
		"document": doc.path,
		"target": "title",
		"severity": "error",
		"message": sprintf("%s has an empty title", [doc.path]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	doc.id
	not regex.match(profile.id_pattern, doc.id)
	finding := {
		"rule": "idea.id_format",
		"document": doc.path,
		"target": doc.id,
		"severity": "error",
		"message": sprintf("%s has id %q, which does not match %s", [doc.path, doc.id, profile.id_pattern]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	doc.id
	some other in input.documents
	other.path != doc.path
	other.id == doc.id
	finding := {
		"rule": "idea.id_duplicate",
		"document": doc.path,
		"target": other.path,
		"severity": "error",
		"message": sprintf("%s duplicates id %q with %s", [doc.path, doc.id, other.path]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	doc.date
	not regex.match("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", doc.date)
	finding := {
		"rule": "idea.date_format",
		"document": doc.path,
		"target": doc.date,
		"severity": "error",
		"message": sprintf("%s has date %q, which is not YYYY-MM-DD", [doc.path, doc.date]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	doc.date
	regex.match("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", doc.date)
	doc.date_invalid
	finding := {
		"rule": "idea.date_invalid",
		"document": doc.path,
		"target": doc.date,
		"severity": "error",
		"message": sprintf("%s has date %q, which is not a valid calendar date", [doc.path, doc.date]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	parts := split(doc.path, "/")
	count(parts) != 3
	finding := {
		"rule": "idea.path_shape",
		"document": doc.path,
		"target": "",
		"severity": "error",
		"message": sprintf("%s must be YYYY/YYYY-MM-DD/<slug>.md", [doc.path]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	parts := split(doc.path, "/")
	count(parts) == 3
	not regex.match("^20[0-9]{2}$", parts[0])
	finding := {
		"rule": "idea.path_shape",
		"document": doc.path,
		"target": parts[0],
		"severity": "error",
		"message": sprintf("%s has an invalid year directory %q", [doc.path, parts[0]]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	parts := split(doc.path, "/")
	count(parts) == 3
	not regex.match("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", parts[1])
	finding := {
		"rule": "idea.path_shape",
		"document": doc.path,
		"target": parts[1],
		"severity": "error",
		"message": sprintf("%s has an invalid date directory %q", [doc.path, parts[1]]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	parts := split(doc.path, "/")
	count(parts) == 3
	doc.date
	doc.date != parts[1]
	finding := {
		"rule": "idea.date_dir_mismatch",
		"document": doc.path,
		"target": parts[1],
		"severity": "error",
		"message": sprintf("%s has date %q but lives in %q", [doc.path, doc.date, parts[1]]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	parts := split(doc.path, "/")
	count(parts) == 3
	doc.date
	not startswith(doc.date, parts[0])
	finding := {
		"rule": "idea.year_dir_mismatch",
		"document": doc.path,
		"target": parts[0],
		"severity": "error",
		"message": sprintf("%s has date %q but lives in year directory %q", [doc.path, doc.date, parts[0]]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	doc.status
	not status_allowed(doc.status)
	finding := {
		"rule": "idea.status_invalid",
		"document": doc.path,
		"target": doc.status,
		"severity": "error",
		"message": sprintf("%s has status %q, which is not allowed", [doc.path, doc.status]),
	}
}

status_allowed(s) if {
	some v in profile.status_values
	v == s
}

deny contains finding if {
	is_idea
	some doc in input.documents
	doc.tags_invalid
	finding := {
		"rule": "idea.tags_invalid",
		"document": doc.path,
		"target": "tags",
		"severity": "error",
		"message": sprintf("%s has tags that are not a list of strings", [doc.path]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	some heading in profile.required_headings
	not heading_present(doc, heading)
	finding := {
		"rule": "idea.heading_missing",
		"document": doc.path,
		"target": heading,
		"severity": "error",
		"message": sprintf("%s is missing the required heading %q", [doc.path, heading]),
	}
}

heading_present(doc, heading) if {
	some s in doc.sections
	s.level == 2
	s.title == heading
}

deny contains finding if {
	is_idea
	some doc in input.documents
	some s in doc.sections
	s.placeholder
	finding := {
		"rule": "idea.placeholder_remaining",
		"document": doc.path,
		"target": s.title,
		"severity": "error",
		"message": sprintf("%s still contains a template placeholder in %q", [doc.path, s.title]),
	}
}

deny contains finding if {
	is_idea
	some doc in input.documents
	doc.id == "ABCDEFGH"
	finding := {
		"rule": "idea.placeholder_remaining",
		"document": doc.path,
		"target": "id",
		"severity": "error",
		"message": sprintf("%s still has the template id ABCDEFGH", [doc.path]),
	}
}
