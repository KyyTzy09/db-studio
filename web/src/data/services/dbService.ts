import type { ConnectionStatus, TableInfo, TableSchema, QueryResult } from '../models';

const API_BASE = '/api';

export async function fetchConnectionStatus(): Promise<ConnectionStatus> {
	const res = await fetch(`${API_BASE}/connection/status`);
	return res.json();
}

export async function fetchTables(): Promise<{ tables: TableInfo[] }> {
	const res = await fetch(`${API_BASE}/tables`);
	if (!res.ok) throw new Error('Failed to fetch tables');
	return res.json();
}

export async function fetchTableSchema(tableName: string): Promise<TableSchema> {
	const res = await fetch(`${API_BASE}/tables/${encodeURIComponent(tableName)}/schema`);
	if (!res.ok) throw new Error(`Failed to fetch schema for ${tableName}`);
	return res.json();
}

export async function fetchTableData(tableName: string): Promise<QueryResult> {
	const res = await fetch(`${API_BASE}/tables/${encodeURIComponent(tableName)}/data`);
	if (!res.ok) throw new Error(`Failed to fetch data for ${tableName}`);
	return res.json();
}

export async function executeRawQuery(query: string, force = false): Promise<{ data?: QueryResult; warning?: string; is428?: boolean }> {
	const res = await fetch(`${API_BASE}/query`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ query, force })
	});

	if (res.status === 428) {
		const json = await res.json();
		return { warning: json.warning || 'Destructive query warning', is428: true };
	}

	if (!res.ok) {
		const json = await res.json();
		throw new Error(json.error || 'Failed to execute query');
	}

	const data = await res.json();
	return { data };
}

export async function insertTableRow(tableName: string, data: Record<string, any>): Promise<void> {
	const res = await fetch(`${API_BASE}/tables/${encodeURIComponent(tableName)}`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(data)
	});
	if (!res.ok) {
		const json = await res.json();
		throw new Error(json.error || 'Failed to insert row');
	}
}

export async function updateTableRow(tableName: string, pk: Record<string, any>, data: Record<string, any>): Promise<void> {
	const res = await fetch(`${API_BASE}/tables/${encodeURIComponent(tableName)}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ pk, data })
	});
	if (!res.ok) {
		const json = await res.json();
		throw new Error(json.error || 'Failed to update row');
	}
}

export async function deleteTableRow(tableName: string, pk: Record<string, any>): Promise<void> {
	const res = await fetch(`${API_BASE}/tables/${encodeURIComponent(tableName)}`, {
		method: 'DELETE',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(pk)
	});
	if (!res.ok) {
		const json = await res.json();
		throw new Error(json.error || 'Failed to delete row');
	}
}

export async function batchInsertOrUpdate(tableName: string, rows: Record<string, any>[], mode: 'insert' | 'upsert'): Promise<{ affected_rows: number }> {
	const res = await fetch(`${API_BASE}/tables/${encodeURIComponent(tableName)}/batch`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ rows, mode })
	});
	if (!res.ok) {
		const json = await res.json();
		throw new Error(json.error || 'Failed to batch process rows');
	}
	return res.json();
}
