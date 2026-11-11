Feature: Transaction Integrity and Audit Logging (US-04)
  As the banking core service
  I want every account transaction to be atomic and logged in JSON format
  So that accounting consistency and financial audit trails are guaranteed

  Rule: Every operation must be atomic and categorized as EXITOSA, FALLO, or ERROR

    Scenario: Log de transacción exitosa por retiro dentro del saldo
      Given Alice tiene una cuenta con un saldo inicial de 234.56
      When Alice solicita un retiro de 200.00
      Then el retiro se realiza con éxito
      And el saldo final de la cuenta es 34.56
      And el historial registra una transacción con el siguiente JSON:
        """json
        {
          "tipo": "RETIRO",
          "monto": 200.00,
          "estado": "EXITOSA",
          "mensaje": "Retiro procesado correctamente",
          "saldo_posterior": 34.56,
          "titular": "Alice"
        }
        """

    Scenario: Log de transacción fallida por saldo insuficiente
      Given Hamza tiene una cuenta con un saldo inicial de 198.76
      When Hamza solicita un retiro de 200.00
      Then el retiro es rechazado por fondos insuficientes
      And el saldo de la cuenta permanece inalterado en 198.76
      And el historial registra una transacción con el siguiente JSON:
        """json
        {
          "tipo": "RETIRO",
          "monto": 200.00,
          "estado": "FALLO",
          "mensaje": "Saldo insuficiente para completar la operacion",
          "saldo_posterior": 198.76,
          "titular": "Hamza"
        }
        """

    Scenario: Log de transacción rechazada por error de validación (monto negativo)
      Given Alice tiene una cuenta con un saldo inicial de 234.56
      When Alice intenta retirar una cantidad invalida de -50.00
      Then el sistema captura un error de validacion
      And el saldo de la cuenta permanece inalterado en 234.56
      And el historial registra una transacción con el siguiente JSON:
        """json
        {
          "tipo": "RETIRO",
          "monto": -50.00,
          "estado": "ERROR",
          "mensaje": "no se pueden retirar cantidades negativas",
          "saldo_posterior": 234.56,
          "titular": "Alice"
        }
        """

    Scenario: Log de transacción exitosa por abono de efectivo
      Given Alice tiene una cuenta con un saldo inicial de 234.56
      When Alice solicita un abono de 100.00
      Then el abono se realiza con éxito
      And el saldo final de la cuenta es 334.56
      And el historial registra una transacción con el siguiente JSON:
        """json
        {
          "tipo": "ABONO",
          "monto": 100.00,
          "estado": "EXITOSA",
          "mensaje": "Abono procesado correctamente",
          "saldo_posterior": 334.56,
          "titular": "Alice"
        }
        """

    Scenario: Log de transacción rechazada por error de abono inválido (monto cero o negativo)
      Given Alice tiene una cuenta con un saldo inicial de 234.56
      When Alice intenta abonar una cantidad invalida de 0.00
      Then el sistema captura un error de validacion
      And el saldo de la cuenta permanece inalterado en 234.56
      And el historial registra una transacción con el siguiente JSON:
        """json
        {
          "tipo": "ABONO",
          "monto": 0.00,
          "estado": "ERROR",
          "mensaje": "no se pueden abonar cantidades menores o iguales a cero",
          "saldo_posterior": 234.56,
          "titular": "Alice"
        }
        """