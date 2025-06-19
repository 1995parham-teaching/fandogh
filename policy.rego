package authz

default allow := false

allow if {
    input.method == "POST"
    input.path == "/api/homes"
}
