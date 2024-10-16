### 目的

本ドキュメントはチーム制作授業でMTG等のリマインドを自動的に行うDiscordBOT開発に関する要件定義書になります。

### 概要

**システム構成**

システム構成は次の通りです。

- BOTサーバー
    - BOTが常駐するメインサーバーです。リマインドやスケジュール登録を行います。
    - Glab教室にあるサーバーを利用します。
- スケジュール登録システム
    - 各チームリーダーのみが参加できるテキストチャンネルを用意し、スラッシュコマンドを利用してスケジュール登録を行います。

### 業務要件

**構築後のフロー**

システム構築後、フローは次の通りです。

1. 各チームリーダーがリマインドしたい日程をスラッシュコマンドを用いて入力
2. BOT側が認識し記録後、結果をDiscordに通知
3. 指定された日時になった場合にメッセージを送信

**利用者一覧**

- 各チームリーダー
- 各チームメンバー (閲覧のみ)

**規模**

3日間 (エンジニア 1人)

### 機能要件

**リマインドスケジュール管理**

- スラッシュコマンドの認識
    - スケジュール追加・削除・編集を行うコマンドです。
- 日時の設定
    - スケジュール操作コマンドに続く形で、チーム・役職・時間・曜日を設定します。

**リマインドメッセージ送信**

- 各チームへのメッセージ送信
    - 各チームのアナウンスチャンネルに役職メンション付きで送信します。

### データ

- スケジュールをJSONデータで管理します。
- 実行ログはCLI上に出力され、記録はしません。
