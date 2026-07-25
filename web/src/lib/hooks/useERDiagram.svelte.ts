import { fetchSchemaGraph } from '../../data/services';
import type { SchemaGraph } from '$lib/api';

export interface NodePosition {
	x: number;
	y: number;
}

export function useERDiagram(onSelectTable?: (tableName: string) => void) {
	let graph = $state<SchemaGraph | null>(null);
	let loading = $state(true);
	let errorMsg = $state<string | null>(null);

	// Canvas Controls State
	let zoom = $state(1.0);
	let panX = $state(50);
	let panY = $state(50);

	// Dragging State
	let isPanning = $state(false);
	let startPanMouse = $state({ x: 0, y: 0 });

	let draggedNode = $state<string | null>(null);
	let dragOffset = $state({ x: 0, y: 0 });

	// Node Positions map { [tableName]: { x, y } }
	let nodePositions = $state<Record<string, NodePosition>>({});
	let selectedTableNode = $state<string | null>(null);

	const STORAGE_KEY = 'dbstudio_er_positions';

	function loadPositionsFromStorage(tables: string[]) {
		try {
			const saved = localStorage.getItem(STORAGE_KEY);
			const parsed = saved ? JSON.parse(saved) : {};
			const positions: Record<string, NodePosition> = {};

			// Auto grid layout defaults
			const cols = Math.ceil(Math.sqrt(tables.length)) || 3;
			const cardWidth = 300;
			const cardHeight = 240;

			tables.forEach((tableName, idx) => {
				if (parsed[tableName]) {
					positions[tableName] = parsed[tableName];
				} else {
					const col = idx % cols;
					const row = Math.floor(idx / cols);
					positions[tableName] = {
						x: col * cardWidth + 40,
						y: row * cardHeight + 40
					};
				}
			});

			nodePositions = positions;
		} catch (e) {
			console.error('Failed to parse node positions from storage:', e);
		}
	}

	function savePositionsToStorage() {
		try {
			localStorage.setItem(STORAGE_KEY, JSON.stringify(nodePositions));
		} catch (e) {
			console.error('Failed to save node positions:', e);
		}
	}

	async function loadGraphData() {
		loading = true;
		errorMsg = null;
		try {
			graph = await fetchSchemaGraph();
			const tableNames = (graph?.nodes || []).map((n) => n.table_name);
			loadPositionsFromStorage(tableNames);
		} catch (err: any) {
			errorMsg = err.message || 'Failed to load ER diagram schema';
		} finally {
			loading = false;
		}
	}

	function handleWheel(e: WheelEvent) {
		e.preventDefault();
		const zoomFactor = e.deltaY < 0 ? 1.1 : 0.9;
		const newZoom = Math.min(Math.max(zoom * zoomFactor, 0.2), 2.5);
		zoom = newZoom;
	}

	function startCanvasPan(e: MouseEvent) {
		if (e.button !== 0 || draggedNode) return;
		isPanning = true;
		startPanMouse = { x: e.clientX - panX, y: e.clientY - panY };
	}

	function startNodeDrag(tableName: string, e: MouseEvent) {
		e.stopPropagation();
		draggedNode = tableName;
		selectedTableNode = tableName;

		const pos = nodePositions[tableName] || { x: 0, y: 0 };
		dragOffset = {
			x: (e.clientX - panX) / zoom - pos.x,
			y: (e.clientY - panY) / zoom - pos.y
		};
	}

	function handleMouseMove(e: MouseEvent) {
		if (isPanning) {
			panX = e.clientX - startPanMouse.x;
			panY = e.clientY - startPanMouse.y;
		} else if (draggedNode) {
			const newX = (e.clientX - panX) / zoom - dragOffset.x;
			const newY = (e.clientY - panY) / zoom - dragOffset.y;

			nodePositions = {
				...nodePositions,
				[draggedNode]: { x: Math.round(newX), y: Math.round(newY) }
			};
		}
	}

	function handleMouseUp() {
		if (draggedNode) {
			savePositionsToStorage();
		}
		isPanning = false;
		draggedNode = null;
	}

	function zoomIn() {
		zoom = Math.min(zoom + 0.15, 2.5);
	}

	function zoomOut() {
		zoom = Math.max(zoom - 0.15, 0.2);
	}

	function resetView() {
		zoom = 1.0;
		panX = 50;
		panY = 50;
	}

	return {
		get graph() {
			return graph;
		},
		get loading() {
			return loading;
		},
		get errorMsg() {
			return errorMsg;
		},
		get zoom() {
			return zoom;
		},
		get panX() {
			return panX;
		},
		get panY() {
			return panY;
		},
		get nodePositions() {
			return nodePositions;
		},
		get selectedTableNode() {
			return selectedTableNode;
		},
		loadGraphData,
		handleWheel,
		startCanvasPan,
		startNodeDrag,
		handleMouseMove,
		handleMouseUp,
		zoomIn,
		zoomOut,
		resetView,
		handleJumpToTable(tableName: string) {
			if (onSelectTable) onSelectTable(tableName);
		}
	};
}
