load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:defs.bzl", "go_proto_library")
load("@rules_cc//cc:defs.bzl", "cc_proto_library")

# Generic proto library rule
proto_library(
    name = "your_proto",
    srcs = ["your_file.proto"],
    deps = [],
)

# Go-specific proto library
go_proto_library(
    name = "your_go_proto",
    importpath = "example.com/project/your_proto",
    proto = ":your_proto",
    deps = [],
)

# C++-specific proto library
cc_proto_library(
    name = "your_cpp_proto",
    deps = [":your_proto"],
)
