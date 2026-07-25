import { batchInsertOrUpdate } from '../../data/services';

export function useImportData(getTableName: () => string) {
	let isOpen = $state(false);
	let isDragging = $state(false);
	let fileName = $state('');
	let parsedRows = $state<Record<string, any>[]>([]);
	let previewRows = $state<Record<string, any>[]>([]);
	let importMode = $state<'insert' | 'upsert'>('insert');
	let isUploading = $state(false);
	let errorMessage = $state('');
	let statusMessage = $state('');

	function openModal() {
		resetState();
		isOpen = true;
	}

	function closeModal() {
		isOpen = false;
	}

	function resetState() {
		isDragging = false;
		fileName = '';
		parsedRows = [];
		previewRows = [];
		importMode = 'insert';
		isUploading = false;
		errorMessage = '';
		statusMessage = '';
	}

	function processFile(file: File) {
		errorMessage = '';
		fileName = file.name;
		const reader = new FileReader();

		if (file.name.endsWith('.json')) {
			reader.onload = (evt) => {
				try {
					const json = JSON.parse(evt.target?.result as string);
					if (Array.isArray(json)) {
						parsedRows = json;
						previewRows = json.slice(0, 5);
					} else {
						errorMessage = 'JSON file must contain an array of objects';
					}
				} catch (err: any) {
					errorMessage = 'Failed to parse JSON file: ' + err.message;
				}
			};
			reader.readAsText(file);
		} else if (file.name.endsWith('.csv')) {
			reader.onload = (evt) => {
				try {
					const text = evt.target?.result as string;
					const rows = parseCSVText(text);
					parsedRows = rows;
					previewRows = rows.slice(0, 5);
				} catch (err: any) {
					errorMessage = 'Failed to parse CSV file: ' + err.message;
				}
			};
			reader.readAsText(file);
		} else {
			errorMessage = 'Only .csv and .json files are supported!';
		}
	}

	function parseCSVText(text: string): Record<string, any>[] {
		const lines = text.split(/\r?\n/).filter((l) => l.trim() !== '');
		if (lines.length < 2) return [];

		const headers = parseCSVLine(lines[0]);
		const results: Record<string, any>[] = [];

		for (let i = 1; i < lines.length; i++) {
			const values = parseCSVLine(lines[i]);
			if (values.length === headers.length) {
				const row: Record<string, any> = {};
				headers.forEach((h, idx) => {
					row[h.trim()] = values[idx].trim();
				});
				results.push(row);
			}
		}
		return results;
	}

	function parseCSVLine(line: string): string[] {
		const result: string[] = [];
		let current = '';
		let inQuotes = false;

		for (let i = 0; i < line.length; i++) {
			const char = line[i];
			if (char === '"') {
				inQuotes = !inQuotes;
			} else if (char === ',' && !inQuotes) {
				result.push(current);
				current = '';
			} else {
				current += char;
			}
		}
		result.push(current);
		return result;
	}

	async function executeImport(onSuccess?: () => void) {
		if (parsedRows.length === 0) return;
		isUploading = true;
		errorMessage = '';
		statusMessage = `Importing ${parsedRows.length} rows...`;

		try {
			const res = await batchInsertOrUpdate(getTableName(), parsedRows, importMode);
			statusMessage = `✅ Success! Processed ${res.affected_rows} rows.`;
			setTimeout(() => {
				if (onSuccess) onSuccess();
				closeModal();
			}, 800);
		} catch (err: any) {
			errorMessage = err.message || 'Failed to import data';
			statusMessage = '';
		} finally {
			isUploading = false;
		}
	}

	return {
		get isOpen() {
			return isOpen;
		},
		set isOpen(v: boolean) {
			isOpen = v;
		},
		get isDragging() {
			return isDragging;
		},
		set isDragging(v: boolean) {
			isDragging = v;
		},
		get fileName() {
			return fileName;
		},
		get parsedRows() {
			return parsedRows;
		},
		get previewRows() {
			return previewRows;
		},
		get importMode() {
			return importMode;
		},
		set importMode(v: 'insert' | 'upsert') {
			importMode = v;
		},
		get isUploading() {
			return isUploading;
		},
		get errorMessage() {
			return errorMessage;
		},
		get statusMessage() {
			return statusMessage;
		},
		openModal,
		closeModal,
		processFile,
		executeImport
	};
}
