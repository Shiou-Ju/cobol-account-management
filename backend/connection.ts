import { Pool } from 'pg';
// 引入 pg Pool

// 創建一個新的 Pool 實例來管理您的 PostgreSQL 連接
export const pgPool = new Pool({
  user: 'postgres',
  // TODO: might be probelmm, name server does not recognize
  host: 'localhost',
  database: 'cobolexample',
  password: 'cobolexamplepw',
  port: 5432,
});
