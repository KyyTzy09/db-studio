<script lang="ts">
	import { onMount } from 'svelte';
	import { useERDiagram } from '../../hooks/useERDiagram.svelte';
	import TableNode from './TableNode.svelte';
	import RelationshipLine from './RelationshipLine.svelte';
	import { Button } from '../shadcn/button';
	import { ZoomIn, ZoomOut, RotateCcw, RefreshCw, Network, Info } from '@lucide/svelte';

	let { onSelectTable } = $props<{
		onSelectTable: (tableName: string) => void;
	}>();

	const controller = useERDiagram(onSelectTable);

	onMount(() => {
		controller.loadGraphData();
	});
</script>

<div class="flex-1 flex flex-col h-full bg-background overflow-hidden relative select-none">
	<!-- ER Diagram Toolbar Header -->
	<div class="h-12 px-6 border-b border-border bg-card/60 flex items-center justify-between z-20">
		<div class="flex items-center gap-2">
			<Network class="size-4 text-primary" />
			<h2 class="text-xs font-bold text-foreground uppercase tracking-wider">
				Visual ER Diagram Map
			</h2>
			{#if controller.graph}
				<span class="text-[10px] font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded-full border border-border">
					{controller.graph.nodes.length} Tables • {controller.graph.edges.length} Foreign Keys
				</span>
			{/if}
		</div>

		<!-- Zoom & Controls -->
		<div class="flex items-center gap-1.5">
			<Button variant="outline" size="icon-xs" onclick={controller.zoomOut} title="Zoom Out">
				<ZoomOut class="size-3.5" />
			</Button>
			<span class="text-[11px] font-mono text-muted-foreground w-12 text-center">
				{Math.round(controller.zoom * 100)}%
			</span>
			<Button variant="outline" size="icon-xs" onclick={controller.zoomIn} title="Zoom In">
				<ZoomIn class="size-3.5" />
			</Button>

			<div class="h-4 w-px bg-border mx-1"></div>

			<Button variant="outline" size="xs" onclick={controller.resetView} title="Reset Canvas View">
				<RotateCcw class="size-3 mr-1" /> Reset
			</Button>

			<Button variant="outline" size="xs" onclick={controller.loadGraphData} title="Refresh ER Graph">
				<RefreshCw class="size-3 mr-1" /> Refresh
			</Button>
		</div>
	</div>

	<!-- Interactive Diagram Canvas -->
	<div
		role="region"
		aria-label="ER Diagram Canvas Area"
		class="flex-1 overflow-hidden relative cursor-grab active:cursor-grabbing bg-background/50"
		onwheel={controller.handleWheel}
		onmousedown={controller.startCanvasPan}
		onmousemove={controller.handleMouseMove}
		onmouseup={controller.handleMouseUp}
		onmouseleave={controller.handleMouseUp}
	>
		{#if controller.loading}
			<div class="absolute inset-0 flex flex-col items-center justify-center gap-3 text-muted-foreground z-30">
				<div class="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
				<span class="text-xs font-medium">Extracting relational schema graph...</span>
			</div>
		{:else if controller.errorMsg}
			<div class="absolute inset-0 flex flex-col items-center justify-center p-6 z-30">
				<div class="max-w-sm bg-destructive/10 border border-destructive/20 rounded-xl p-4 text-center">
					<span class="text-xs font-semibold text-destructive">{controller.errorMsg}</span>
					<Button variant="destructive" size="xs" class="mt-3" onclick={controller.loadGraphData}>
						Retry Fetching Graph
					</Button>
				</div>
			</div>
		{:else if controller.graph && controller.graph.nodes.length === 0}
			<div class="absolute inset-0 flex flex-col items-center justify-center text-muted-foreground text-xs gap-2 z-30">
				<Info class="size-8 text-muted-foreground/60" />
				<span>Tidak ada tabel yang ditemukan pada database saat ini.</span>
			</div>
		{:else if controller.graph}
			<!-- Transformed Canvas Layer -->
			<div
				class="absolute inset-0 origin-top-left transition-transform duration-75"
				style="transform: translate({controller.panX}px, {controller.panY}px) scale({controller.zoom});"
			>
				<!-- SVG Edges Connecting FKs -->
				<svg class="absolute inset-0 w-[5000px] h-[5000px] pointer-events-none z-0">
					{#each controller.graph.edges as edge (edge.id)}
						<RelationshipLine
							{edge}
							sourcePos={controller.nodePositions[edge.source_table]}
							targetPos={controller.nodePositions[edge.target_table]}
						/>
					{/each}
				</svg>

				<!-- HTML Draggable Table Card Nodes -->
				{#each controller.graph.nodes as schema (schema.table_name)}
					<TableNode
						{schema}
						x={controller.nodePositions[schema.table_name]?.x || 0}
						y={controller.nodePositions[schema.table_name]?.y || 0}
						isSelected={controller.selectedTableNode === schema.table_name}
						onMouseDown={(e) => controller.startNodeDrag(schema.table_name, e)}
						onSelectTable={(name) => controller.handleJumpToTable(name)}
					/>
				{/each}
			</div>
		{/if}
	</div>
</div>
