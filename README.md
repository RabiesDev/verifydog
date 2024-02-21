# VerifyDog

サーバーに参加したメンバーの認証を行うシンプルなBot

## 前提条件

- [GoLang のバージョン 1.22 以上](https://go.dev/)
- [Docker および Docker Compose](https://www.docker.com/)

## インストールと実行

1. このリポジトリをクローン

    ```bash
    git clone https://github.com/RabiesDev/verifydog
    ```

2. プロジェクトのディレクトリに移動

    ```bash
    cd verifydog
    ```

3. .env に必要な情報を設定します

4. Docker Compose を使用してデータベースを起動

    ```bash
    docker-compose up
    ```

5. Botを起動する

   ```bash
   go run main.go
   ```

## 使用方法

* `/setup-verification` を実行すると認証用Embedが送信されます
* ユーザーが認証を完了すると、指定したロールが付与されます

## 貢献

* もしバグを見つけたり、改善の提案があれば、Issueを作成してください
* 貢献をしたい方は、Pull Requestを送ってください

## ライセンス

このプロジェクトは GPL-3.0 ライセンスの下で提供されています
