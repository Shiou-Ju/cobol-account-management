import express from 'express';
import { pgPool } from './connection';

const app = express();
const PORT = 3000;

app.get('/', async (_req, res) => {
  // res.send('Hello World!');
  const { rows } = await pgPool.query('SELECT NOW() as now');
  res.send(`Database time is: ${rows[0].now}`);
});

app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
