# wakka

## これは

俺が, 俺のために作った CircleCI の API を叩くコマンドラインツールです.

* https://circleci.com/docs/api/

wakka は, Circle (円形) ≒ 輪 (わっか) から名付けています. すいません.

## 使い方

1. User Settings > Personal API Token にて Token を生成
2. 環境変数に API Token と CircleCI Username (Github Username) を設定

```sh
export CIRCLECI_USERNAME=username
export CIRCLECI_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

3. ビルドする

```sh
$ brew install upx
$ mkdir pkg
$ make build
```

4. 生成されたバイナリをパスが通ったディレクトリに放り込む

```sh
$ mv pkg/wakka_darwin_amd64 ~/bin/wakka
```

5. 動作確認

```sh
$ wakka version
```

## サポートしている機能

### projects

フォローしているプロジェクト一覧を取得.

```sh
$ wakka projects
```

### variables

プロジェクトの環境変数を管理する.

```sh
# 環境変数一覧を取得
$ wakka variables -project=${YOUR_PROJECT_NAME}

# 環境変数を追加
$ wakka variables -project=${YOUR_PROJECT_NAME} -add -name=VARIABLE_NAME -value=VARIABLE_VALUE

# 環境変数を削除
$ wakka variables -project=${YOUR_PROJECT_NAME} -del -name=VARIABLE_NAME
```

## すいません

色々と荒削りです. 気づいたら直します.
