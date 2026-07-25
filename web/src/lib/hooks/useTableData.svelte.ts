import { fetchTableData, fetchTableSchema, deleteTableRow } from '../../data/services';
import type { QueryResult, TableSchema } from '$lib/api';

export function useTableData(getTableName: () => string) {
	let dataResult = $state<QueryResult | null>(null);
	let schema = $state<TableSchema | null>(null);
	let loading = $state(true);
	let errorMsg = $state<string | null>(null);
	let searchQuery = $state('');
	let currentPage = $state(1);
	let pageSize = $state(15);

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
		openDeleteModal,
		closeDeleteModal,
		confirmDelete
	};
}
