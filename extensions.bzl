# Extension file to load legacy non-module bazel mods

# Load http dependencies
load("//:deps.bzl", "http_dependencies")
http_dependencies()

# Load rules
load("//:rules.bzl", "dependency_rules")

# Load dependencies
def _non_module_dependencies_impl(_ctx):
    dependency_rules()

# Use module_extension to extend the bzlmod system with the legacy rules
non_module_dependencies = module_extension(
    implementation = _non_module_dependencies_impl,
)
