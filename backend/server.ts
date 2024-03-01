import express from 'express';
import { pgPool } from './connection';
import { exec } from 'child_process';
import {
  TransactionPayload,
  TransactionRecord,
  TransactionType,
  VALID_TRANSACTION_TYPES,
} from './interfaces';
import { promisify } from 'util';
import path from 'path';
import { existsSync, mkdirSync } from 'fs';

function isValidTransactionType(
  type: TransactionType,
): type is TransactionType {
  return VALID_TRANSACTION_TYPES.includes(type);
}

const execAsync = promisify(exec);

const app = express();
const PORT = 3000;

app.use(express.json());

const isDevMode = process.env.NODE_ENV === 'development';

app.get('/api/time', async (_req, res) => {
  const { rows } = await pgPool.query('SELECT NOW() as now');
  res.send(`Database time is: ${rows[0].now}`);
});

app.get('/api/user/:user', async (req, res) => {
  const userName = req.params.user;

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

app.get('/api/user/:user/transactions', async (req, res) => {
  const userName = req.params.user;

  const GET_ALL_USER_TRANSACTIONS = `
    SELECT * FROM transactions
    WHERE "user" = $1
    ORDER BY "date" DESC
    LIMIT 10;
  `;

  try {
    const { rows } = await pgPool.query(GET_ALL_USER_TRANSACTIONS, [userName]);
    const hasTransactions = rows.length > 0;

    if (!hasTransactions) {
      return res
        .status(404)
        .send(`Transactions not found for the user "${userName}"`);
    }

    res.json(rows);
  } catch (error) {
    console.error('Error fetching transactions:', error);
    res.status(500).send('Server error');
  }
});

app.get('/api/users', async (_req, res) => {
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

app.post('/api/transaction', async (req, res) => {
  const { user, transaction, type } = req.body as TransactionPayload;

  const hasMissingData = !user || transaction === undefined || !type;

  if (hasMissingData) {
    return res.status(400).send('Missing required fields');
  }

  if (!isValidTransactionType(type)) {
    return res.status(400).send('invalid transaction type');
  }

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

  const outputDir = path.join(__dirname, 'cobol-output');

  if (!existsSync(outputDir)) {
    mkdirSync(outputDir, { recursive: true });
  }

  const cobolFileName = path.join(outputDir, 'ProcessTransaction');

  try {
    const { stdout: versionStdout } = await execAsync('cobc --version');
    console.log(`COBOL Compiler Version:\n\n${versionStdout}`);
  } catch (versionError) {
    console.error(`Error getting COBOL compiler version: ${versionError}`);
    return;
  }

  const command = `echo "${userNameToCobol}\n${currentBalance}\n${transaction}\n${transactionTypeToCobol}" | cobc -o ${cobolFileName} -xj ${cobolFilePath}`;

  try {
    const { stdout: cobolStdout } = await execAsync(command);
    console.log(`COBOL Program Output:\n\n${cobolStdout}`);

    const resultMatch = cobolStdout.match(/Result: (\d+\.\d+)/);

    if (resultMatch) {
      const newBalance = parseFloat(resultMatch[1]);

      const adjustedTransactionAmount =
        type.toUpperCase() === 'WITHDRAW'
          ? -Math.abs(transaction)
          : Math.abs(transaction);

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

if (!isDevMode) {
  const frontendDistPath = path.resolve(__dirname, '../../frontend/dist');

  app.use(express.static(frontendDistPath));

  app.get('*', (_req, res) => {
    const frontEndPageToBeServed = path.join(frontendDistPath, 'index.html');

    res.sendFile(frontEndPageToBeServed);
  });
}

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
