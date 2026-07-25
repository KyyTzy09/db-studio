import { fetchTables } from '../../data/services/dbService';
import { tablesList } from '../stores/dbStore';

export interface ColumnSpec {
	name: string;
	data_type: string;
	is_primary_key: boolean;
	is_nullable: boolean;
	default_value: string;
	auto_increment: boolean;
	is_foreign_key?: boolean;
	fk_table?: string;
	fk_column?: string;
}

export function useSchemaEditor() {
	let isCreateTableOpen = $state(false);
	let isAddColumnOpen = $state(false);
	let targetTableForColumn = $state('');

	let tableName = $state('');
	let columns = $state<ColumnSpec[]>([
		{ name: 'id', data_type: 'INTEGER', is_primary_key: true, is_nullable: false, default_value: '', auto_increment: true },
		{ name: 'created_at', data_type: 'TIMESTAMP', is_primary_key: false, is_nullable: false, default_value: 'CURRENT_TIMESTAMP', auto_increment: false }
	]);

	let newColumn = $state<ColumnSpec>({
		name: '',
		data_type: 'VARCHAR(255)',
		is_primary_key: false,
		is_nullable: true,
		default_value: '',
		auto_increment: false
	});

	let isSubmitting = $state(false);

	function openCreateTableModal() {
		tableName = '';
		columns = [
			{ name: 'id', data_type: 'INTEGER', is_primary_key: true, is_nullable: false, default_value: '', auto_increment: true },
			{ name: 'created_at', data_type: 'TIMESTAMP', is_primary_key: false, is_nullable: false, default_value: 'CURRENT_TIMESTAMP', auto_increment: false }
		];
		isCreateTableOpen = true;
	}

	function closeCreateTableModal() {
		isCreateTableOpen = false;
	}

	function openAddColumnModal(table: string) {
		targetTableForColumn = table;
		newColumn = {
			name: '',
			data_type: 'VARCHAR(255)',
			is_primary_key: false,
			is_nullable: true,
			default_value: '',
			auto_increment: false
		};
		isAddColumnOpen = true;
	}

	function closeAddColumnModal() {
		isAddColumnOpen = false;
	}

	function addColumnRow() {
		columns.push({
			name: `col_${columns.length + 1}`,
			data_type: 'VARCHAR(255)',
			is_primary_key: false,
			is_nullable: true,
			default_value: '',
			auto_increment: false,
			is_foreign_key: false,
			fk_table: '',
			fk_column: 'id'
		});
	}

	function removeColumnRow(index: number) {
		if (columns.length <= 1) return;
		columns.splice(index, 1);
	}

	function toggleAutoIncrement(col: ColumnSpec, value: boolean) {
		col.auto_increment = value;
		if (value) {
			col.data_type = 'INTEGER';
		}
	}

	// Live SQL Preview Calculation
	let generatedSql = $derived.by(() => {
		if (!tableName.trim()) return '-- Enter table name to see preview';
		const colDefs: string[] = [];
		const pks: string[] = [];

		for (const c of columns) {
			if (!c.name.trim()) continue;
			let def = `${c.name} ${c.data_type}`;
			if (c.auto_increment) def += ' AUTO_INCREMENT';
			if (!c.is_nullable) def += ' NOT NULL';
			if (c.default_value) def += ` DEFAULT ${c.default_value}`;
			colDefs.push(def);
			if (c.is_primary_key) pks.push(c.name);
		}

		if (pks.length > 0) {
			colDefs.push(`PRIMARY KEY (${pks.join(', ')})`);
		}

		for (const c of columns) {
			if (c.is_foreign_key && c.fk_table && c.fk_column) {
				colDefs.push(`FOREIGN KEY (${c.name}) REFERENCES ${c.fk_table}(${c.fk_column})`);
			}
		}

		return `CREATE TABLE ${tableName.trim()} (\n  ${colDefs.join(',\n  ')}\n);`;
	});

	async function submitCreateTable() {
		if (!tableName.trim()) {
			alert('Table name is required');
			return false;
		}

		const validColumns = columns.filter((c) => c.name.trim().length > 0);
		if (validColumns.length === 0) {
			alert('At least one valid column is required');
			return false;
		}

		isSubmitting = true;
		try {
			const res = await fetch('/api/tables', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					table_name: tableName.trim(),
					columns: validColumns
				})
			});

			if (!res.ok) {
				const err = await res.json();
				throw new Error(err.error || 'Failed to create table');
			}

			const { tables } = await fetchTables();
			tablesList.set(tables || []);
			closeCreateTableModal();
			return true;
		} catch (err: any) {
			alert(err.message || 'Error creating table');
			return false;
		} finally {
			isSubmitting = false;
		}
	}

	async function submitAddColumn() {
		if (!newColumn.name.trim()) {
			alert('Column name is required');
			return false;
		}

		isSubmitting = true;
		try {
			const res = await fetch(`/api/tables/${encodeURIComponent(targetTableForColumn)}/columns`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(newColumn)
			});

			if (!res.ok) {
				const err = await res.json();
				throw new Error(err.error || 'Failed to add column');
			}

			closeAddColumnModal();
			window.location.reload();
			return true;
		} catch (err: any) {
			alert(err.message || 'Error adding column');
			return false;
		} finally {
			isSubmitting = false;
		}
	}

	async function submitDropColumn(table: string, columnName: string) {
		if (!confirm(`Are you sure you want to drop column '${columnName}' from table '${table}'? Data in this column will be lost!`)) {
			return false;
		}

		try {
			const res = await fetch(`/api/tables/${encodeURIComponent(table)}/columns/${encodeURIComponent(columnName)}`, {
				method: 'DELETE'
			});

			if (!res.ok) {
				const err = await res.json();
				throw new Error(err.error || 'Failed to drop column');
			}

			window.location.reload();
			return true;
		} catch (err: any) {
			alert(err.message || 'Error dropping column');
			return false;
		}
	}

	return {
		get isCreateTableOpen() { return isCreateTableOpen; },
		get isAddColumnOpen() { return isAddColumnOpen; },
		get tableName() { return tableName; },
		set tableName(val: string) { tableName = val; },
		get columns() { return columns; },
		get newColumn() { return newColumn; },
		get isSubmitting() { return isSubmitting; },
		get generatedSql() { return generatedSql; },
		get targetTableForColumn() { return targetTableForColumn; },

		openCreateTableModal,
		closeCreateTableModal,
		openAddColumnModal,
		closeAddColumnModal,
		addColumnRow,
		removeColumnRow,
		toggleAutoIncrement,
		submitCreateTable,
		submitAddColumn,
		submitDropColumn
	};
}
