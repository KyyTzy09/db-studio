import { updateTableRow } from '../../data/services';
import type { ColumnInfo } from '../../data/models';

export function useEditRow(getTableName: () => string, getColumns: () => ColumnInfo[]) {
	let isOpen = $state(false);
	let formData = $state<Record<string, any>>({});
	let pkData = $state<Record<string, any>>({});
	let isSubmitting = $state(false);
	let errorMessage = $state('');

	function openModal(row: Record<string, any>) {
		errorMessage = '';
		formData = { ...row };
		const pk: Record<string, any> = {};

		const cols = getColumns();
		cols.forEach((col) => {
			if (col.is_primary_key && row[col.name] !== undefined) {
				pk[col.name] = row[col.name];
			}
		});

		if (Object.keys(pk).length === 0) {
			pkData = { ...row };
		} else {
			pkData = pk;
		}
		isOpen = true;
	}

	function closeModal() {
		isOpen = false;
	}

	async function submitForm(onSuccess?: () => void) {
		isSubmitting = true;
		errorMessage = '';
		try {
			const cols = getColumns();
			const cleanData: Record<string, any> = {};
			for (const col of cols) {
				if (col.is_primary_key) continue;
				const val = formData[col.name];
				if (val !== '' && val !== null && val !== undefined) {
					cleanData[col.name] = val;
				}
			}

			await updateTableRow(getTableName(), pkData, cleanData);
			if (onSuccess) onSuccess();
			closeModal();
		} catch (err: any) {
			errorMessage = err.message || 'Failed to update row';
		} finally {
			isSubmitting = false;
		}
	}

	return {
		get isOpen() {
			return isOpen;
		},
		set isOpen(v: boolean) {
			isOpen = v;
		},
		get formData() {
			return formData;
		},
		set formData(v: Record<string, any>) {
			formData = v;
		},
		get isSubmitting() {
			return isSubmitting;
		},
		get errorMessage() {
			return errorMessage;
		},
		set errorMessage(v: string) {
			errorMessage = v;
		},
		openModal,
		closeModal,
		submitForm
	};
}
