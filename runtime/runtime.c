#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

// Rush object types
typedef enum {
    RUSH_NULL = 0,
    RUSH_INTEGER = 1,
    RUSH_FLOAT = 2,
    RUSH_BOOLEAN = 3,
    RUSH_STRING = 4,
} RushObjectType;

// Rush object structure (matches LLVM definition)
typedef struct {
    int type;
    long long value;
} RushObject;

// Simple print function for string data (used by LLVM AOT)  
void rush_print(const char* str, size_t len) {
    write(STDOUT_FILENO, str, len);
}

// Print function with newline (matches interpreter behavior)
void rush_print_line(const char* str, size_t len) {
    write(STDOUT_FILENO, str, len);
    write(STDOUT_FILENO, "\n", 1);
}

// Print function for Rush objects
void rush_print_object(RushObject obj) {
    switch (obj.type) {
        case RUSH_INTEGER:
            printf("%lld", obj.value);
            break;
            
        case RUSH_FLOAT: {
            // Convert int64 bits back to double
            double* float_ptr = (double*)&obj.value;
            printf("%g", *float_ptr);
            break;
        }
        
        case RUSH_BOOLEAN:
            printf("%s", obj.value ? "true" : "false");
            break;
            
        case RUSH_STRING: {
            // Convert int64 back to char pointer
            char* str_ptr = (char*)obj.value;
            printf("%s", str_ptr);
            break;
        }
        
        case RUSH_NULL:
        default:
            printf("null");
            break;
    }
}

// Print with newline
void rush_println(RushObject obj) {
    rush_print_object(obj);
    printf("\n");
}

// Create Rush objects from C types
RushObject rush_make_integer(long long value) {
    RushObject obj = {RUSH_INTEGER, value};
    return obj;
}

RushObject rush_make_float(double value) {
    RushObject obj;
    obj.type = RUSH_FLOAT;
    obj.value = *(long long*)&value;
    return obj;
}

RushObject rush_make_boolean(int value) {
    RushObject obj = {RUSH_BOOLEAN, value ? 1 : 0};
    return obj;
}

RushObject rush_make_string(const char* value) {
    RushObject obj;
    obj.type = RUSH_STRING;
    obj.value = (long long)value;
    return obj;
}

RushObject rush_make_null() {
    RushObject obj = {RUSH_NULL, 0};
    return obj;
}

// Arithmetic operations
RushObject rush_add(RushObject left, RushObject right) {
    if (left.type == RUSH_INTEGER && right.type == RUSH_INTEGER) {
        return rush_make_integer(left.value + right.value);
    } else if (left.type == RUSH_FLOAT || right.type == RUSH_FLOAT) {
        double left_val = (left.type == RUSH_FLOAT) ? *(double*)&left.value : (double)left.value;
        double right_val = (right.type == RUSH_FLOAT) ? *(double*)&right.value : (double)right.value;
        return rush_make_float(left_val + right_val);
    }
    return rush_make_null();
}

RushObject rush_subtract(RushObject left, RushObject right) {
    if (left.type == RUSH_INTEGER && right.type == RUSH_INTEGER) {
        return rush_make_integer(left.value - right.value);
    } else if (left.type == RUSH_FLOAT || right.type == RUSH_FLOAT) {
        double left_val = (left.type == RUSH_FLOAT) ? *(double*)&left.value : (double)left.value;
        double right_val = (right.type == RUSH_FLOAT) ? *(double*)&right.value : (double)right.value;
        return rush_make_float(left_val - right_val);
    }
    return rush_make_null();
}

RushObject rush_multiply(RushObject left, RushObject right) {
    if (left.type == RUSH_INTEGER && right.type == RUSH_INTEGER) {
        return rush_make_integer(left.value * right.value);
    } else if (left.type == RUSH_FLOAT || right.type == RUSH_FLOAT) {
        double left_val = (left.type == RUSH_FLOAT) ? *(double*)&left.value : (double)left.value;
        double right_val = (right.type == RUSH_FLOAT) ? *(double*)&right.value : (double)right.value;
        return rush_make_float(left_val * right_val);
    }
    return rush_make_null();
}

RushObject rush_divide(RushObject left, RushObject right) {
    if (left.type == RUSH_INTEGER && right.type == RUSH_INTEGER) {
        if (right.value == 0) return rush_make_null(); // Division by zero
        return rush_make_integer(left.value / right.value);
    } else if (left.type == RUSH_FLOAT || right.type == RUSH_FLOAT) {
        double left_val = (left.type == RUSH_FLOAT) ? *(double*)&left.value : (double)left.value;
        double right_val = (right.type == RUSH_FLOAT) ? *(double*)&right.value : (double)right.value;
        if (right_val == 0.0) return rush_make_null(); // Division by zero
        return rush_make_float(left_val / right_val);
    }
    return rush_make_null();
}