load("@rules_go//go:def.bzl", "go_library")

cc_library(
    name = "cc_logger",
    hdrs = ["logger.hpp"],
    copts = ["-std=c++23"],
    visibility = ["//visibility:public"],
    deps = ["//shared/util:chan"],
)

go_library(
    name = "logger",
    srcs = ["logger.go" ],
    importpath = "github.com/megakuul/cthulhu/shared/logger",
    visibility = ["//visibility:public"],
)
