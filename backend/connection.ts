import { Pool } from 'pg';

export const pgPool = new Pool({
  user: 'postgres',
  host: 'localhost',
  database: 'cobolexample',
  password: 'cobolexamplepw',
  port: 5432,
});
