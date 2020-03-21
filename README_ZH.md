# Telebot-Ex

基于 [telebot](https://gopkg.in/tucnak/telebot.v2) 的 Telegram Bot 框架。

## 功能
- 以插件形式注册 Handler
- 为插件指定 Interceptor
- 统一的配置加载

## Motivation

我想在单个 Bot 上跑功能较为隔离的多个 Handler；对于有的 Handler，我需要做鉴权。
全糊在一起就很蛋疼，把公共逻辑抽出来看起来舒服一些吧。