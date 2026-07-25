<script lang="ts">
	import type { ForeignKeyRelation } from '$lib/api';

	let { edge, sourcePos, targetPos } = $props<{
		edge: ForeignKeyRelation;
		sourcePos: { x: number; y: number } | undefined;
		targetPos: { x: number; y: number } | undefined;
	}>();

	// Calculate Bezier curve path connecting node centers
	let pathData = $derived.by(() => {
		if (!sourcePos || !targetPos) return '';

		// Approximate node widths & heights
		const width = 256;
		const height = 120;

		const x1 = sourcePos.x + width / 2;
		const y1 = sourcePos.y + height / 2;
		const x2 = targetPos.x + width / 2;
		const y2 = targetPos.y + height / 2;

		const dx = Math.max(Math.abs(x2 - x1) / 2, 50);

		return `M ${x1} ${y1} C ${x1 + dx} ${y1}, ${x2 - dx} ${y2}, ${x2} ${y2}`;
	});

	let midPoint = $derived.by(() => {
		if (!sourcePos || !targetPos) return { x: 0, y: 0 };
		const width = 256;
		const height = 120;
		return {
			x: (sourcePos.x + targetPos.x + width) / 2,
			y: (sourcePos.y + targetPos.y + height) / 2
		};
	});
</script>

{#if sourcePos && targetPos && pathData}
	<g class="group">
		<!-- Background Glow Path -->
		<path
			d={pathData}
			fill="none"
			stroke="var(--color-primary)"
			stroke-width="4"
			stroke-opacity="0.15"
			class="transition-all group-hover:stroke-opacity-40"
		/>

		<!-- Main Relationship Path -->
		<path
			d={pathData}
			fill="none"
			stroke="var(--color-primary)"
			stroke-width="2"
			stroke-dasharray="6,4"
			class="transition-all group-hover:stroke-width-3"
		/>

		<!-- Source Marker Circle -->
		<circle
			cx={sourcePos.x + 128}
			cy={sourcePos.y + 60}
			r="4"
			fill="var(--color-primary)"
		/>

		<!-- Target Marker Circle -->
		<circle
			cx={targetPos.x + 128}
			cy={targetPos.y + 60}
			r="4"
			fill="var(--color-warning)"
		/>

		<!-- Column Label Badge -->
		<g transform="translate({midPoint.x}, {midPoint.y})">
			<rect
				x="-45"
				y="-10"
				width="90"
				height="20"
				rx="4"
				fill="var(--color-card)"
				stroke="var(--color-border)"
				stroke-width="1"
				class="shadow-xs"
			/>
			<text
				x="0"
				y="3"
				text-anchor="middle"
				fill="var(--color-foreground)"
				font-size="10"
				font-family="monospace"
				font-weight="500"
			>
				{edge.source_column}
			</text>
		</g>
	</g>
{/if}
