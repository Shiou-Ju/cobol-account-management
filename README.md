# Mandarin & English Introduction

## 簡介

簡易的存提款交易系統，利用 COBOL 程式計算交易結果，輸出到 Express server 中間層，介接 PostgreSQL 資料庫以及 Vue3 使用者介面


### 快速啟動
確保已安裝 Docker 和 Docker Compose 並執行以下命令：

```bash
docker-compose up --build
```

完成後通過 http://localhost:3000 訪問。


### 用戶界面介紹

- **用戶列表（`AppHome.vue`）**：所有用戶列表，可連接到每位用戶頁面，並顯示用戶最後一筆交易時間。
- **個別用戶（`SingleUserInfo.vue`）**：顯示用戶最近 10 筆交易記錄和當前餘額，可由此進入交易頁面。
- **操作交易（`TransactionForm.vue`）**：允許用戶選擇進行存款或提款操作，交易完成後頁面會刷新。


***

## Introduction

A simple deposit and withdrawal transaction system, utilizing COBOL program to calculate transaction results, outputting to an Express server, which interacts with PostgreSQL database and Vue3 user interface.

### Quick Start
Ensure Docker and Docker Compose are installed and run the following command:

```bash
docker-compose up --build
```

After completion, visit http://localhost:3000.


### User Interface Introduction

- **User List (`AppHome.vue`)**: List of all users, with a link to each user's info page, displaying the last transaction time of the user.
- **Individual User (`SingleUserInfo.vue`)**: Displays the user's recent 10 transaction records and current balance, from where you can enter the transaction operation page.
- **Transaction Operation (`TransactionForm.vue`)**: Allows users to perform deposit or withdrawal operations, and the page will refresh after the transaction is completed.
