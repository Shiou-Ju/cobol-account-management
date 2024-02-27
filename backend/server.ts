import express from 'express';
import { pgPool } from './connection';
import { exec } from 'child_process';

// TODO: not unser root dir
// File '/Users/bamboo/Repos/cobol-account-management/utils/interfaces.ts' is not under 'rootDir' '/Users/bamboo/Repos/cobol-account-management/backend'. 'rootDir' is expected to contain all source files.
// The file is in the program because:
// Matched by include pattern '../utils/**/*' in '/Users/bamboo/Repos/cobol-account-management/backend/tsconfig.json'ts

import {
  TransactionPayload,
  TransactionRecord,
  TransactionType,
  VALID_TRANSACTION_TYPES,
} from '../utils/interfaces';
import { promisify } from 'util';
import path from 'path';
import { existsSync, mkdirSync } from 'fs';

const execAsync = promisify(exec);

const app = express();
const PORT = 3000;

app.use(express.json());

function isValidTransactionType(type: string): type is TransactionType {
  // TODO: if as is ok here
  return VALID_TRANSACTION_TYPES.includes(type as TransactionType);
}

app.get('/', async (_req, res) => {
  const { rows } = await pgPool.query('SELECT NOW() as now');
  res.send(`Database time is: ${rows[0].now}`);
});

// TODO: add another api for transaction list

app.get('/user/:user', async (req, res) => {
  const userName = req.params.user;

  // TODO: maybe more than one
  // TODO: modify frontend to make transaction list as well
  const GET_LATEST_USER_TRANSACTION = `
    SELECT * FROM (
      SELECT
        transactions.*,
        ROW_NUMBER() OVER (PARTITION BY "user" ORDER BY "date" DESC) as rn
      FROM transactions
      WHERE "user" = $1
    ) as ranked_transactions
    WHERE rn = 1;
  `;

  try {
    const { rows } = await pgPool.query(GET_LATEST_USER_TRANSACTION, [
      userName,
    ]);

    const hasUser = rows.length > 0;

    if (!hasUser) {
      res.status(404).send('User not found');
    }

    const userInfo = rows[0];
    res.json(userInfo);
  } catch (error) {
    console.error(error);
    res.status(500).send('Server error');
  }
});

app.get('/users', async (_req, res) => {
  try {
    const sql = `SELECT * FROM (
      SELECT 
        transactions.*,
        ROW_NUMBER() OVER(PARTITION BY "user" ORDER BY "date" DESC) as rn
      FROM 
        transactions
    ) t
    WHERE t.rn = 1;
    `;

    const { rows } = await pgPool.query(sql);
    res.json(rows);
  } catch (error) {
    console.error(error);
    res.status(500).send('Server error');
  }
});

app.post('/transaction', async (req, res) => {
  // TODO: ts new feature? satisfy
  const { user, transaction, type } = req.body as TransactionPayload;

  const hasMissingData = !user || transaction === undefined || !type;

  if (hasMissingData) {
    return res.status(400).send('Missing required fields');
  }

  if (!isValidTransactionType(type)) {
    return res.status(400).send('invalid transaction type');
  }

  // TODO: transaction must be the latest one
  // const userResult = await pgPool.query(
  //   'SELECT * FROM transactions WHERE "user" = $1 LIMIT 1;',
  //   [user],
  // );

  const GET_LATEST_USER_TRANSACTION = `
SELECT * FROM (
  SELECT
    transactions.*,
    ROW_NUMBER() OVER (PARTITION BY "user" ORDER BY "date" DESC) as rn
  FROM transactions
  WHERE "user" = $1
) as ranked_transactions
WHERE rn = 1;
`;

  const userResult = await pgPool.query(GET_LATEST_USER_TRANSACTION, [user]);

  const hasUser = userResult.rows.length > 0;

  if (!hasUser) {
    return res.status(404).send('User not found');
  }

  const currentBalance = (userResult.rows[0] as TransactionRecord).balance;

  const cobolFilePath = path.join(__dirname, 'cobol', 'ProcessTransaction.cbl');

  const userNameToCobol = user.toUpperCase();
  const transactionTypeToCobol = type.toUpperCase();

  // const command = `echo "${userNameToCobol}\n${currentBalance}\n${transaction}\n${transactionTypeToCobol}" | cobc -o cobol-output/ProcessTransaction -xj ${cobolFilePath}`;

  const outputDir = path.join(__dirname, 'cobol-output');

  if (!existsSync(outputDir)) {
    mkdirSync(outputDir, { recursive: true });
  }

  const cobolFileName = path.join(outputDir, 'ProcessTransaction');

  const command = `echo "${userNameToCobol}\n${currentBalance}\n${transaction}\n${transactionTypeToCobol}" | cobc -o ${cobolFileName} -xj ${cobolFilePath}`;

  try {
    const { stdout: cobolStdout } = await execAsync(command);
    console.log(`COBOL Program Output:\n\n${cobolStdout}`);

    const resultMatch = cobolStdout.match(/Result: (\d+\.\d+)/);

    if (resultMatch) {
      const newBalance = parseFloat(resultMatch[1]);
      //TODO: Update database logic goes here

      // 根據 type 調整 transaction 的正負值
      const adjustedTransactionAmount =
        type.toUpperCase() === 'WITHDRAW'
          ? -Math.abs(transaction)
          : Math.abs(transaction);

      // 構建 SQL 查詢字符串來插入新的交易記錄
      const INSERT_NEW_TRANSACTION = `
INSERT INTO transactions ("user", "transaction", "balance", "date")
VALUES ($1, $2, $3, NOW())
RETURNING *;
`;

      const transactionResult = await pgPool.query(INSERT_NEW_TRANSACTION, [
        user,
        adjustedTransactionAmount,
        newBalance,
      ]);
      const newTransactionRecord = transactionResult.rows[0];

      return res.json({
        status: 'success',
        cobolStdout,
        newTransactionRecord,
      });
    }
  } catch (error) {
    console.error(error);
    res.status(500).send('Server error');
  }
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
