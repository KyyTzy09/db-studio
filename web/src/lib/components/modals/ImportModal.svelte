<script lang="ts">
	import { batchInsertOrUpdate } from '../../../data/services';
	import { Button } from '../shadcn/button';
	import { Upload, FileText, X } from '@lucide/svelte';

	let { tableName, isOpen = $bindable(false), onClose, onSuccess } = $props<{
		tableName: string;
		isOpen?: boolean;
		onClose?: () => void;
		onSuccess?: () => void;
	}>();

	let isDragging = $state(false);
	let fileName = $state('');
	let parsedRows = $state<Record<string, any>[]>([]);
	let previewRows = $state<Record<string, any>[]>([]);
	let importMode = $state<'insert' | 'upsert'>('insert');
	let isUploading = $state(false);
	let errorMessage = $state('');
	let statusMessage = $state('');

	$effect(() => {
		if (isOpen) {
			resetState();
		}
	});

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

	function handleClose() {
		isOpen = false;
		if (onClose) onClose();
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		if (e.dataTransfer && e.dataTransfer.files.length > 0) {
			processFile(e.dataTransfer.files[0]);
		}
	}

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			processFile(target.files[0]);
		}
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

	async function handleExecuteImport() {
		if (parsedRows.length === 0) return;
		isUploading = true;
		errorMessage = '';
		statusMessage = `Importing ${parsedRows.length} rows...`;

		try {
			const res = await batchInsertOrUpdate(tableName, parsedRows, importMode);
			statusMessage = `✅ Success! Processed ${res.affected_rows} rows.`;
			setTimeout(() => {
				if (onSuccess) onSuccess();
				handleClose();
			}, 800);
		} catch (err: any) {
			errorMessage = err.message || 'Failed to import data';
			statusMessage = '';
		} finally {
			isUploading = false;
		}
	}
</script>

{#if isOpen}
	<div
		onclick={(e) => { if (e.target === e.currentTarget) handleClose(); }}
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 cursor-pointer"
		role="dialog"
		aria-modal="true"
		tabindex="-1"
		onkeydown={(e) => { if (e.key === 'Escape') handleClose(); }}
	>
		<div
			class="w-full max-w-2xl rounded-xl border border-border bg-card p-6 shadow-2xl cursor-default"
			onclick={(e) => e.stopPropagation()}
			role="document"
			tabindex="-1"
		>
			<div class="flex items-center justify-between border-b border-border pb-4">
				<h3 class="text-base font-bold text-foreground flex items-center gap-1.5">
					<Upload class="size-4 text-primary" /> Bulk Import Data to <span class="text-primary font-mono">{tableName}</span>
				</h3>
				<Button type="button" variant="ghost" size="icon-xs" onclick={handleClose}>
					<X class="size-4" />
				</Button>
			</div>

			{#if errorMessage}
				<div class="mt-4 rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-xs text-destructive">
					⚠️ {errorMessage}
				</div>
			{/if}

			{#if statusMessage}
				<div class="mt-4 rounded-lg bg-success/10 border border-success/20 p-3 text-xs text-success font-medium">
					{statusMessage}
				</div>
			{/if}

			<!-- Drag & Drop Zone -->
			<div
				ondragover={handleDragOver}
				ondragleave={handleDragLeave}
				ondrop={handleDrop}
				class="mt-4 border-2 border-dashed rounded-xl p-6 text-center transition flex flex-col items-center justify-center cursor-pointer {
					isDragging
						? 'border-primary bg-primary/10'
						: 'border-border bg-background/60 hover:border-primary/50 hover:bg-background'
				}"
				role="region"
				aria-label="Drag and Drop Upload Area"
			>
				<input type="file" accept=".csv,.json" onchange={handleFileSelect} class="hidden" id="fileInput" />
				<label for="fileInput" class="cursor-pointer w-full flex flex-col items-center">
					<FileText class="size-8 text-primary mb-2" />
					{#if fileName}
						<span class="text-xs font-semibold text-primary font-mono">{fileName}</span>
						<span class="text-[11px] text-muted-foreground mt-1">({parsedRows.length} rows detected)</span>
					{:else}
						<span class="text-xs font-semibold text-foreground">Drag & Drop your .csv or .json file here</span>
						<span class="text-[11px] text-muted-foreground mt-1">or click to browse files from your computer</span>
					{/if}
				</label>
			</div>

			<!-- Live Preview & Settings -->
			{#if parsedRows.length > 0}
				<div class="mt-4 space-y-4 text-xs">
					<!-- Import Mode Toggle -->
					<div class="flex items-center justify-between bg-secondary/50 rounded-lg p-3 border border-border">
						<span class="font-semibold text-foreground">Import Mode Strategy:</span>
						<div class="flex gap-2">
							<Button
								type="button"
								variant={importMode === 'insert' ? 'default' : 'outline'}
								size="xs"
								onclick={() => (importMode = 'insert')}
							>
								Insert Only
							</Button>
							<Button
								type="button"
								variant={importMode === 'upsert' ? 'default' : 'outline'}
								size="xs"
								onclick={() => (importMode = 'upsert')}
							>
								Upsert (Update if PK Exists)
							</Button>
						</div>
					</div>

					<!-- Sample Table Preview -->
					<div>
						<div class="text-[11px] font-semibold text-muted-foreground mb-2">
							Sample Data Preview (First 5 of {parsedRows.length} rows):
						</div>
						<div class="max-h-36 overflow-auto rounded-lg border border-border bg-background text-xs">
							<table class="w-full text-left font-mono">
								<thead class="bg-secondary text-muted-foreground border-b border-border">
									<tr>
										{#each Object.keys(previewRows[0] || {}) as key}
											<th class="p-2 border-r border-border">{key}</th>
										{/each}
									</tr>
								</thead>
								<tbody>
									{#each previewRows as row}
										<tr class="border-b border-border/60 hover:bg-secondary/40">
											{#each Object.keys(previewRows[0] || {}) as key}
												<td class="p-2 border-r border-border text-foreground">{row[key]}</td>
											{/each}
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>
				</div>
			{/if}

			<div class="pt-6 flex justify-end gap-3 border-t border-border mt-6">
				<Button type="button" variant="outline" size="sm" onclick={handleClose}>
					Cancel
				</Button>
				<Button
					type="button"
					size="sm"
					onclick={handleExecuteImport}
					disabled={parsedRows.length === 0 || isUploading}
				>
					{isUploading ? 'Processing Import...' : `Import ${parsedRows.length} Rows`}
				</Button>
			</div>
		</div>
	</div>
{/if}
