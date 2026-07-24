import { writable } from 'svelte/store';
import type { ConnectionStatus, TableInfo } from '$lib/api';

export type ActiveTab = 'data' | 'schema' | 'sql';

export const connectionStatus = writable<ConnectionStatus | null>(null);
export const tablesList = writable<TableInfo[]>([]);
export const selectedTable = writable<string | null>(null);
export const activeTab = writable<ActiveTab>('data');
export const isLoading = writable<boolean>(false);
export const errorMessage = writable<string | null>(null);
