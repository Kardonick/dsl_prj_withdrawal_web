Feature: Cash deposit (Abono de efectivo)

  Rule: Customers can only deposit positive amounts greater than zero

    Scenario: Successful cash deposit
      Given Alice has 234.56 in their account
      When Alice deposits 100.00
      Then the deposit is successful
      And the account balance should be 334.56

    Scenario: Declined deposit with negative amount
      Given Alice has 234.56 in their account
      When Alice tries to deposit -50.00
      Then the deposit fails with a validation error
      And the account balance remains 234.56

    Scenario: Declined deposit with zero amount
      Given Alice has 234.56 in their account
      When Alice tries to deposit 0.00
      Then the deposit fails with a validation error
      And the account balance remains 234.56