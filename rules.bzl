# This file contains legacy non-module rules

load("@com_github_nelhage_rules_boost//:boost/boost.bzl", "boost_deps")

# Call toolchain initialization rules
def _non_module_init_rules_impl(_ctx):
    boost_deps()

# Use module_extension to extend the bzlmod system with the legacy rules
non_module_init_rules = module_extension(
    implementation = _non_module_init_rules_impl,
)
