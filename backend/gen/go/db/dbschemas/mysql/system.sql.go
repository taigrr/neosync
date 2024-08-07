// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: system.sql

package mysql_queries

import (
	"context"
	"database/sql"
)

const getDatabaseSchema = `-- name: GetDatabaseSchema :many
SELECT
	c.table_schema,
	c.table_name,
	c.column_name,
	c.ordinal_position,
	COALESCE(c.column_default, 'NULL') as column_default, -- must coalesce because sqlc doesn't appear to work for system structs to output a *string
	c.is_nullable,
	c.data_type,
	c.character_maximum_length,
  c.numeric_precision,
  c.numeric_scale,
	c.extra
FROM
	information_schema.columns AS c
	JOIN information_schema.tables AS t ON c.table_schema = t.table_schema
		AND c.table_name = t.table_name
WHERE
	c.table_schema NOT IN('sys', 'performance_schema', 'mysql')
	AND t.table_type = 'BASE TABLE'
`

type GetDatabaseSchemaRow struct {
	TableSchema            string
	TableName              string
	ColumnName             string
	OrdinalPosition        int64
	ColumnDefault          string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	Extra                  sql.NullString
}

func (q *Queries) GetDatabaseSchema(ctx context.Context, db DBTX) ([]*GetDatabaseSchemaRow, error) {
	rows, err := db.QueryContext(ctx, getDatabaseSchema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetDatabaseSchemaRow
	for rows.Next() {
		var i GetDatabaseSchemaRow
		if err := rows.Scan(
			&i.TableSchema,
			&i.TableName,
			&i.ColumnName,
			&i.OrdinalPosition,
			&i.ColumnDefault,
			&i.IsNullable,
			&i.DataType,
			&i.CharacterMaximumLength,
			&i.NumericPrecision,
			&i.NumericScale,
			&i.Extra,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getForeignKeyConstraints = `-- name: GetForeignKeyConstraints :many
SELECT
rc.constraint_name
,
kcu.table_schema AS schema_name
,
kcu.table_name as table_name
,
kcu.column_name as column_name
,
c.is_nullable as is_nullable
,
kcu.referenced_table_schema AS foreign_schema_name
,
kcu.referenced_table_name AS foreign_table_name
,
kcu.referenced_column_name AS foreign_column_name
FROM
	information_schema.referential_constraints rc
JOIN information_schema.key_column_usage kcu
	ON
	kcu.constraint_name = rc.constraint_name
	AND kcu.constraint_schema = rc.constraint_schema
JOIN information_schema.columns as c
	ON
	c.table_schema = kcu.table_schema
	AND c.table_name = kcu.table_name
	AND c.column_name = kcu.column_name
WHERE
	kcu.table_schema = ?
ORDER BY
	rc.constraint_name,
	kcu.ordinal_position
`

type GetForeignKeyConstraintsRow struct {
	ConstraintName    string
	SchemaName        string
	TableName         string
	ColumnName        string
	IsNullable        string
	ForeignSchemaName string
	ForeignTableName  string
	ForeignColumnName string
}

func (q *Queries) GetForeignKeyConstraints(ctx context.Context, db DBTX, tableSchema string) ([]*GetForeignKeyConstraintsRow, error) {
	rows, err := db.QueryContext(ctx, getForeignKeyConstraints, tableSchema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetForeignKeyConstraintsRow
	for rows.Next() {
		var i GetForeignKeyConstraintsRow
		if err := rows.Scan(
			&i.ConstraintName,
			&i.SchemaName,
			&i.TableName,
			&i.ColumnName,
			&i.IsNullable,
			&i.ForeignSchemaName,
			&i.ForeignTableName,
			&i.ForeignColumnName,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMysqlRolePermissions = `-- name: GetMysqlRolePermissions :many
SELECT
    tp.TABLE_SCHEMA AS table_schema,
    tp.TABLE_NAME AS table_name,
    tp.PRIVILEGE_TYPE AS privilege_type
FROM
    information_schema.TABLE_PRIVILEGES AS tp
WHERE
    tp.TABLE_SCHEMA NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
    AND (tp.GRANTEE = CONCAT("'", SUBSTRING_INDEX(CURRENT_USER(), '@', 1), "'@'%'"))
ORDER BY
    tp.TABLE_SCHEMA,
    tp.TABLE_NAME
`

type GetMysqlRolePermissionsRow struct {
	TableSchema   string
	TableName     string
	PrivilegeType string
}

func (q *Queries) GetMysqlRolePermissions(ctx context.Context, db DBTX) ([]*GetMysqlRolePermissionsRow, error) {
	rows, err := db.QueryContext(ctx, getMysqlRolePermissions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetMysqlRolePermissionsRow
	for rows.Next() {
		var i GetMysqlRolePermissionsRow
		if err := rows.Scan(&i.TableSchema, &i.TableName, &i.PrivilegeType); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPrimaryKeyConstraints = `-- name: GetPrimaryKeyConstraints :many
SELECT
	table_schema AS schema_name,
	table_name as table_name,
	column_name as column_name,
	constraint_name as constraint_name
FROM
	information_schema.key_column_usage
WHERE
	table_schema = ?
	AND constraint_name = 'PRIMARY'
ORDER BY
	table_name,
	column_name
`

type GetPrimaryKeyConstraintsRow struct {
	SchemaName     string
	TableName      string
	ColumnName     string
	ConstraintName string
}

func (q *Queries) GetPrimaryKeyConstraints(ctx context.Context, db DBTX, tableSchema string) ([]*GetPrimaryKeyConstraintsRow, error) {
	rows, err := db.QueryContext(ctx, getPrimaryKeyConstraints, tableSchema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetPrimaryKeyConstraintsRow
	for rows.Next() {
		var i GetPrimaryKeyConstraintsRow
		if err := rows.Scan(
			&i.SchemaName,
			&i.TableName,
			&i.ColumnName,
			&i.ConstraintName,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUniqueConstraints = `-- name: GetUniqueConstraints :many
SELECT
    tc.table_schema AS schema_name,
    tc.table_name AS table_name,
    tc.constraint_name AS constraint_name,
    kcu.column_name AS column_name
FROM
    information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
    AND tc.table_name = kcu.table_name
WHERE
    tc.table_schema = ?
    AND tc.constraint_type = 'UNIQUE'
ORDER BY
    tc.table_name,
    kcu.column_name
`

type GetUniqueConstraintsRow struct {
	SchemaName     string
	TableName      string
	ConstraintName string
	ColumnName     string
}

func (q *Queries) GetUniqueConstraints(ctx context.Context, db DBTX, tableSchema string) ([]*GetUniqueConstraintsRow, error) {
	rows, err := db.QueryContext(ctx, getUniqueConstraints, tableSchema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetUniqueConstraintsRow
	for rows.Next() {
		var i GetUniqueConstraintsRow
		if err := rows.Scan(
			&i.SchemaName,
			&i.TableName,
			&i.ConstraintName,
			&i.ColumnName,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
