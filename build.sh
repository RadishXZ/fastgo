#!/bin/bash

# 获取脚本所在目录作为根目录
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")

# 定义编译后的输出目录
OUTPUT_DIR=${PROJ_ROOT_DIR}/_output

# 指定版本信息的路径, 后续会通过ldflags参数将版本注入进去
VERSION_PACKAGE=github.com/RadishXZ/fastgo/pkg/version

# 确定VERSION值, 如果环境变量中没有设置VERSION, 则使用git标签作为版本号
if [[ -z "${VERSION}" ]];then
    VERSION=$(git describe --tags --always --match='v*')
fi

# 检查代码仓库状态: 判断工作目录是否干净
# 默认状态设为"dirty" (有提交更改)
GIT_TREE_STATE="dirty"

# 使用git status检查是否有未提交的更改
is_clean=$(git status --porcelain 2>/dev/null)
# 如果is_clean为空, 说明没有未提交的更改, 状态为"clean"
if [[ -z ${is_clean} ]];then
    GIT_TREE_STATE="clean"
fi

# 获取当前 git commit的完整哈希值
GIT_COMMIT=$(git rev-parse HEAD)

# 构建链接器标志 ldflags
# 通过 -X 选项向VERSION_PACKAGE包中注入以下变量的值
# - gitVersion: 版本号
# - gitTreeState: 代码仓库状态(clean或dirty)
# - buildDate: 构建日期和时间(UTC格式)
GO_LDFLAGS="-X ${VERSION_PACKAGE}.gitVersion=${VERSION} \
  -X ${VERSION_PACKAGE}.gitCommit=${GIT_COMMIT} \
  -X ${VERSION_PACKAGE}.gitTreeState=${GIT_TREE_STATE} \
  -X ${VERSION_PACKAGE}.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

# 执行Go构建命令
# -v: 显示详细编译信息
# -o: 指定输出文件路径和名称
# 最后参数是入口文件路径
go build -v -ldflags "${GO_LDFLAGS}" -o ${OUTPUT_DIR}/fg-apiserver -v cmd/fg-apiserver/main.go