# Cthulhu

![Cthulhu Icon](/cthulhu.svg "Cthulhu")

*This is just a proof of concept*

Cloud based hypervisor system for juju controlled environments.

### Development


##### Dependency Management

All build steps are managed by the `bazel` build tool. Other tools shall only be used as bazel plugin.

Go dependencies are managed in the `go.mod` file at the root of the project, this file is processed by `gazelle` to generate the required rules under the hood.

C++ dependencies are managed in the `MODULE.bazel` file at the root of the project.

##### Build Files

Every component must have its own Bazel *BUILD* file.

C++ rules must be manually configured, for Go and Proto rules, you can use `bazel run //:gazelle` to generate the code automatically.

##### Completions / IntelliSense

###### GO

For Go development I suggest using the *gopls* lang-server and import the repository root. Then you should have completion over the whole project.

###### C++

For C++ development I recommend using *clangd* for intellisense / documentation.
To generate the *compile_commands.json* file there are various options, I recommend to use this tool:

[bazel-compile-commands](https://github.com/kiron1/bazel-compile-commands)

[bazel-compile-commands releases](https://github.com/kiron1/bazel-compile-commands/releases)

(In case std is by default not ++23, you can use this option to replace the std: `bazel-compile-commands -R c++14=c++23`)
