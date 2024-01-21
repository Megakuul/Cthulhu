load("@gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/megakuul/cthulhu
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=go_deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)
