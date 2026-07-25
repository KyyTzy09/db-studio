import { fetchTableData, fetchTableSchema, deleteTableRow } from '../../data/services';
import type { QueryResult, TableSchema } from '$lib/api';

export type FilterOperator = 'contains' | 'equals' | 'gt' | 'lt' | 'starts' | 'is_null';

export interface FilterRule {
	id: string;
	column: string;
	operator: FilterOperator;
	value: string;
}

export interface SortState {
	column: string | null;
	direction: 'asc' | 'desc' | null;
}

export function useTableData(getTableName: () => string) {
	let dataResult = $state<QueryResult | null>(null);
	let schema = $state<TableSchema | null>(null);
	let loading = $state(true);
	let errorMsg = $state<string | null>(null);
	let searchQuery = $state('');
	let currentPage = $state(1);
	let pageSize = $state(15);

	// Filtering State
	let filterRules = $state<FilterRule[]>([]);
	let showFilterPanel = $state(false);

	// Sorting State
	let sortState = $state<SortState>({ column: null, direction: null });

	let showDeleteModal = $state(false);
	let selectedRowToDelete = $state<Record<string, any> | null>(null);
	let deleting = $state(false);

	async function loadData() {
		const name = getTableName();
		if (!name) return;
		loading = true;
		errorMsg = null;
		currentPage = 1;
		try {
			const [dataRes, schemaRes] = await Promise.all([
				fetchTableData(name).catch(() => null),
				fetchTableSchema(name).catch(() => null)
			]);
			dataResult = dataRes;
			schema = schemaRes;
		} catch (err: any) {
			errorMsg = err.message || 'Failed to load table data';
		} finally {
			loading = false;
		}
	}

	function addFilterRule(rule: Omit<FilterRule, 'id'>) {
		filterRules = [
			...filterRules,
			{ ...rule, id: Math.random().toString(36).substring(2, 9) }
		];
		currentPage = 1;
	}

	function removeFilterRule(id: string) {
		filterRules = filterRules.filter((r) => r.id !== id);
		currentPage = 1;
	}

	function clearAllFilters() {
		filterRules = [];
		searchQuery = '';
		currentPage = 1;
	}

	function toggleSort(colName: string) {
		if (sortState.column !== colName) {
			sortState = { column: colName, direction: 'asc' };
		} else if (sortState.direction === 'asc') {
			sortState = { column: colName, direction: 'desc' };
		} else {
			sortState = { column: null, direction: null };
		}
	}

	function openDeleteModal(row: Record<string, any>) {
		selectedRowToDelete = row;
		showDeleteModal = true;
	}

	function closeDeleteModal() {
		showDeleteModal = false;
		selectedRowToDelete = null;
	}

	async function confirmDelete() {
		if (!selectedRowToDelete || !dataResult) return;
		deleting = true;
		try {
			await deleteTableRow(getTableName(), selectedRowToDelete);
			closeDeleteModal();
			await loadData();
		} catch (err: any) {
			alert(err.message || 'Failed to delete row');
		} finally {
			deleting = false;
		}
	}

	// Derived Rows Filtering & Sorting Pipeline
	let rowsList = $derived(dataResult?.rows || []);

	let filteredRows = $derived.by(() => {
		let list = [...rowsList];

		// 1. Global Search Query
		if (searchQuery.trim()) {
			const q = searchQuery.toLowerCase();
			list = list.filter((row) =>
				Object.values(row || {}).some((val) =>
					String(val ?? '').toLowerCase().includes(q)
				)
			);
		}

		// 2. Multi-Condition Filter Rules (AND logic)
		if (filterRules.length > 0) {
			list = list.filter((row) => {
				return filterRules.every((rule) => {
					const rawVal = row[rule.column];
					if (rule.operator === 'is_null') {
						return rawVal === null || rawVal === undefined || rawVal === '';
					}
					if (rawVal === null || rawVal === undefined) return false;

					const strVal = String(rawVal).toLowerCase();
					const targetVal = rule.value.toLowerCase();

					switch (rule.operator) {
						case 'contains':
							return strVal.includes(targetVal);
						case 'equals':
							return strVal === targetVal;
						case 'starts':
							return strVal.startsWith(targetVal);
						case 'gt': {
							const numRaw = Number(rawVal);
							const numTarget = Number(rule.value);
							if (!isNaN(numRaw) && !isNaN(numTarget)) {
								return numRaw > numTarget;
							}
							return strVal > targetVal;
						}
						case 'lt': {
							const numRaw = Number(rawVal);
							const numTarget = Number(rule.value);
							if (!isNaN(numRaw) && !isNaN(numTarget)) {
								return numRaw < numTarget;
							}
							return strVal < targetVal;
						}
						default:
							return true;
					}
				});
			});
		}

		// 3. Sorting
		if (sortState.column && sortState.direction) {
			const col = sortState.column;
			const dir = sortState.direction === 'asc' ? 1 : -1;
			list.sort((a, b) => {
				const valA = a[col];
				const valB = b[col];
				if (valA === null || valA === undefined) return 1;
				if (valB === null || valB === undefined) return -1;

				const numA = Number(valA);
				const numB = Number(valB);
				if (!isNaN(numA) && !isNaN(numB)) {
					return (numA - numB) * dir;
				}
				return String(valA).localeCompare(String(valB)) * dir;
			});
		}

		return list;
	});

	let totalPages = $derived(Math.ceil(filteredRows.length / pageSize) || 1);

	let paginatedRows = $derived(
		filteredRows.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	return {
		get dataResult() {
			return dataResult;
		},
		get schema() {
			return schema;
		},
		get loading() {
			return loading;
		},
		get errorMsg() {
			return errorMsg;
		},
		get searchQuery() {
			return searchQuery;
		},
		set searchQuery(v: string) {
			searchQuery = v;
		},
		get currentPage() {
			return currentPage;
		},
		set currentPage(v: number) {
			currentPage = v;
		},
		get pageSize() {
			return pageSize;
		},
		get filterRules() {
			return filterRules;
		},
		get showFilterPanel() {
			return showFilterPanel;
		},
		set showFilterPanel(v: boolean) {
			showFilterPanel = v;
		},
		get sortState() {
			return sortState;
		},
		get filteredRows() {
			return filteredRows;
		},
		get totalPages() {
			return totalPages;
		},
		get paginatedRows() {
			return paginatedRows;
		},
		get showDeleteModal() {
			return showDeleteModal;
		},
		set showDeleteModal(v: boolean) {
			showDeleteModal = v;
		},
		get selectedRowToDelete() {
			return selectedRowToDelete;
		},
		get deleting() {
			return deleting;
		},
		loadData,
		addFilterRule,
		removeFilterRule,
		clearAllFilters,
		toggleSort,
		openDeleteModal,
		closeDeleteModal,
		confirmDelete
	};
}
