import { writable, get } from 'svelte/store';
import type { ConnectionStatus, TableInfo, ColumnInfo } from '$lib/api';
import { fetchTables, fetchTableSchema } from '../../data/services/dbService';

export type ActiveTab = 'data' | 'schema' | 'sql' | 'diagram';

export const connectionStatus = writable<ConnectionStatus | null>(null);
export const tablesList = writable<TableInfo[]>([]);
export const selectedTable = writable<string | null>(null);
export const activeTab = writable<ActiveTab>('data');
export const isLoading = writable<boolean>(false);
export const errorMessage = writable<string | null>(null);

// Global schema cache store: Record<tableName, ColumnInfo[]>
export const tableSchemas = writable<Record<string, ColumnInfo[]>>({});

export async function refreshGlobalTables() {
	try {
		const data = await fetchTables();
		if (data.tables) {
			tablesList.set(data.tables);
			if (data.tables.length > 0 && !get(selectedTable)) {
				selectedTable.set(data.tables[0].name);
			}
		}
	} catch (err) {
		console.error('Failed to refresh global tables:', err);
	}
}

export async function fetchColumnsForTable(tableName: string): Promise<ColumnInfo[]> {
	if (!tableName) return [];
	const current = get(tableSchemas);
	if (current[tableName]) {
		return current[tableName];
	}
	try {
		const schema = await fetchTableSchema(tableName);
		if (schema && schema.columns) {
			tableSchemas.update((s) => ({ ...s, [tableName]: schema.columns }));
			return schema.columns;
		}
	} catch (err) {
		console.error(`Failed to fetch schema for ${tableName}:`, err);
	}
	return [];
}
