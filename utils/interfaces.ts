export type TransactionRecord = {
  user: string;
  transaction: number;
  balance: number;
  date: string;
};

export const VALID_TRANSACTION_TYPES = ['DEPOSIT', 'WITHDRAW'] as const;

export type TransactionType = typeof VALID_TRANSACTION_TYPES[number];

export type TransactionPayload = {
  user: string;
  transaction: number;
  type: TransactionType;
};
