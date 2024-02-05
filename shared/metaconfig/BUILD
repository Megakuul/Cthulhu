# gazelle:exclude *.hpp

load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_metaconfig",
    srcs = ["metaconfig.go"],
    importpath = "github.com/megakuul/cthulhu/shared/metaconfig",
    visibility = ["//visibility:public"],
)

cc_library(
    name = "cc_metaconfig",
    hdrs = ["metaconfig.hpp"],
    copts = ["-std=c++23"],
    visibility = ["//visibility:public"],
)
