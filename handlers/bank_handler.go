/*
Source file: handlers/bank_handler.go
Módulo: Controladores HTTP y API RESTful (Acápite 4 - Refactorización Web)
Curso: Metodologías de Desarrollo de Software / Arquitectura de Software
Fecha: nov-2026
Autor: J.Medianero.A / Refactorización Web
*/
package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"prj_withdrawal/headers"
)

// OperacionRequest define el payload JSON recibido para retiros y abonos.
type OperacionRequest struct {
	Monto float64 `json:"monto"`
}

// CuentaResponse define la respuesta con el estado actual de la cuenta.
type CuentaResponse struct {
	Titular string  `json:"titular"`
	Saldo   float64 `json:"saldo"`
}

// BankHandler encapsula la cuenta, el cerrojo de sincronización y el archivo de persistencia.
type BankHandler struct {
	mu         sync.Mutex
	cuenta     *headers.Cuenta
	archivoLog string
}

// NewBankHandler inicializa el controlador con su respectiva cuenta y ruta de log.
func NewBankHandler(cuenta *headers.Cuenta, archivoLog string) *BankHandler {
	return &BankHandler{
		cuenta:     cuenta,
		archivoLog: archivoLog,
	}
}

// responderJSON centraliza la serialización y cabeceras HTTP de respuesta.
func responderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// GetCuenta maneja: GET /api/cuenta
func (h *BankHandler) GetCuenta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	res := CuentaResponse{
		Titular: h.cuenta.Nombre(),
		Saldo:   h.cuenta.ObtenerSaldo(),
	}
	responderJSON(w, http.StatusOK, res)
}

// Retirar maneja: POST /api/retirar
func (h *BankHandler) Retirar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req OperacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "El cuerpo JSON de la solicitud es inválido",
		})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Procesa la lógica de negocio pura (US-01 / US-04)
	h.cuenta.Retirar(req.Monto)
	_ = h.cuenta.GuardarHistorialEnArchivo(h.archivoLog)

	txn, err := h.cuenta.UltimaTransaccion()
	if err != nil {
		responderJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Error interno al recuperar la transacción",
		})
		return
	}

	responderJSON(w, http.StatusOK, txn)
}

// Abonar maneja: POST /api/abonar
func (h *BankHandler) Abonar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req OperacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "El cuerpo JSON de la solicitud es inválido",
		})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Procesa la lógica de negocio pura (US-01b / US-04)
	h.cuenta.Abonar(req.Monto)
	_ = h.cuenta.GuardarHistorialEnArchivo(h.archivoLog)

	txn, err := h.cuenta.UltimaTransaccion()
	if err != nil {
		responderJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Error interno al recuperar la transacción",
		})
		return
	}

	responderJSON(w, http.StatusOK, txn)
}

// GetTransacciones maneja: GET /api/transacciones
func (h *BankHandler) GetTransacciones(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	txs := h.cuenta.Historial()
	responderJSON(w, http.StatusOK, txs)
}
