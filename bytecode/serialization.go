package bytecode

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"rush/interpreter"
)

const (
	// Magic number for Rush bytecode files
	MagicNumber uint32 = 0x52555348 // "RUSH" in hex
	// Version of bytecode format
	FormatVersion uint32 = 1
	// Cache directory name
	CacheDir = ".rush_cache"
)

// SerializedBytecode represents serialized bytecode with metadata
type SerializedBytecode struct {
	Magic        uint32
	Version      uint32
	Timestamp    int64
	SourceHash   [32]byte
	Instructions Instructions
	Constants    []SerializedValue
}

// SerializedValue represents a serialized constant value
type SerializedValue struct {
	Type ValueType
	Data []byte
}

// ValueType enum for serialization
type ValueType byte

const (
	IntegerType ValueType = iota
	FloatType
	StringType
	BooleanType
	NullType
	ArrayType
	HashType
	FunctionType
)

// Serialize converts bytecode and constants to binary format
func Serialize(instructions Instructions, constants []interpreter.Value, sourceHash [32]byte) ([]byte, error) {
	var buf bytes.Buffer

	// Write header
	err := binary.Write(&buf, binary.BigEndian, MagicNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to write magic number: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, FormatVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}

	timestamp := int64(time.Now().Unix())
	err = binary.Write(&buf, binary.BigEndian, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, sourceHash)
	if err != nil {
		return nil, fmt.Errorf("failed to write source hash: %w", err)
	}

	// Write instructions
	instructionsLen := uint32(len(instructions))
	err = binary.Write(&buf, binary.BigEndian, instructionsLen)
	if err != nil {
		return nil, fmt.Errorf("failed to write instructions length: %w", err)
	}

	_, err = buf.Write(instructions)
	if err != nil {
		return nil, fmt.Errorf("failed to write instructions: %w", err)
	}

	// Write constants
	constantsLen := uint32(len(constants))
	err = binary.Write(&buf, binary.BigEndian, constantsLen)
	if err != nil {
		return nil, fmt.Errorf("failed to write constants length: %w", err)
	}

	for _, constant := range constants {
		serializedValue, err := serializeValue(constant)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize constant: %w", err)
		}

		err = binary.Write(&buf, binary.BigEndian, serializedValue.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to write constant type: %w", err)
		}

		dataLen := uint32(len(serializedValue.Data))
		err = binary.Write(&buf, binary.BigEndian, dataLen)
		if err != nil {
			return nil, fmt.Errorf("failed to write constant data length: %w", err)
		}

		_, err = buf.Write(serializedValue.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to write constant data: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// Deserialize converts binary format back to bytecode and constants
func Deserialize(data []byte) (Instructions, []interpreter.Value, [32]byte, error) {
	buf := bytes.NewReader(data)

	// Read and verify header
	var magic uint32
	err := binary.Read(buf, binary.BigEndian, &magic)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("failed to read magic number: %w", err)
	}

	if magic != MagicNumber {
		return nil, nil, [32]byte{}, fmt.Errorf("invalid magic number: expected %x, got %x", MagicNumber, magic)
	}

	var version uint32
	err = binary.Read(buf, binary.BigEndian, &version)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("failed to read version: %w", err)
	}

	if version != FormatVersion {
		return nil, nil, [32]byte{}, fmt.Errorf("unsupported format version: %d", version)
	}

	// Skip timestamp for now
	var timestamp int64
	err = binary.Read(buf, binary.BigEndian, &timestamp)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("failed to read timestamp: %w", err)
	}

	// Read source hash
	var sourceHash [32]byte
	err = binary.Read(buf, binary.BigEndian, &sourceHash)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("failed to read source hash: %w", err)
	}

	// Read instructions
	var instructionsLen uint32
	err = binary.Read(buf, binary.BigEndian, &instructionsLen)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("failed to read instructions length: %w", err)
	}

	instructions := make(Instructions, instructionsLen)
	_, err = io.ReadFull(buf, instructions)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("failed to read instructions: %w", err)
	}

	// Read constants
	var constantsLen uint32
	err = binary.Read(buf, binary.BigEndian, &constantsLen)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("failed to read constants length: %w", err)
	}

	constants := make([]interpreter.Value, constantsLen)
	for i := uint32(0); i < constantsLen; i++ {
		var valueType ValueType
		err = binary.Read(buf, binary.BigEndian, &valueType)
		if err != nil {
			return nil, nil, [32]byte{}, fmt.Errorf("failed to read constant type: %w", err)
		}

		var dataLen uint32
		err = binary.Read(buf, binary.BigEndian, &dataLen)
		if err != nil {
			return nil, nil, [32]byte{}, fmt.Errorf("failed to read constant data length: %w", err)
		}

		data := make([]byte, dataLen)
		_, err = io.ReadFull(buf, data)
		if err != nil {
			return nil, nil, [32]byte{}, fmt.Errorf("failed to read constant data: %w", err)
		}

		value, err := deserializeValue(valueType, data)
		if err != nil {
			return nil, nil, [32]byte{}, fmt.Errorf("failed to deserialize constant: %w", err)
		}

		constants[i] = value
	}

	return instructions, constants, sourceHash, nil
}

