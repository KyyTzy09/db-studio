export interface ConnectionConfig {
	id?: string;
	name: string;
	driver: string;
	host?: string;
	port?: number;
	user?: string;
	password?: string;
	database?: string;
	file_path?: string;
	project_path?: string;
}

export interface ConnectionStatus {
	connected: boolean;
	error?: string;
	config?: ConnectionConfig;
}
