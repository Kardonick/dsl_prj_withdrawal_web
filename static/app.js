/*
Source file: static/app.js
Módulo: Lógica Cliente Asíncrona (Fetch API con Backend Go)
Curso: Metodologías de Desarrollo de Software / Arquitectura de Software
Fecha: nov-2026
Autor: J.Medianero.A / Refactorización Web (Acápite 4)
*/

document.addEventListener('DOMContentLoaded', () => {
    cargarEstadoCuenta();
    cargarHistorial();
});

// Asigna montos inmediatos al formulario
function setQuickAmount(valor) {
    document.getElementById('amountInput').value = valor;
}

// Consulta GET /api/cuenta
async function cargarEstadoCuenta() {
    try {
        const res = await fetch('/api/cuenta');
        if (!res.ok) throw new Error('No se pudo obtener el estado de la cuenta');
        const data = await res.json();

        document.getElementById('holderName').textContent = data.titular;
        document.getElementById('currentBalance').textContent = Number(data.saldo).toFixed(2);
    } catch (err) {
        console.error('Error cargando cuenta:', err);
        mostrarAlerta('danger', 'Error de Conexión', 'No se pudo comunicar con el servidor bancario.');
    }
}

// Procesa POST /api/retirar o POST /api/abonar
async function ejecutarOperacion(tipo) {
    const amountInput = document.getElementById('amountInput');
    const monto = parseFloat(amountInput.value);

    if (isNaN(monto)) {
        mostrarAlerta('warning', 'Monto Requerido', 'Por favor ingresa un monto numérico válido.');
        return;
    }

    const endpoint = tipo === 'abonar' ? '/api/abonar' : '/api/retirar';
    const btn = tipo === 'abonar' ? document.getElementById('btnDeposit') : document.getElementById('btnWithdraw');

    btn.disabled = true;

    try {
        const res = await fetch(endpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ monto: monto })
        });

        const txn = await res.json();

        // Actualiza saldo y tabla de transacciones de inmediato
        await cargarEstadoCuenta();
        await cargarHistorial();

        // Notificación basada en el estado de negocio (US-04)
        if (txn.estado === 'EXITOSA') {
            mostrarAlerta(
                'success',
                `Operación Exitosa (${txn.tipo})`,
                `${txn.mensaje}. Monto procesado: €${Number(txn.monto).toFixed(2)}. Saldo restante: €${Number(txn.saldo_posterior).toFixed(2)}.`
            );
        } else if (txn.estado === 'FALLO') {
            mostrarAlerta(
                'warning',
                `Operación Rechazada (${txn.tipo})`,
                `${txn.mensaje}. El saldo disponible permanece en €${Number(txn.saldo_posterior).toFixed(2)}.`
            );
        } else {
            mostrarAlerta(
                'danger',
                `Error de Validación (${txn.tipo})`,
                `${txn.mensaje}. La transacción fue cancelada y no afectó el saldo.`
            );
        }

        amountInput.value = '';
    } catch (err) {
        console.error(`Error en operación ${tipo}:`, err);
        mostrarAlerta('danger', 'Error Interno', 'Ocurrió un error al despachar la solicitud HTTP.');
    } finally {
        btn.disabled = false;
    }
}

// Consulta GET /api/transacciones y renderiza la tabla
async function cargarHistorial() {
    const tbody = document.getElementById('txTableBody');

    try {
        const res = await fetch('/api/transacciones');
        if (!res.ok) throw new Error('No se pudo obtener el historial');
        const transacciones = await res.json();

        if (!transacciones || transacciones.length === 0) {
            tbody.innerHTML = `
        <tr>
          <td colspan="7" class="text-center py-4 text-muted">
            No existen transacciones registradas aún en el archivo de auditoría.
          </td>
        </tr>`;
            return;
        }

        // Renderiza en orden inverso (las más recientes arriba)
        tbody.innerHTML = transacciones.slice().reverse().map(t => {
            let badgeClass = 'badge-neutral';
            if (t.estado === 'EXITOSA') badgeClass = 'badge-exitosa';
            else if (t.estado === 'FALLO') badgeClass = 'badge-fallo';
            else if (t.estado === 'ERROR') badgeClass = 'badge-error';

            return `
        <tr>
          <td><strong>${t.id}</strong></td>
          <td>${t.fecha_hora}</td>
          <td><strong>${t.tipo}</strong></td>
          <td>€${Number(t.monto).toFixed(2)}</td>
          <td><span class="badge ${badgeClass}">${t.estado}</span></td>
          <td>€${Number(t.saldo_posterior).toFixed(2)}</td>
          <td>${t.mensaje}</td>
        </tr>
      `;
        }).join('');

    } catch (err) {
        console.error('Error cargando historial:', err);
        tbody.innerHTML = `
      <tr>
        <td colspan="7" class="text-center py-4 text-muted">
          Error al cargar los registros desde transacciones.json.
        </td>
      </tr>`;
    }
}

// Sistema de alertas visuales
function mostrarAlerta(tipo, titulo, mensaje) {
    const alertEl = document.getElementById('alertFeedback');
    const titleEl = document.getElementById('alertTitle');
    const msgEl = document.getElementById('alertMessage');
    const iconEl = document.getElementById('alertIcon');

    alertEl.className = `alert alert-${tipo}`;
    titleEl.textContent = titulo;
    msgEl.textContent = mensaje;

    const iconSVGs = {
        success: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>`,
        warning: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`,
        danger: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>`
    };

    iconEl.innerHTML = iconSVGs[tipo] || iconSVGs.danger;

    alertEl.classList.remove('d-none');
}