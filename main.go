/*
Source file: main.go [Servidor Web HTTP Transaccional - Acápite 4]
Curso: Lenguajes de Dominio Específico (DSL)
Fecha: nov-2026
Autor: J.Medianero.A / Refactorización Web
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"prj_withdrawal/handlers"
	"prj_withdrawal/headers"
)

func main() {
	// Banderas de configuración del servidor
	puerto := flag.String("port", "8080", "Puerto de escucha del servidor HTTP")
	archivoLog := flag.String("archivo", "transacciones.json", "Ruta del archivo JSON de persistencia")
	titular := flag.String("nombre", "Alice", "Nombre del titular de la cuenta inicial")
	saldoInicial := flag.Float64("saldo", 234.56, "Saldo inicial de la cuenta")

	flag.Parse()

	// 1. Inicialización de la entidad de negocio (Dominio)
	cuenta := headers.NewCuenta(*titular, *saldoInicial)

	// 2. Inicialización del controlador HTTP con sincronización de concurrencia
	bankHandler := handlers.NewBankHandler(cuenta, *archivoLog)

	// 3. Registro de endpoints de la API RESTful
	http.HandleFunc("/api/cuenta", bankHandler.GetCuenta)
	http.HandleFunc("/api/retirar", bankHandler.Retirar)
	http.HandleFunc("/api/abonar", bankHandler.Abonar)
	http.HandleFunc("/api/transacciones", bankHandler.GetTransacciones)

	// 4. Servidor de archivos estáticos para la interfaz web (static/)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	direccion := ":" + *puerto
	fmt.Println("======================================================================")
	fmt.Println("       BANCO FIIS - SERVICIO TRANSACCIONAL WEB (ACÁPITE 4)            ")
	fmt.Println("======================================================================")
	fmt.Printf("Servidor HTTP ejecutándose en: http://localhost:%s\n", *puerto)
	fmt.Printf("Titular: %s | Saldo: %.2f € | Persistencia: %s\n", *titular, *saldoInicial, *archivoLog)
	fmt.Println("Endpoints API disponibles:")
	fmt.Printf("  • GET  http://localhost:%s/api/cuenta\n", *puerto)
	fmt.Printf("  • POST http://localhost:%s/api/retirar\n", *puerto)
	fmt.Printf("  • POST http://localhost:%s/api/abonar\n", *puerto)
	fmt.Printf("  • GET  http://localhost:%s/api/transacciones\n", *puerto)
	fmt.Println("Presiona Ctrl + C para finalizar el proceso.")

	if err := http.ListenAndServe(direccion, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
	}
}
