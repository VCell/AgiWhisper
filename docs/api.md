# 1 talk_manual

手动采集音频的接口，负责处理一从轮音频采集到回答的过程。

客户端需把音频分成最大8k的分片，逐帧传给服务端。当遇到停顿（断句）的时候，应该设置action=split；当用户点击提问时，应该设置action=ask，此时，后端会停止接收音频接收，开始处理。

处理结果的返回分三个阶段：
1 返回一帧type=question,此时返回data是从音频中提取出的问题；
2 传递agi的所有流式返回，type=slice,data是答案的分片
3 返回完整答案，type=answer,data是所有2阶段的data的拼接


## 1.1 接口概述

- **协议类型**：WebSocket
- **路径**：`/talk_manual`
- **请求内容类型**：JSON
- **响应内容类型**：JSON

## 1.2 请求帧格式

客户端发送的请求帧是一个JSON对象，包含以下字段：

- `action`: 操作，字符串枚举类型，值包括 `"split"`, `"ask"`.
- `audio`: 字符串类型，格式音频的base64

### 例子

```json
{
  "status": "split",
  "audioData": "Base64 encoded audio data"
}
```

## 1.3 返回帧格式

服务端的返回帧是一个JSON对象，包含以下字段：

- `type`: 返回帧类型，包括 `"question"`, `"answer"`, `"slice"`.
- `data`: 返回帧的数据.



### 例子

```json
{
  "key": "abc",
  "info": "abcdefg",
}
```

# 2 talk_auto

自动采集音频的接口。客户端持续收集音频，由服务端决定哪些音频有意义。

## 2.1 接口概述

- **协议类型**：WebSocket
- **路径**：`/talk_auto`
- **请求内容类型**：JSON
- **响应内容类型**：JSON

## 2.2 请求帧格式

客户端发送的请求帧是一个JSON对象，包含以下字段：

- `status`: 状态，字符串枚举类型，值包括 `"unfinished"`, `"finished"`.自动模式下，status=finished代表请求结束，服务端会主动断开
- `audio`: 字符串类型，pcm格式音频的base64
- `id`: 当前分片的编号。预期编号是递增的int类型
- `cmd`: 指令，字符串枚举类型，值包括 `"query"`.没有指令时可以不传或传空字符串

### 例子

```json
{
  "status": "unfinished",
  "audioData": "Base64 encoded audio data",
  "status": 1,
  "cmd": "query"
}
```

## 2.3 返回帧格式

服务端的返回帧是一个JSON对象，包含以下字段：

- `code`: 错误码，处理正常的话为0或空.
- `key`: 识别到的关键词或者提取出的问题.
- `info`: 对应的扩展信息
