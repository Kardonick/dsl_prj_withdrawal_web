/*
Source file: test_bdd.go
Módulo: Harness de Pruebas BDD / TDD con Validación de Integridad y Logs JSON (US-04)
Curso: Lenguajes de Dominio Específico (DSL)
Fecha: nov-2026
Autor: J.Medianero.A / Incremento US-04 (Auditoría Contable y JSON)
*/
package headers

import (
	"encoding/json"
	"fmt"
	"math"
)

// imprimirSaldo imprime el saldo formateado a 2 decimales.
func imprimirSaldo(cuenta *Cuenta) {
	fmt.Printf("Saldo actual: %.2f €\n", cuenta.ObtenerSaldo())
}

// validarYMostrarJSON extrae la última transacción, verifica su validez sintáctica e imprime el JSON.
func validarYMostrarJSON(cuenta *Cuenta, estadoEsperado string) bool {
	txn, err := cuenta.UltimaTransaccion()
	if err != nil {
		fmt.Printf("ERROR BDD: No se encontró registro de transacción en el historial: %v\n", err)
		return false
	}

	jsonStr, err := cuenta.ExportarUltimaTransaccionJSON()
	if err != nil {
		fmt.Printf("ERROR BDD: Falló la serialización a JSON: %v\n", err)
		return false
	}

	// Verificación de integridad sintáctica del JSON
	if !json.Valid([]byte(jsonStr)) {
		fmt.Println("ERROR BDD: El JSON generado es sintácticamente inválido.")
		return false
	}

	fmt.Println("Log JSON generado (Auditoría US-04):")
	fmt.Println(jsonStr)

	// Verificación del estado contable
	if txn.Estado != estadoEsperado {
		fmt.Printf("ERROR BDD: Estado esperado '%s', pero se obtuvo '%s'.\n", estadoEsperado, txn.Estado)
		return false
	}

	fmt.Printf("Validación US-04: Log JSON válido con estado '%s' confirmado.\n", txn.Estado)
	return true
}

// -------------------------------------------------------------
// Escenarios BDD: Retiro de Efectivo (US-01, US-02, US-03, US-04)
// -------------------------------------------------------------

// TestEscenario1Exito simula un retiro exitoso, verifica saldo y log JSON (Estado: EXITOSA).
func TestEscenario1Exito(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("\n--- Escenario 1: Successful withdrawal within balance (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if exito {
		fmt.Println("Resultado: ÉXITO. La retirada fue procesada.")
		imprimirSaldo(cuenta)

		saldoEsperado := saldoInicial - cantidadRetirar
		if math.Abs(cuenta.ObtenerSaldo()-saldoEsperado) < 0.001 {
			fmt.Printf("Validación: El saldo restante es correcto (%.2f).\n", saldoEsperado)
		} else {
			fmt.Println("ERROR: El saldo restante no coincide.")
		}

		validarYMostrarJSON(cuenta, EstadoExitosa)
	} else {
		fmt.Println("Resultado: FALLO. La retirada fue rechazada incorrectamente.")
	}
}

// TestEscenario2Fallo simula un retiro rechazado por fondos insuficientes y log JSON (Estado: FALLO).
func TestEscenario2Fallo(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("\n--- Escenario 2: Declined withdrawal in excess of balance (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if !exito {
		fmt.Println("Resultado: ÉXITO. La retirada fue rechazada correctamente por saldo insuficiente.")
		imprimirSaldo(cuenta)

		// Verificación de atomicidad: el saldo debe permanecer inalterado
		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Validación: El saldo permanece inalterado (%.2f).\n", saldoInicial)
		} else {
			fmt.Println("ERROR: El saldo se modificó indebidamente.")
		}

		validarYMostrarJSON(cuenta, EstadoFallo)
	} else {
		fmt.Println("Resultado: FALLO. La transacción debió ser rechazada.")
	}
}

// TestEscenario3Error simula una entrada inválida y log JSON (Estado: ERROR).
func TestEscenario3Error(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("\n--- Escenario 3: Validation Error - Negative Amount (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: ÉXITO. El sistema atrapó el error correctamente: %v\n", err)
		imprimirSaldo(cuenta)

		// Verificación de atomicidad
		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Validación: El saldo no sufrió alteraciones (%.2f).\n", saldoInicial)
		} else {
			fmt.Println("ERROR: El saldo fue alterado tras un error.")
		}

		validarYMostrarJSON(cuenta, EstadoError)
	} else if !exito {
		fmt.Println("Resultado: RECHAZADO. Fondos insuficientes (no es un error de validación).")
	} else {
		fmt.Println("Resultado: FALLO. Se procesó una operación con monto negativo.")
	}
}

// -------------------------------------------------------------
// Escenarios BDD: Abono de Efectivo (US-01b, US-04)
// -------------------------------------------------------------

