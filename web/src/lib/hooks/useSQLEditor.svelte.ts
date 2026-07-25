import { executeRawQuery, fetchTables } from '../../data/services';
import type { QueryResult } from '$lib/api';
import { tablesList } from '../stores/dbStore';

export function useSQLEditor() {
	let queryInput = $state('SELECT * FROM information_schema.tables LIMIT 10;');
	let queryResult = $state<QueryResult | null>(null);
	let executing = $state(false);
	let errorMsg = $state<string | null>(null);

	let showDangerModal = $state(false);
	let dangerWarning = $state('');

	async function runQuery(force = false) {
		if (!queryInput.trim()) return;
		executing = true;
		errorMsg = null;

		try {
			const res = await executeRawQuery(queryInput, force);
			if (res.is428 && res.warning) {
				dangerWarning = res.warning;
				showDangerModal = true;
				return;
			}
			if (res.data) {
				queryResult = res.data;
				try {
					const { tables } = await fetchTables();
					tablesList.set(tables || []);
				} catch (e) {
					// Ignore background refresh errors
				}
			}
		} catch (err: any) {
			errorMsg = err.message || 'Failed to execute SQL query';
		} finally {
			executing = false;
		}
	}

	function confirmDangerousQuery() {
		showDangerModal = false;
		runQuery(true);
	}

	return {
		get queryInput() {
			return queryInput;
		},
		set queryInput(v: string) {
			queryInput = v;
		},
		get queryResult() {
			return queryResult;
		},
		get executing() {
			return executing;
		},
		get errorMsg() {
			return errorMsg;
		},
		get showDangerModal() {
			return showDangerModal;
		},
		set showDangerModal(v: boolean) {
			showDangerModal = v;
		},
		get dangerWarning() {
			return dangerWarning;
		},
		runQuery,
		confirmDangerousQuery
	};
}
