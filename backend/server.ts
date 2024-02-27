import express from 'express';
import { pgPool } from './connection';
// TODO: not unser root dir
// File '/Users/bamboo/Repos/cobol-account-management/utils/interfaces.ts' is not under 'rootDir' '/Users/bamboo/Repos/cobol-account-management/backend'. 'rootDir' is expected to contain all source files.
// The file is in the program because:
// Matched by include pattern '../utils/**/*' in '/Users/bamboo/Repos/cobol-account-management/backend/tsconfig.json'ts

import { TransactionPayload } from '../utils/interfaces';

const app = express();
const PORT = 3000;

app.get('/', async (_req, res) => {
  const { rows } = await pgPool.query('SELECT NOW() as now');
  res.send(`Database time is: ${rows[0].now}`);
});

app.get('/user/:user', async (req, res) => {
  const userName = req.params.user;

  try {
    const { rows } = await pgPool.query(
      'SELECT * FROM transactions WHERE "user" = $1;',
      [userName],
    );

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
    const { rows } = await pgPool.query('SELECT * FROM transactions;');
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

  try {
    // TODO: update using cobol program

    res.send('Transaction processed successfully');
  } catch (error) {
    console.error(error);
    res.status(500).send('Server error');
  }
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
