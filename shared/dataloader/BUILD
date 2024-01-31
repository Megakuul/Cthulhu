# gazelle:exclude *.hpp

load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_dataloader",
    srcs = ["dataloader.go"],
    importpath = "github.com/megakuul/cthulhu/shared/dataloader",
    visibility = ["//visibility:public"],
    deps = ["@com_github_spf13_viper//:viper"],
)

cc_library(
    name = "cc_dataloader",
    hdrs = ["logger.hpp"],
    copts = ["-std=c++23"],
    visibility = ["//visibility:public"],
)
