<script lang="ts">
	import type { useImportData } from '../../hooks/useImportData.svelte';
	import { Button } from '../shadcn/button';
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		DialogFooter,
		DialogDescription
	} from '../shadcn/dialog';
	import { Upload, FileText } from '@lucide/svelte';

	let { tableName, controller, onSuccess } = $props<{
		tableName: string;
		controller: ReturnType<typeof useImportData>;
		onSuccess?: () => void;
	}>();

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		controller.isDragging = true;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		controller.isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		controller.isDragging = false;
		if (e.dataTransfer && e.dataTransfer.files.length > 0) {
			controller.processFile(e.dataTransfer.files[0]);
		}
	}

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			controller.processFile(target.files[0]);
		}
	}
</script>

<Dialog bind:open={controller.isOpen}>
	<DialogContent class="max-w-2xl">
		<DialogHeader>
			<DialogTitle class="flex items-center gap-1.5 text-base font-bold">
				<Upload class="size-4 text-primary" /> Bulk Import Data to <span class="text-primary font-mono">{tableName}</span>
			</DialogTitle>
			<DialogDescription>
				Upload a .csv or .json file to insert or upsert records in bulk.
			</DialogDescription>
		</DialogHeader>

		{#if controller.errorMessage}
			<div class="rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-xs text-destructive">
				⚠️ {controller.errorMessage}
			</div>
		{/if}

		{#if controller.statusMessage}
			<div class="rounded-lg bg-success/10 border border-success/20 p-3 text-xs text-success font-medium">
				{controller.statusMessage}
			</div>
		{/if}

		<!-- Drag & Drop Zone -->
		<div
			ondragover={handleDragOver}
			ondragleave={handleDragLeave}
			ondrop={handleDrop}
			class="border-2 border-dashed rounded-xl p-6 text-center transition flex flex-col items-center justify-center cursor-pointer {
				controller.isDragging
					? 'border-primary bg-primary/10'
					: 'border-border bg-background/60 hover:border-primary/50 hover:bg-background'
			}"
			role="region"
			aria-label="Drag and Drop Upload Area"
		>
			<input type="file" accept=".csv,.json" onchange={handleFileSelect} class="hidden" id="fileInput" />
			<label for="fileInput" class="cursor-pointer w-full flex flex-col items-center">
				<FileText class="size-8 text-primary mb-2" />
				{#if controller.fileName}
					<span class="text-xs font-semibold text-primary font-mono">{controller.fileName}</span>
					<span class="text-[11px] text-muted-foreground mt-1">({controller.parsedRows.length} rows detected)</span>
				{:else}
					<span class="text-xs font-semibold text-foreground">Drag & Drop your .csv or .json file here</span>
					<span class="text-[11px] text-muted-foreground mt-1">or click to browse files from your computer</span>
				{/if}
			</label>
		</div>

		<!-- Live Preview & Settings -->
		{#if controller.parsedRows.length > 0}
			<div class="space-y-4 text-xs">
				<!-- Import Mode Toggle -->
				<div class="flex items-center justify-between bg-secondary/50 rounded-lg p-3 border border-border">
					<span class="font-semibold text-foreground">Import Mode Strategy:</span>
					<div class="flex gap-2">
						<Button
							type="button"
							variant={controller.importMode === 'insert' ? 'default' : 'outline'}
							size="xs"
							onclick={() => (controller.importMode = 'insert')}
						>
							Insert Only
						</Button>
						<Button
							type="button"
							variant={controller.importMode === 'upsert' ? 'default' : 'outline'}
							size="xs"
							onclick={() => (controller.importMode = 'upsert')}
						>
							Upsert (Update if PK Exists)
						</Button>
					</div>
				</div>

				<!-- Sample Table Preview -->
				<div>
					<div class="text-[11px] font-semibold text-muted-foreground mb-2">
						Sample Data Preview (First 5 of {controller.parsedRows.length} rows):
					</div>
					<div class="max-h-36 overflow-auto rounded-lg border border-border bg-background text-xs">
						<table class="w-full text-left font-mono">
							<thead class="bg-secondary text-muted-foreground border-b border-border">
								<tr>
									{#each Object.keys(controller.previewRows[0] || {}) as key}
										<th class="p-2 border-r border-border">{key}</th>
									{/each}
								</tr>
							</thead>
							<tbody>
								{#each controller.previewRows as row}
									<tr class="border-b border-border/60 hover:bg-secondary/40">
										{#each Object.keys(controller.previewRows[0] || {}) as key}
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

		<DialogFooter class="pt-2 border-t border-border">
			<Button type="button" variant="outline" size="sm" onclick={controller.closeModal}>
				Cancel
			</Button>
			<Button
				type="button"
				size="sm"
				onclick={() => controller.executeImport(onSuccess)}
				disabled={controller.parsedRows.length === 0 || controller.isUploading}
			>
				{controller.isUploading ? 'Processing Import...' : `Import ${controller.parsedRows.length} Rows`}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
