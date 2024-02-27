export type TransactionRecord = {
  user: string;
  transaction: number;
  balance: number;
  date: string;
};

export type TransactionPayload = {
  user: string;
  transaction: number;
  type: 'Deposit' | 'Withdrawal';
};
