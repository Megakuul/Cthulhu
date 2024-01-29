load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "granit_lib",
    srcs = ["main.go"],
    importpath = "github.com/megakuul/cthulhu/granit",
    visibility = ["//visibility:private"],
    deps = ["//shared/dataloader:go_dataloader"],
)

go_binary(
    name = "granit",
    embed = [":granit_lib"],
    visibility = ["//visibility:public"],
)
