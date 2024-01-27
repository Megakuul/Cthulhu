load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_logger",
    srcs = ["dataloader.go"],
    importpath = "github.com/megakuul/cthulhu/shared/dataloader",
    visibility = ["//visibility:public"],
    deps = ["@com_github_spf13_viper//:go_default_library"],
)

cc_library(
	name = "cc_logger",
	hdrs = ["logger.hpp"],
    copts = ["-std=c++23"],
	visibility = ["//visibility:public"]
)
