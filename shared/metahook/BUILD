# gazelle:exclude *.hpp

load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_metahook",
    srcs = ["metahook.go"],
    importpath = "github.com/megakuul/cthulhu/shared/metahook",
    visibility = ["//visibility:public"],
)

cc_library(
    name = "cc_metahook",
    hdrs = ["metahook.hpp"],
    copts = ["-std=c++23"],
    visibility = ["//visibility:public"],
)
