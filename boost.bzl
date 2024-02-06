# Extension to load boost rules from nelhage rulset

load("@com_github_nelhage_rules_boost//:boost/boost.bzl", "boost_deps")

def boost_dependency():
    boost_deps()
