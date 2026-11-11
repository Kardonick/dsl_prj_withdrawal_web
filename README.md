# Core Bancario - prj_withdrawal

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Methodology](https://img.shields.io/badge/Methodology-BDD%20%7C%20TDD-blue)](https://cucumber.io)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2F%20Layered-orange)](#arquitectura-del-sistema)
[![Status](https://img.shields.io/badge/Status-Completed%20(Ac%C3%A1pite%204)-brightgreen)](#roadmap-incremental)

Sistema bancario transaccional desarrollado en **Golang** bajo un enfoque riguroso de **Desarrollo Guiado por el Comportamiento (BDD)** articulado con **Desarrollo Guiado por Pruebas (TDD)**. El proyecto modela operaciones atómicas de retiro y abono en cuentas bancarias, genera trazas estructuradas de auditoría en formato JSON y desacopla la lógica de dominio hacia una **API RESTful concurrente** con un panel de control web responsivo.

---

## Información Institucional y Académica

* **Institución:** Universidad Nacional de Ingeniería (UNI)
* **Facultad:** Facultad de Ingeniería Industrial y de Sistemas (FIIS)
* **Escuela:** Escuela Profesional de Ingeniería de Software
* **Asignatura:** Lenguajes de Dominio Específico (DSL)
* **Catedrático:** Ing. Jorge Medianero Acosta
* **Autor:** Dair David Ramos Jacay
* **Código de Estudiante:** `20232156D`

---

## Roadmap Incremental

El proyecto fue concebido y ejecutado bajo un ciclo de vida estrictamente iterativo e incremental:

```text
[ Fase Base: codeLAB-02 ]
   └─ US-01 a US-03: Retiro de efectivo, validación de saldos y manejo de montos negativos.
[ Acápite 2: Primera Entrega ]
   └─ US-01b: Abono en efectivo con control de montos estrictamente positivos ($monto > 0$).
[ Acápite 3: Segunda Entrega ]
   └─ US-04: Integridad de transacción (Atomicidad contable) y logging estructurado JSON (EXITOSA, FALLO, ERROR).
[ Acápite 4: Tercera Entrega (Hito Actual) ]
   └─ Refactorización Web: API RESTful en Go (net/http, sync.Mutex) + SPA Frontend Dashboard.
```

---

## Arquitectura del Sistema

La solución adopta una **Arquitectura en Capas (Layered Architecture)** con separación estricta de responsabilidades (*Separation of Concerns*) e independencia de frameworks pesados:

```text
                       [ Navegador Web / Cliente HTTP ]
                                      │
                                      ▼
                        ┌──────────────────────────┐
                        │   handlers/ (API REST)   │
                        │   Controladores HTTP     │
                        │   sync.Mutex Concurrente │
                        └─────────────┬────────────┘
                                      │ Invoca
                                      ▼
                        ┌──────────────────────────┐
                        │   headers/ (Dominio)     │
                        │   Entidad Cuenta         │
                        │   Entidad Transaccion    │
                        │   Reglas de Negocio BDD  │
                        └─────────────┬────────────┘
                                      │ Serializa / Lee
                                      ▼
                        ┌──────────────────────────┐
                        │    transacciones.json    │
                        │    Persistencia en Disco │
                        └──────────────────────────┘
```

1. **Capa de Dominio (`headers/`):** Contiene la lógica nuclear bancaria inmutable (`Cuenta` y `Transaccion`), garantizando la regla de *todo o nada* (atomicidad) y el cálculo de balances.
2. **Capa de Pruebas BDD (`headers/test_bdd.go`):** Emula los pasos declarativos Given-When-Then de las especificaciones Gherkin, validando integridad sintáctica (`json.Valid`) y aritmética flotante IEEE 754 mediante tolerancia absoluta ($\varepsilon = 0.001$).
3. **Capa de Controladores HTTP (`handlers/`):** Expone los endpoints de la API, deserializa los payloads JSON, gestiona bloqueos de concurrencia mediante `sync.Mutex` y orquesta la persistencia.
4. **Capa de Presentación Web (`static/`):** Dashboard interactivo consolidado en `index.html` con HTML5 semántico y lógica cliente asíncrona (Fetch API), estilizado mediante variables y componentes adaptables en `styles.css`.

---

## Estructura del Repositorio

```text
prj_withdrawal/
├── go.mod                     # Declaración del módulo Go (v1.21+)
├── main.go                    # Servidor HTTP y configuración de enrutamiento
├── transacciones.json         # Fuente de persistencia transaccional en disco
├── caracteristicas/           # Documentación viva BDD (Gherkin Features)
│   ├── retiro-cash.feature    # Escenarios de retiro (US-01 a US-03)
│   ├── abono-cash.feature     # Escenarios de abono (US-01b)
│   └── transaccion-log.feature# Escenarios de integridad y log JSON (US-04)
├── handlers/                  # Controladores RESTful concurrentes
│   └── bank_handler.go        # Handlers HTTP para Cuenta, Retiro, Abono y Logs
├── headers/                   # Núcleo de Dominio y Harness de pruebas
│   ├── class.go               # Modelos Cuenta y Transacción con serializador JSON
│   └── test_bdd.go            # Arnés BDD ejecutable y validadores de esquema
└── static/                    # Frontend de la aplicación Web
    ├── index.html             # Dashboard, formularios de operaciones y lógica asíncrona (Fetch API)
    └── styles.css             # Hoja de estilos moderna y adaptable
```

---

## Especificación BDD y Reglas de Negocio (US-04)

Toda operación bancaria evaluada por el sistema clasifica de forma unívoca en uno de los tres estados auditables:

| Estado | Condición de Activación | Impacto en Saldo | Mensaje de Auditoría |
| :---: | :--- | :---: | :--- |
| **`EXITOSA`** | Precondiciones cumplidas ($monto \le saldo$ en retiro o $monto > 0$ en abono). | Modificado ($\pm monto$) | Operación procesada correctamente |
| **`FALLO`** | Fondos insuficientes en solicitud de retiro ($monto > saldo$). | Inalterado | Saldo insuficiente para completar la operación |
| **`ERROR`** | Datos de entrada inválidos ($monto \le 0$ en abono o $monto < 0$ en retiro). | Inalterado | Error de validación según el tipo |

### Ejemplo de Log JSON Estructurado

```json
{
  "id": "TXN-0001",
  "fecha_hora": "2026-11-10 02:00:15",
  "tipo": "RETIRO",
  "monto": 200.00,
  "estado": "EXITOSA",
  "mensaje": "Retiro procesado correctamente",
  "saldo_posterior": 34.56,
  "titular": "Alice"
}
```

---

## Documentación de la API RESTful

La API se comunica exclusivamente en formato JSON (`Content-Type: application/json`).

### 1. Consultar Estado de la Cuenta
* **Método:** `GET`
* **Ruta:** `/api/cuenta`
* **Respuesta Exitosa (`200 OK`):**
  ```json
  {
    "titular": "Alice",
    "saldo": 234.56
  }
  ```

### 2. Procesar Retiro de Efectivo
* **Método:** `POST`
* **Ruta:** `/api/retirar`
* **Cuerpo de Solicitud:**
  ```json
  {
    "monto": 50.00
  }
  ```
* **Respuesta Exitosa (`200 OK`):** Retorna el objeto `Transaccion` con el resultado y estado (`EXITOSA`, `FALLO` o `ERROR`).

### 3. Procesar Abono de Efectivo
* **Método:** `POST`
* **Ruta:** `/api/abonar`
* **Cuerpo de Solicitud:**
  ```json
  {
    "monto": 100.00
  }
  ```
* **Respuesta Exitosa (`200 OK`):** Retorna el objeto `Transaccion` con el nuevo balance acreditado.

### 4. Consultar Historial de Auditoría
* **Método:** `GET`
* **Ruta:** `/api/transacciones`
* **Respuesta Exitosa (`200 OK`):** Array con todas las transacciones históricas registradas en `transacciones.json`.

---

## Instalación y Puesta en Marcha

### Prerrequisitos
* Tener instalado [Go 1.21+](https://go.dev/dl/) o superior en tu sistema operativo.
* Un navegador web moderno (Chrome, Firefox, Edge, Safari).

### Pasos de Ejecución

1. **Clonar el repositorio:**
   ```bash
   git clone https://github.com/TU_USUARIO/prj_withdrawal.git
   cd prj_withdrawal
   ```

2. **Compilar y ejecutar el servidor web:**
   ```bash
   go run main.go
   ```
   *(Opcional: puedes personalizar el puerto, titular o saldo inicial mediante banderas)*
   ```bash
   go run main.go -port 8080 -nombre "Alice" -saldo 234.56 -archivo transacciones.json
   ```

3. **Acceder a la aplicación:**
   Abre tu navegador e ingresa a:
   ```text
   http://localhost:8080
   ```

---

## Pruebas Automatizadas desde CLI

Para validar los contratos de la API directamente desde la consola sin abrir el navegador, puedes ejecutar las siguientes pruebas según tu sistema operativo:

### En Windows (PowerShell con `curl.exe --%`)

El operador `--%` (*stop-parsing*) evita que PowerShell interprete o modifique las comillas del payload JSON:

```bash
# 1. Validar abono válido (Estado: EXITOSA)
curl.exe --% -X POST http://localhost:8080/api/abonar -H "Content-Type: application/json" -d "{\"monto\": 100.0}"

# 2. Validar rechazo por saldo insuficiente (Estado: FALLO)
curl.exe --% -X POST http://localhost:8080/api/retirar -H "Content-Type: application/json" -d "{\"monto\": 5000.0}"

# 3. Validar captura de monto negativo (Estado: ERROR)
curl.exe --% -X POST http://localhost:8080/api/retirar -H "Content-Type: application/json" -d "{\"monto\": -40.0}"

# 4. Consultar historial acumulado de auditoría
curl.exe -X GET http://localhost:8080/api/transacciones
```

### En Linux / macOS / Git Bash

```bash
# 1. Validar abono válido (Estado: EXITOSA)
curl -X POST http://localhost:8080/api/abonar \
     -H "Content-Type: application/json" \
     -d '{"monto": 100.0}'

# 2. Validar rechazo por saldo insuficiente (Estado: FALLO)
curl -X POST http://localhost:8080/api/retirar \
     -H "Content-Type: application/json" \
     -d '{"monto": 5000.0}'

# 3. Validar captura de monto negativo (Estado: ERROR)
curl -X POST http://localhost:8080/api/retirar \
     -H "Content-Type: application/json" \
     -d '{"monto": -40.0}'

# 4. Consultar historial acumulado de auditoría
curl -X GET http://localhost:8080/api/transacciones
```

---
