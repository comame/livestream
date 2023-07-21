## 使い方

1. rtmp://livestream.comame.dev:1935/app に対して配信する
1. https://livestream.comame.xyz/hls/[stream-key]/index.m3u8 を開く
    - https://livestream.comame.xyz/viewer/[stream-key]

## なかみ


### `nginx.conf`

[arut/nginx-rtmp-module](https://github.com/arut/nginx-rtmp-module) で RTMP を受け取り、HLS に変換して配信する

### `viewer`

HLS を移すだけの Web ページ

###  `cleanup/`

VPN 越しに配信すると `index.m3u8` に記載されているはずの `.ts` ファイルが削除されてしまう問題が発生したので、Nginx の `hls_cleanup` を無効にし、自分でクリーンアップすることにしたやつ
