import express from 'express';
import { pgPool } from './connection';

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

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
