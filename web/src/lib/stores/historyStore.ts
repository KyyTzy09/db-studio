import { writable } from 'svelte/store';
import type { QueryHistoryItem, QuerySnippet } from '$lib/api';
import {
	fetchQueryHistory as apiFetchHistory,
	clearQueryHistory as apiClearHistory,
	fetchSnippets as apiFetchSnippets,
	saveSnippet as apiSaveSnippet,
	deleteSnippet as apiDeleteSnippet
} from '../../data/services/dbService';

export const queryHistory = writable<QueryHistoryItem[]>([]);
export const querySnippets = writable<QuerySnippet[]>([]);
export const isDrawerOpen = writable<boolean>(false);
export const activeDrawerTab = writable<'history' | 'snippets'>('history');

export async function refreshHistory() {
	try {
		const res = await apiFetchHistory();
		queryHistory.set(res.history || []);
	} catch (err) {
		console.error('Failed to load query history:', err);
	}
}

export async function clearHistory() {
	try {
		await apiClearHistory();
		queryHistory.set([]);
	} catch (err) {
		console.error('Failed to clear history:', err);
	}
}

export async function refreshSnippets() {
	try {
		const res = await apiFetchSnippets();
		querySnippets.set(res.snippets || []);
	} catch (err) {
		console.error('Failed to load query snippets:', err);
	}
}

export async function addOrUpdateSnippet(snippet: { id?: string; title: string; description?: string; query: string }) {
	await apiSaveSnippet(snippet);
	await refreshSnippets();
}

export async function removeSnippet(id: string) {
	await apiDeleteSnippet(id);
	await refreshSnippets();
}
