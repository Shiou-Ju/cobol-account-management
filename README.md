# Mandarin & English Introduction

## 簡介

簡易的存提款交易系統，使用者除了存提款外，可於聊天室進行即時通訊。

1. 存提款：利用 COBOL 程式計算交易結果，輸出到 Express server 中間層，介接 PostgreSQL 資料庫以及 Vue3 使用者介面。
2. 聊天室：透過 Go 介接 Redis 實作的 Websocket 聊天室，使用者之間可以進行互動。

## 聊天室特色
1. Go server 並無使用框架，並自製 Middleware 來處理 CORS。
2. 聊天室在選擇角色或離開時，會透過鎖來避免資源衝突。
3. 一旦角色選定後，該用戶不可再被選取。
4. Heart-beat 偵測 Websocket 連線狀態，決定釋放哪些用戶，使其可再被選取。


### 快速啟動
1. 請將 `.env.example` 複製並命名為 `.env` 
2. `.env` 請填入適當的 URL，例如在本機執行，請填入：
   - `VUE_APP_API_BASE_URL=http://localhost:3001`
   - `VUE_APP_WEBSOCKET_BASE_URL=ws://localhost:3001`
3. 確保已安裝 Docker 和 Docker Compose 並執行以下命令：

```bash
docker-compose -f docker-compose.yml up --build -d
```

建構完成後，可通過 http://localhost:3000 訪問。

## 部署時路由注意事項
以 Nginx 為例，請注意須將 `/go-api/` 的請求導到機器的 `3001` port 上。


### 用戶界面介紹

- **用戶列表（`AppHome.vue`）**：所有用戶列表，可連接到每位用戶頁面，並顯示用戶最後一筆交易時間。
- **個別用戶（`SingleUserInfo.vue`）**：顯示用戶最近 10 筆交易記錄和當前餘額，可由此進入交易頁面。
- **操作交易（`TransactionForm.vue`）**：允許用戶選擇進行存款或提款操作，交易完成後頁面會刷新。
- **聊天室（`ChatRoom.vue`）**：使用者必須先選擇用戶才能進入聊天室，選中的用戶狀態會在離開聊天室時被清除。


### 其他檔案說明
1. `docker-compose.dev.yml`：可迅速建立本機開發所需要的資料庫
2. `build-go.sh`：Go server 本機開發時 build 的腳本。

***

## Introduction

A simple deposit and withdrawal transaction system that allows users not only to deposit and withdraw but also to communicate in real-time via a chatroom.

1. Banking Transactions: Utilizes a COBOL program to calculate transaction results, outputting to an Express server, which interacts with a PostgreSQL database and a Vue3 user interface.
2. Chatroom: Implemented via Go interfacing with Redis for a WebSocket-based chatroom, enabling real-time interaction among users.

## Chatroom Features
1. The Go server does not use a framework and has a custom middleware handling CORS.
2. To prevent resource conflicts, locking mechanisms are employed when selecting roles or exiting the chatroom.
3. Once a role is selected, that user cannot be selected again.
4. The Go server also employs heartbeat checks on WebSocket connection statuses to decide which users to release, making them selectable again.

### Quick Start
1. Copy `.env.example` and rename it to `.env`
2. Fill in the appropriate URLs in `.env`, take local use for example:
   - `VUE_APP_API_BASE_URL=http://localhost:3001`
   - `VUE_APP_WEBSOCKET_BASE_URL=ws://localhost:3001`
3. Ensure Docker and Docker Compose are installed and run the following command:

```bash
docker-compose -f docker-compose.yml up --build -d
```

After construction, it can be accessed through http://localhost:3000

## Deployment Routing Considerations
For instance, with Nginx, ensure that requests to `/go-api/` are directed to the `3001` port on your machine.

### User Interface Introduction

- **User List (`AppHome.vue`)**: List of all users, with a link to each user's info page, displaying the last transaction time of the user.
- **Individual User (`SingleUserInfo.vue`)**: Displays the user's recent 10 transaction records and current balance, from where you can enter the transaction operation page.
- **Transaction Operation (`TransactionForm.vue`)**: Allows users to perform deposit or withdrawal operations, and the page will refresh after the transaction is completed.
- **Chatroom (`ChatRoom.vue`)**: Users must select a user before entering the chatroom, and the selected user's status is cleared upon exiting the chatroom.

### Additional Files Description
1. `docker-compose.dev.yml`: Quickly sets up the databases required for local development.
2. `build-go.sh`: A script for building the Go server during local development.
