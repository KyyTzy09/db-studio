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