// TestEscenarioAbonoExito simula un abono válido y log JSON (Estado: EXITOSA).
func TestEscenarioAbonoExito(nombre string, saldoInicial float64, cantidadAbonar float64) {
	fmt.Printf("\n--- Escenario Abono 1: Successful deposit (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta abonar %.2f\n", nombre, cantidadAbonar)
	exito, err := cuenta.Abonar(cantidadAbonar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if exito {
		fmt.Println("Resultado: ÉXITO. El abono fue procesado.")
		imprimirSaldo(cuenta)

		saldoEsperado := saldoInicial + cantidadAbonar
		if math.Abs(cuenta.ObtenerSaldo()-saldoEsperado) < 0.001 {
			fmt.Printf("Validación: El saldo resultante es correcto (%.2f).\n", saldoEsperado)
		} else {
			fmt.Println("ERROR: El saldo resultante no coincide.")
		}

		validarYMostrarJSON(cuenta, EstadoExitosa)
	} else {
		fmt.Println("Resultado: FALLO. El abono no pudo completarse.")
	}
}

// TestEscenarioAbonoError simula un abono con monto <= 0 y log JSON (Estado: ERROR).
func TestEscenarioAbonoError(nombre string, saldoInicial float64, cantidadAbonar float64) {
	fmt.Printf("\n--- Escenario Abono 2: Invalid deposit amount (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta abonar %.2f\n", nombre, cantidadAbonar)
	exito, err := cuenta.Abonar(cantidadAbonar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: ÉXITO. El sistema atrapó el error correctamente: %v\n", err)
		imprimirSaldo(cuenta)

		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Validación: El saldo permanece inalterado (%.2f).\n", saldoInicial)
		} else {
			fmt.Println("ERROR: El saldo fue alterado indebidamente.")
		}

		validarYMostrarJSON(cuenta, EstadoError)
	} else if exito {
		fmt.Println("Resultado: FALLO. Se procesó un abono con monto inválido.")
	}
}

// -------------------------------------------------------------
// Suite Completa de Integridad US-04: Transacciones + Historial JSON
// -------------------------------------------------------------

// TestSuiteIntegridadUS04 ejecuta un ciclo completo de operaciones sobre una misma cuenta,
// validando la generación acumulativa del historial en formato JSON y su persistencia a archivo.
func TestSuiteIntegridadUS04(nombre string, saldoInicial float64) {
	fmt.Printf("\n============================================================\n")
	fmt.Printf("=== SUITE DE PRUEBA: INTEGRIDAD US-04 + AUDITORÍA JSON ===\n")
	fmt.Printf("============================================================\n")

	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Cuenta creada: %s con saldo inicial %.2f €\n", nombre, saldoInicial)

	// 1. Operación EXITOSA (Abono)
	fmt.Println("\n[Paso 1] Ejecutando Abono de 100.00 € (Esperado: EXITOSA)")
	cuenta.Abonar(100.00)

	// 2. Operación FALLO (Retiro que excede saldo)
	fmt.Println("\n[Paso 2] Intentando retirar 500.00 € (Esperado: FALLO por saldo insuficiente)")
	cuenta.Retirar(500.00)

	// 3. Operación ERROR (Retiro con monto negativo)
	fmt.Println("\n[Paso 3] Intentando retirar -30.00 € (Esperado: ERROR de validación)")
	cuenta.Retirar(-30.00)

	// 4. Operación EXITOSA (Retiro dentro de saldo)
	fmt.Println("\n[Paso 4] Ejecutando Retiro de 50.00 € (Esperado: EXITOSA)")
	cuenta.Retirar(50.00)

	// Validación global del historial acumulado
	fmt.Println("\n------------------------------------------------------------")
	fmt.Println("HISTORIAL COMPLETO DE AUDITORÍA EN FORMATO JSON (ARRAY):")
	fmt.Println("------------------------------------------------------------")

	jsonCompleto, err := cuenta.ExportarHistorialJSON()
	if err != nil {
		fmt.Printf("ERROR: Fallo al exportar el array JSON: %v\n", err)
		return
	}
	fmt.Println(jsonCompleto)

	// Verificación de integridad contable final
	// Saldo esperado: 234.56 + 100.00 (abono) - 50.00 (retiro) = 284.56 €
	saldoEsperadoFinal := saldoInicial + 100.00 - 50.00
	if math.Abs(cuenta.ObtenerSaldo()-saldoEsperadoFinal) < 0.001 {
		fmt.Printf("\nValidación Final de Atomicidad: Saldo consistente en %.2f €\n", cuenta.ObtenerSaldo())
	} else {
		fmt.Println("\nERROR: Inconsistencia contable en el saldo final.")
	}

	// Prueba de guardado y lectura de archivo JSON
	archivoLog := "transacciones.json"
	if err := cuenta.GuardarHistorialEnArchivo(archivoLog); err != nil {
		fmt.Printf("ERROR guardando en %s: %v\n", archivoLog, err)
	} else {
		fmt.Printf("Persistencia: Archivo '%s' guardado exitosamente en disco.\n", archivoLog)

		// Verificación de lectura
		txsLeidas, err := LeerHistorialDesdeArchivo(archivoLog)
		if err != nil {
			fmt.Printf("ERROR leyendo desde %s: %v\n", archivoLog, err)
		} else {
			fmt.Printf("Lectura confirmada: %d transacciones cargadas desde '%s'.\n", len(txsLeidas), archivoLog)
		}
	}
	fmt.Println("============================================================")
}
