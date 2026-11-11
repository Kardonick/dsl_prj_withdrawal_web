/*
Source file: class.go
Módulo: Dominio del Negocio Bancario (Cuenta y Transacción con Auditoría JSON)
Curso: Lenguajes de Dominio Específico (DSL)
Fecha: nov-2026
Autor: J.Medianero.A (adaptado a Go) / Incremento US-04 (Integridad y Logging JSON)
*/
package headers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Constantes para los estados y tipos de transacción (US-04)
const (
	EstadoExitosa = "EXITOSA"
	EstadoFallo   = "FALLO"
	EstadoError   = "ERROR"

	TipoRetiro = "RETIRO"
	TipoAbono  = "ABONO"
)

// Transaccion representa el log estructurado de auditoría contable (US-04).
type Transaccion struct {
	ID             string  `json:"id"`
	FechaHora      string  `json:"fecha_hora"`
	Tipo           string  `json:"tipo"`
	Monto          float64 `json:"monto"`
	Estado         string  `json:"estado"`
	Mensaje        string  `json:"mensaje"`
	SaldoPosterior float64 `json:"saldo_posterior"`
	Titular        string  `json:"titular"`
}

// String formatea la transacción para depuración rápida por consola.
func (t Transaccion) String() string {
	return fmt.Sprintf("[%s] %s | %s: %.2f € | Estado: %s | Saldo Post: %.2f € | %s",
		t.ID, t.FechaHora, t.Tipo, t.Monto, t.Estado, t.SaldoPosterior, t.Mensaje)
}

// Cuenta representa la entidad de negocio descrita en el Modelo de Dominio BDD.
type Cuenta struct {
	nombre    string
	saldo     float64
	historial []Transaccion
}

// NewCuenta inicializa una nueva cuenta bancaria con saldo inicial.
func NewCuenta(nombre string, saldoInicial float64) *Cuenta {
	return &Cuenta{
		nombre:    nombre,
		saldo:     saldoInicial,
		historial: make([]Transaccion, 0),
	}
}

// ObtenerSaldo devuelve el saldo actual disponible.
func (c *Cuenta) ObtenerSaldo() float64 {
	return c.saldo
}

// Nombre devuelve el titular de la cuenta.
func (c *Cuenta) Nombre() string {
	return c.nombre
}

// Historial devuelve una copia de la lista de transacciones registradas.
func (c *Cuenta) Historial() []Transaccion {
	copia := make([]Transaccion, len(c.historial))
	copy(copia, c.historial)
	return copia
}

// UltimaTransaccion retorna un puntero a la última transacción generada.
func (c *Cuenta) UltimaTransaccion() (*Transaccion, error) {
	if len(c.historial) == 0 {
		return nil, errors.New("no existen transacciones registradas en el historial")
	}
	return &c.historial[len(c.historial)-1], nil
}

// String implementa fmt.Stringer para la cuenta.
func (c *Cuenta) String() string {
	return fmt.Sprintf("Cuenta{%s, saldo=%.2f, transacciones=%d}", c.nombre, c.saldo, len(c.historial))
}

// registrarTransaccion genera un registro de auditoría y lo anexa al historial interno.
func (c *Cuenta) registrarTransaccion(tipo string, monto float64, estado string, mensaje string) Transaccion {
	idTxn := fmt.Sprintf("TXN-%04d", len(c.historial)+1)
	fechaHora := time.Now().Format("2006-01-02 15:04:05")

	txn := Transaccion{
		ID:             idTxn,
		FechaHora:      fechaHora,
		Tipo:           tipo,
		Monto:          monto,
		Estado:         estado,
		Mensaje:        mensaje,
		SaldoPosterior: c.saldo,
		Titular:        c.nombre,
	}

	c.historial = append(c.historial, txn)
	return txn
}

/*
Retirar procesa una sustracción de fondos de forma atómica y auditable (US-01 / US-04).
- Si cantidad < 0: Registra ERROR y retorna error.
- Si cantidad > saldo: Registra FALLO y no modifica el saldo.
- Si es válido: Descuenta el monto, registra EXITOSA y actualiza el saldo.
*/
func (c *Cuenta) Retirar(cantidad float64) (bool, error) {
	if cantidad < 0 {
		c.registrarTransaccion(TipoRetiro, cantidad, EstadoError, "no se pueden retirar cantidades negativas")
		return false, errors.New("no se pueden retirar cantidades negativas")
	}

	if cantidad > c.saldo {
		c.registrarTransaccion(TipoRetiro, cantidad, EstadoFallo, "Saldo insuficiente para completar la operacion")
		return false, nil // Escenario: Withdrawal declined
	}

	c.saldo -= cantidad
	c.registrarTransaccion(TipoRetiro, cantidad, EstadoExitosa, "Retiro procesado correctamente")
	return true, nil // Escenario: Successful withdrawal
}

/*
Abonar procesa un ingreso de fondos de forma atómica y auditable (US-01b / US-04).
- Si cantidad <= 0: Registra ERROR y retorna error.
- Si es válido: Suma el monto, registra EXITOSA y actualiza el saldo.
*/
func (c *Cuenta) Abonar(cantidad float64) (bool, error) {
	if cantidad <= 0 {
		c.registrarTransaccion(TipoAbono, cantidad, EstadoError, "no se pueden abonar cantidades menores o iguales a cero")
		return false, errors.New("no se pueden abonar cantidades menores o iguales a cero")
	}

	c.saldo += cantidad
	c.registrarTransaccion(TipoAbono, cantidad, EstadoExitosa, "Abono procesado correctamente")
	return true, nil
}

// -------------------------------------------------------------
// Serialización y Persistencia en Formato JSON (US-04)
// -------------------------------------------------------------

// ExportarHistorialJSON serializa el historial completo a un JSON con indentación legible.
func (c *Cuenta) ExportarHistorialJSON() (string, error) {
	bytes, err := json.MarshalIndent(c.historial, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error serializando historial a JSON: %w", err)
	}
	return string(bytes), nil
}

// ExportarUltimaTransaccionJSON devuelve el JSON formateado de la última operación realizada.
func (c *Cuenta) ExportarUltimaTransaccionJSON() (string, error) {
	ultima, err := c.UltimaTransaccion()
	if err != nil {
		return "", err
	}
	bytes, err := json.MarshalIndent(ultima, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error serializando última transacción a JSON: %w", err)
	}
	return string(bytes), nil
}

// GuardarHistorialEnArchivo persiste el log de transacciones en un archivo físico en disco.
func (c *Cuenta) GuardarHistorialEnArchivo(nombreArchivo string) error {
	jsonStr, err := c.ExportarHistorialJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(nombreArchivo, []byte(jsonStr), 0644)
}

// LeerHistorialDesdeArchivo carga y deserializa transacciones desde un archivo JSON externo.
func LeerHistorialDesdeArchivo(nombreArchivo string) ([]Transaccion, error) {
	data, err := os.ReadFile(nombreArchivo)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo %s: %w", nombreArchivo, err)
	}

	var transacciones []Transaccion
	if err := json.Unmarshal(data, &transacciones); err != nil {
		return nil, fmt.Errorf("error deserializando contenido JSON: %w", err)
	}

	return transacciones, nil
}
