<!-- -*- mode: gfm; auto-fill: t; fill-column: 78; -*- -->
# 📚 Validator Predicate Reference

This document describes the predicates supported by the `validator` package.
These predicates can be used with DAG-style rule evaluation engines and are
designed to operate on Go structs via reflection and dynamic field access.

## 🧠 Predicate Structure

Each predicate operates on a specific field within an input structure.
Predicates are structured as:

```yaml
field: "temperature"
predicate: "field-value-<"
value: 25
```

The `predicate` is a symbolic operator used to evaluate the value of a field.

### 🧩 Field vs Value Predicates

- **Field Type Predicates**: Operate on the Go type of a struct field.
- **Field Value Predicates**: Operate on the runtime value of a struct field.

---

# ✅ Supported Predicates

## 🏷 Field Type Predicates

These predicates inspect the *declared type* of a field.

### `field-type-equal` — Field Type Equals

**Description:**
Returns true if the field's Go type name matches the given string.

**Example:**
```yaml
field: "status"
predicate: "field-type-equal"
value: "string"
```

### `field-type-in` — Field Type In

**Description:**
Returns true if the field's type is in the list of accepted types.

**Example:**
```yaml
field: "payload"
predicate: "field-type-in"
value: ["map[string]any", "[]byte"]
```

## 🔢 Field Value Predicates

These compare the *runtime value* of a field.

### `field-value-equal` — Field Value Equals

**Description:**
Returns true if the value equals the expected value. Supports `int`, `float`,
`complex`, `string`, `bool`.

### `field-value-not-equal` — Field Value Not Equal

**Description:**
Logical inverse of `FVEQ`. Returns true if values differ.

### `field-value-<` — Field Value Less Than

**Description:**
Field's value is less than the predicate's value.

### `field-value-<=` — Field Value Less Than or Equal

**Description:**
Field's value is less than or equal to the predicate's value.

### `field-value->` — Field Value Greater Than

**Description:**
Field's value is greater than the predicate's value.

### `field-value->=` — Field Value Greater Than or Equal

**Description:**
Field's value is greater than or equal to the predicate's value.

### `field-value-in` — Field Value In

**Description:**
Returns true if the value is present in the provided list.

**Example:**
```yaml
field: "mode"
predicate: "field-value-in"
value: ["safe", "dry_run", "check"]
```

## 🔍 String / Regex Predicates

### `field-value-regex-match` — Field Value Regex Match

**Description:**
Returns true if the string value matches the provided regex.

**Example:**
```yaml
field: "email"
predicate: "field-value-regex-match"
value: "^[a-z0-9._%+-]+@example\.com$"
```

## ⚖ Logical Predicates

### `field-value-is-true` — Field Value Is Logically True

**Description:**
Returns true if the field's value is non-zero/non-empty.

> Structures are **never** considered logically true.

### `field-value-is-false` — Field Value Is Logically False

**Description:**
Returns true if the value is zero, empty, or `false`.

> Equivalent to Go’s `reflect.Value.IsZero()` in most cases.

### `field-value-is-nil` — Field Value Is Nil

**Description:**
Returns true only if the field’s reference value is `nil`.

> This does not match zero values like empty strings or zero numbers.

---

# 🛠 Validator Actions

These are built-in actions that can be triggered during validation when a
predicate chain evaluates to `true`.
They conform to the `dag.ActionFn` function signature and are compiled by the
`actions.Builder()` method.

## ✨ `none`

This action does nothing.

### Parameters

None.

### Example

```yaml
perform: "none"
params: {}
```

## 📢 `log`

This action writes a structured log message to the current DAG context logger.

### Parameters

| Name      | Type     | Required | Description                     |
|-----------|----------|----------|---------------------------------|
| `message` | `string` | ✅       | The message to print to the log |

### Example

```yaml
perform: "log"
params:
  message: "Validation path hit: status is critical"
```

**Output (at INFO level):**

```json
{
    "message": "Validation path hit: status is critical",
    "src": "log_action",
    "structure": { ...input fields... }
}
```

If the `message` parameter is missing or not a string, the action fails to
compile.

---

# ➕ Adding Custom Actions

To extend the validator, implement the `dag.Actions` interface and pass your
handler into:

```go
compiler := dag.NewCompilerWithPredicates(ctx, &yourCustomActions{}, dag.BuildPredicateDict())
```

---

# 💡 Examples

Let's consider the following Go structure:

``` go
    type DummyStructure struct {
        One   any
        Two   map[string]int
        Three any
        Four  string
        Five  string
    }
```

Let's say we have the following contrived constraints:

 * Field `One` must be:
   * of type `int64`
   * a value of `40`, `41`, `42`, `43`
 * Field `Two` must be:
   * of type `map[string]int`
   * logically true -- that is, not zero or empty.
 * Field `Three` must be:
   * `nil`
 * Field `Four` must be:
   * of type `string`
   * a value of `OK`, `CRITICAL`, `WARNING`
 * Field `Five` must be:
   * of type `string`
   * match the regular expression `.*coffee.*`

We can define the YAML rules thusly:

``` yaml

- name: "'One' must be int64 and between 40-43"
  conditions:
    - attribute: one
      operator: field-type-in
      value: [int8, int16, int32, int64]
    - attribute: one
      operator: field-value-in
      value: [40, 41, 42, 43]
  action:
    perform: none

- name: "'Two' must be map[string]int and not empty"
  conditions:
    - attribute: two
      operator: field-type-equal
      value: map[string]int
    - attribute: two
      operator: field-value-is-true
  action:
    perform: none

- name: "'three' must be nil"
  conditions:
    - attribute: three
      operator: field-value-is-nil
  action:
    perform: none

- name: "'four' must be string and member"
  conditions:
    - attribute: four
      operator: field-type-equal
      value: string
    - attribute: four
      operator: field-value-in
      value: [OK, CRITICAL, WARNING]
  action:
    perform: none

- name: "'five' must match regex"
  conditions:
    - attribute: five
      operator: field-type-equal
      value: string
    - attribute: five
      operator: field-value-regex-match
      value: ".*coffee.*"
  action:
    perform: none
```

In this case, I am using the special action `none` which simply results in
no action being performed.  In the future there is a possibility there will
be more actions one can do with the validator.

## Visualisation

![Example graph visualisation](./example.png "Example graph visualisation")

---

# 📌 Notes

- `FVEQ`/`FVNEQ` support cross-type comparisons (e.g. comparing `int64` to
  `int`), with width checks.
- `FVNIL` and `FVFALSE` are **not** the same — nil is reference-level, false is
  logical.
- Struct fields are never `FVFALSE` or `FVTRUE` — those require recursive
  validation.

---

# 🧑 See Also

 * [../dag/manual.md](../dag/manual.md)
 * [validator/predicates.go](../../validator/predicates.go)
 * [../go/validator.md](../go/validator.md)
