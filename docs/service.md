# DTA (Dependency Tree Analyzer) service

DTA is a service that is able to generate a dependency tree for a specific commit of a project. Before generating the dependency tree, it queries `btd` tool for the used build tool. After the tree is generated, it will send the data to one of the parsers, depending on the used build tool.

The response from the parse services is then forwarded to the `collector` service that will download the libraries from the web.