export interface TableInfo {
	name: string;
	type: string;
	row_count?: number;
}

export interface ColumnInfo {
	name: string;
	data_type: string;
	is_nullable: boolean;
	is_primary_key: boolean;
	is_foreign_key?: boolean;
	is_auto_increment?: boolean;
	default_value?: string;
}

export interface TableSchema {
	table_name: string;
	columns: ColumnInfo[];
}

export interface QueryResult {
	columns: string[];
	rows: Record<string, any>[];
	affected_rows: number;
	execution_ms: number;
}

export interface ForeignKeyRelation {
	id: string;
	source_table: string;
	source_column: string;
	target_table: string;
	target_column: string;
}

export interface SchemaGraph {
	nodes: TableSchema[];
	edges: ForeignKeyRelation[];
}
