cc_binary(
    name = "wave",
    srcs = glob(["src/*.cc"]) + glob(["include/*.h"]),
    copts = [
        "-Iwave/include",
        "-std=c++23",
    ],
    deps = ["@boost//:asio", "//shared/logger:cc_logger"],
)