// serializeValue converts a Rush value to serialized form
func serializeValue(value interpreter.Value) (SerializedValue, error) {
	var buf bytes.Buffer

	switch v := value.(type) {
	case *interpreter.Integer:
		err := binary.Write(&buf, binary.BigEndian, v.Value)
		if err != nil {
			return SerializedValue{}, err
		}
		return SerializedValue{Type: IntegerType, Data: buf.Bytes()}, nil

	case *interpreter.Float:
		err := binary.Write(&buf, binary.BigEndian, v.Value)
		if err != nil {
			return SerializedValue{}, err
		}
		return SerializedValue{Type: FloatType, Data: buf.Bytes()}, nil

	case *interpreter.String:
		_, err := buf.WriteString(v.Value)
		if err != nil {
			return SerializedValue{}, err
		}
		return SerializedValue{Type: StringType, Data: buf.Bytes()}, nil

	case *interpreter.Boolean:
		var b byte
		if v.Value {
			b = 1
		}
		buf.WriteByte(b)
		return SerializedValue{Type: BooleanType, Data: buf.Bytes()}, nil

	case *interpreter.Null:
		return SerializedValue{Type: NullType, Data: []byte{}}, nil

	case *interpreter.CompiledFunction:
		encoder := gob.NewEncoder(&buf)
		err := encoder.Encode(struct {
			Instructions  []byte
			NumLocals     int
			NumParameters int
		}{
			Instructions:  v.Instructions,
			NumLocals:     v.NumLocals,
			NumParameters: v.NumParameters,
		})
		if err != nil {
			return SerializedValue{}, err
		}
		return SerializedValue{Type: FunctionType, Data: buf.Bytes()}, nil

	default:
		return SerializedValue{}, fmt.Errorf("unsupported value type for serialization: %T", value)
	}
}

// deserializeValue converts serialized data back to a Rush value
func deserializeValue(valueType ValueType, data []byte) (interpreter.Value, error) {
	buf := bytes.NewReader(data)

	switch valueType {
	case IntegerType:
		var value int64
		err := binary.Read(buf, binary.BigEndian, &value)
		if err != nil {
			return nil, err
		}
		return &interpreter.Integer{Value: value}, nil

	case FloatType:
		var value float64
		err := binary.Read(buf, binary.BigEndian, &value)
		if err != nil {
			return nil, err
		}
		return &interpreter.Float{Value: value}, nil

	case StringType:
		return &interpreter.String{Value: string(data)}, nil

	case BooleanType:
		if len(data) != 1 {
			return nil, fmt.Errorf("invalid boolean data length")
		}
		return &interpreter.Boolean{Value: data[0] != 0}, nil

	case NullType:
		return &interpreter.Null{}, nil

	case FunctionType:
		decoder := gob.NewDecoder(buf)
		var fnData struct {
			Instructions  []byte
			NumLocals     int
			NumParameters int
		}
		err := decoder.Decode(&fnData)
		if err != nil {
			return nil, err
		}
		return &interpreter.CompiledFunction{
			Instructions:  fnData.Instructions,
			NumLocals:     fnData.NumLocals,
			NumParameters: fnData.NumParameters,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported value type for deserialization: %d", valueType)
	}
}

// HashSource creates a SHA-256 hash of source code
func HashSource(source string) [32]byte {
	return sha256.Sum256([]byte(source))
}

// Cache management functions

// GetCacheDir returns the cache directory path
func GetCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, CacheDir)
	
	// Create cache directory if it doesn't exist
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

// GetCacheFilePath returns the cache file path for a source file
func GetCacheFilePath(sourceFile string) (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}

	// Create a filename based on the source file path hash
	sourceHash := sha256.Sum256([]byte(sourceFile))
	filename := fmt.Sprintf("%x.rushc", sourceHash[:8])
	
	return filepath.Join(cacheDir, filename), nil
}

// SaveToCache saves bytecode to cache file
func SaveToCache(sourceFile string, instructions Instructions, constants []interpreter.Value, sourceHash [32]byte) error {
	cacheFile, err := GetCacheFilePath(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to get cache file path: %w", err)
	}

	data, err := Serialize(instructions, constants, sourceHash)
	if err != nil {
		return fmt.Errorf("failed to serialize bytecode: %w", err)
	}

	err = os.WriteFile(cacheFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// LoadFromCache loads bytecode from cache file
func LoadFromCache(sourceFile string, currentSourceHash [32]byte) (Instructions, []interpreter.Value, error) {
	cacheFile, err := GetCacheFilePath(sourceFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cache file path: %w", err)
	}

	// Check if cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("cache file does not exist")
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	instructions, constants, cachedSourceHash, err := Deserialize(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize bytecode: %w", err)
	}

	// Verify source hash matches
	if cachedSourceHash != currentSourceHash {
		return nil, nil, fmt.Errorf("source file has been modified, cache is stale")
	}

	return instructions, constants, nil
}

// ClearCache removes all cache files
func ClearCache() error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	err = os.RemoveAll(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	return nil
}

// GetCacheStats returns statistics about the cache
func GetCacheStats() (int, int64, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return 0, 0, err
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read cache directory: %w", err)
	}

	var totalSize int64
	fileCount := 0

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".rushc" {
			fileCount++
			
			info, err := entry.Info()
			if err == nil {
				totalSize += info.Size()
			}
		}
	}

	return fileCount, totalSize, nil
}