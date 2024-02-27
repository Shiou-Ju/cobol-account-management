IDENTIFICATION DIVISION.
PROGRAM-ID. ProcessTransaction.

DATA DIVISION.
WORKING-STORAGE SECTION.
01 TransactionData.
   05 T-User         PIC X(20).
   05 T-CurrentBalance PIC 9(10)V99.
   05 T-Transaction  PIC 9(10)V99.
   05 T-Type         PIC X(10).
   05 T-Result       PIC 9(10)V99.

PROCEDURE DIVISION.
   ACCEPT T-User FROM CONSOLE.
   ACCEPT T-CurrentBalance FROM CONSOLE.
   ACCEPT T-Transaction FROM CONSOLE.
   ACCEPT T-Type FROM CONSOLE.

   DISPLAY "User: " T-User.
   DISPLAY "Current Balance: " T-CurrentBalance.
   DISPLAY "Transaction: " T-Transaction.
   DISPLAY "Type: " T-Type.
   
   EVALUATE T-Type
       WHEN "DEPOSIT" 
           COMPUTE T-Result = T-CurrentBalance + T-Transaction
           DISPLAY "Deposit block entered."
       WHEN "WITHDRAW" 
           COMPUTE T-Result = T-CurrentBalance - T-Transaction
           DISPLAY "Withdraw block entered."
   END-EVALUATE.
   
   DISPLAY "Result: " T-Result.
   
   STOP RUN.
