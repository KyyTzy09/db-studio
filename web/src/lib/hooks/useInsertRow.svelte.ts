import { insertTableRow } from '../../data/services';
import type { ColumnInfo } from '../../data/models';

export function useInsertRow(getTableName: () => string, getColumns: () => ColumnInfo[]) {
	let isOpen = $state(false);
	let formData = $state<Record<string, any>>({});
	let isSubmitting = $state(false);
	let errorMessage = $state('');

	function openModal() {
		const cols = getColumns();
		const initial: Record<string, any> = {};
		cols.forEach((col) => {
			initial[col.name] = '';
		});
		formData = initial;
		errorMessage = '';
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
			const payload: Record<string, any> = {};
			for (const col of cols) {
				const val = formData[col.name];
				if (val !== '' && val !== null && val !== undefined) {
					payload[col.name] = val;
				}
			}

			await insertTableRow(getTableName(), payload);
			if (onSuccess) onSuccess();
			closeModal();
		} catch (err: any) {
			errorMessage = err.message || 'Failed to insert row';
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
