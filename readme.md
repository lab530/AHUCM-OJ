# AHUCM OJ

## 环境

### Python3

配置脚本依赖于 Python3 与其对应的 pip，确保安装的 Python 3 版本在 3.8 及其以上。

例：在 Ubuntu 20.04 上，使用以下指令直接安装：

```bash
sudo apt install python3 python3-pip
```

pip 换中国源参考：

1. <https://mirrors.ustc.edu.cn/help/pypi.html>
2. <https://mirrors.tuna.tsinghua.edu.cn/help/pypi/>

### Docker

配置官网提供的较新版本的 Docker 和其对应的 Docker compose。

例：在 Ubuntu 20.04 上，参考 <https://docs.docker.com/engine/install/ubuntu/> 上的安装方法，而非直接使用 Ubuntu 官方 apt 源中的 Docker。

### Go

开发使用的版本为 1.21。

下载链接：<https://go.dev/dl/>

### Vue.js

node 版本: v18.19.0 链接<https://www.nodejs.com.cn/download.html>

npm  版本: 10.5.0 （以实际的 node v18.19.0 版本自带的 npm 版本为主）

npm 换源

```bash
npm config set registry https://registry.npmmirror.com
```
查看当前 npm 的下载源：

```bash
npm config get registry
```
全局安装 yarn

``` bash
cd frontend
npm install -g yarn
```

安装项目依赖

```bash
yarn install
```
运行项目

```bash
yarn serve
```

### Rust

开发环境为 `rustc 1.76.0-nightly (49b3924bd 2023-11-27)`，安装使用更新版本的 nightly 版本。

安装参考：<https://www.rust-lang.org/learn/get-started>
