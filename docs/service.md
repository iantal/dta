# DTA (Dependency Tree Analyzer) service

DTA is a service that is able to generate a dependency tree for a specific commit of a project. Before generating the dependency tree, it queries `btd` tool for the used build tool. After the tree is generated, it will send the data to one of the parsers, depending on the used build tool.

The response from the parse services is then forwarded to the `collector` service that will download the libraries from the web.

## Database
---
### `Table: Statuses`
### `Primary Key`

- `Columns`: ProjectID, Commit

### `Columns`

| `Name`       | `Type`      | `Nullable` | `Default` | `Comment` |
| ------------ | ----------- | ---------- | --------- | --------- |
| ProjectID    | varchar(50) | `false`    |           |           |
| Commit       | varchar(50) | `false`    |           |           |
| Name         | varchar(50) | `false`    |           |           |
| Status       | varchar(50) | `false`    |           |           |
---

Status:
---
* DOWNLOAD_SUCCESS
* DOWNLOAD_FAILURE
---
* BUILD_TOOL_SUCCESS
* BUILD_TOOL_FAILURE
---
* DEPENDENCY_TREE_SUCCESS
* DEPENDENCY_TREE_FAILURE
---
* PARSE_SUCCESS
* PARSE_FAILURE

Name:
---
`If` multi-project build `=>` subproject name `else` `=>` "root"
