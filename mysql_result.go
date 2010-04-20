/**
 * GoMySQL - A MySQL client library for Go
 * Copyright 2010 Phil Bayfield
 * This software is licensed under a Creative Commons Attribution-Share Alike 2.0 UK: England & Wales License
 * Further information on this license can be found here: http://creativecommons.org/licenses/by-sa/2.0/uk/
 */
package mysql

/**
 * All results stored using MySQLResult
 */
type MySQLResult struct {
	AffectedRows	uint64
	InsertId	uint64
	WarningCount	uint16
	Message		string
	
	Fields		[]*MySQLField
	FieldCount	uint64
	fieldsRead	uint64
	fieldsEOF	bool
	
	Rows		[]*MySQLRow
	RowCount	uint64
	rowsEOF		bool
	
	pointer		int
}

/**
 * Fetch the current row (as an array)
 */
func (res *MySQLResult) FetchRow() []interface{} {
	if res.RowCount > 0 {
		if len(res.Rows) > res.pointer {
			row := res.Rows[res.pointer].Data
			res.pointer ++
			return row
		}
	}
	return nil
}

/**
 * Fetch a map of the current row
 */
func (res *MySQLResult) FetchMap() map[string] interface{} {
	if res.RowCount > 0 {
		if len(res.Rows) > res.pointer {
			row := res.Rows[res.pointer].Data
			rowMap := make(map[string] interface{})
			for key, val := range row {
				rowMap[res.Fields[key].Name] = val
			}
			res.pointer ++
			return rowMap
		}
	}
	return nil
}

/**
 * Reset pointer
 */
func (res *MySQLResult) Reset() {
	res.pointer = 0
}

/**
 * Field definition
 */
type MySQLField struct {
	Name		string
	Length		uint32
	Type		byte
	Flags		*MySQLFieldFlags
	Decimals	uint8
}

/**
 * Field flags
 */
type MySQLFieldFlags struct {
	NotNull		bool
	PrimaryKey	bool
	UniqueKey	bool
	MultiKey	bool
	Blob		bool
	Unsigned	bool
	Zerofill	bool
	Binary		bool
	Enum		bool
	AutoIncrement	bool
	Timestamp	bool
	Set		bool
}

/**
 * Process flags setting found flags as boolean true
 * @todo This would probably faster using binary
 */
func (field *MySQLFieldFlags) process(flags uint16) {
	// MySQL 5.1 returns values larger than defined in protocol docs, ignore these for now
	if flags >= FLAG_UNKNOWN_4 { flags -= FLAG_UNKNOWN_4 }
	if flags >= FLAG_UNKNOWN_3 { flags -= FLAG_UNKNOWN_3 }
	if flags >= FLAG_UNKNOWN_2 { flags -= FLAG_UNKNOWN_2 }
	if flags >= FLAG_UNKNOWN_1 { flags -= FLAG_UNKNOWN_1 }
	// Populate struct with known flags
	if flags >= FLAG_SET {
		field.Set = true
		flags -= FLAG_SET
	}
	if flags >= FLAG_TIMESTAMP {
		field.Timestamp = true
		flags -= FLAG_TIMESTAMP
	}
	if flags >= FLAG_AUTO_INCREMENT {
		field.AutoIncrement = true
		flags -= FLAG_AUTO_INCREMENT
	}
	if flags >= FLAG_ENUM {
		field.Enum = true
		flags -= FLAG_ENUM
	}
	if flags >= FLAG_BINARY {
		field.Binary = true
		flags -= FLAG_BINARY
	}
	if flags >= FLAG_ZEROFILL {
		field.Zerofill = true
		flags -= FLAG_ZEROFILL
	}
	if flags >= FLAG_UNSIGNED {
		field.Unsigned = true
		flags -= FLAG_UNSIGNED
	}
	if flags >= FLAG_BLOB {
		field.Blob = true
		flags -= FLAG_BLOB
	}
	if flags >= FLAG_MULTIPLE_KEY {
		field.MultiKey = true
		flags -= FLAG_MULTIPLE_KEY
	}
	if flags >= FLAG_UNIQUE_KEY {
		field.UniqueKey = true
		flags -= FLAG_UNIQUE_KEY
	}
	if flags >= FLAG_PRI_KEY {
		field.PrimaryKey = true
		flags -= FLAG_PRI_KEY
	}
	if flags >= FLAG_NOT_NULL {
		field.NotNull = true
		flags -= FLAG_NOT_NULL
	}
}

/**
 * Row definition
 */
type MySQLRow struct {
	Data		[]interface{}
}
