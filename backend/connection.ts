import { Pool } from 'pg';

const isDevMode = process.env.NODE_ENV === 'development';

const host = isDevMode ? 'localhost' : 'db';

export const pgPool = new Pool({
  user: 'postgres',
  host: host,
  database: 'cobolexample',
  password: 'cobolexamplepw',
  port: 5432,
});
