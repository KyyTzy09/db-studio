/**
 * Utility functions for exporting data to CSV and JSON formats
 */

export function exportToJSON(filename: string, data: Record<string, any>[]): void {
	const jsonStr = JSON.stringify(data, null, 2);
	const blob = new Blob([jsonStr], { type: 'application/json;charset=utf-8;' });
	triggerDownload(blob, `${filename}.json`);
}

export function exportToCSV(filename: string, columns: string[], rows: Record<string, any>[]): void {
	if (!columns || columns.length === 0) return;

	const escapeCSVCell = (val: any): string => {
		if (val === null || val === undefined) return '';
		const str = String(val);
		if (str.includes(',') || str.includes('"') || str.includes('\n')) {
			return `"${str.replace(/"/g, '""')}"`;
		}
		return str;
	};

	const headerRow = columns.map(escapeCSVCell).join(',');
	const bodyRows = rows.map((row) =>
		columns.map((col) => escapeCSVCell(row[col])).join(',')
	);

	const csvContent = [headerRow, ...bodyRows].join('\n');
	const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
	triggerDownload(blob, `${filename}.csv`);
}

function triggerDownload(blob: Blob, fileName: string): void {
	const url = URL.createObjectURL(blob);
	const a = document.createElement('a');
	a.href = url;
	a.download = fileName;
	document.body.appendChild(a);
	a.click();
	document.body.removeChild(a);
	URL.revokeObjectURL(url);
}
