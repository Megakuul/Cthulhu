# This file contains legacy non-module loaded http dependencies

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Load HTTP archive dependencies
def http_dependencies():
    # Boost nelhage ruleset
    http_archive(
        name = "com_github_nelhage_rules_boost",
        url = "https://github.com/nelhage/rules_boost/archive/76ed276e0ea602f83d35687084f5f8751ecd2bcb.tar.gz",
        strip_prefix = "rules_boost-76ed276e0ea602f83d35687084f5f8751ecd2bcb",
    )

# Encapsulate dependency loader into a extension wrapper
def _non_module_dependencies_impl(_ctx):
    http_dependencies()

# Use module_extension to extend the bzlmod system with the legacy rules
non_module_dependencies = module_extension(
    implementation = _non_module_dependencies_impl,
)
