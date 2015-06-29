package pm

import (
    "encoding/json"
)

type Args interface {
    GetInt(key string) int
    GetString(key string) string
    GetFloat(key string) float64
    GetMap(key string) map[string]interface{}
    GetArray(key string) []interface{}
    GetStringArray(key string) []string
    GetIntArray(key string) []int
    Set(key string, value interface{})
    Clone(deep bool) Args
}

type MapArgs struct {
    data map[string]interface{}
}

func NewMapArgs(data map[string]interface{}) Args{
    return &MapArgs{
        data: data,
    }
}

func (args *MapArgs) MarshalJSON() ([]byte, error) {
    return json.Marshal(args.data)
}

func (args *MapArgs) UnmarshalJSON(data []byte) error {
    return json.Unmarshal(data, &args.data)
}

func (args *MapArgs) GetInt(key string) int {
    s, ok := args.data[key]
    if ok {
        return s.(int)
    }
    return 0
}

func (args *MapArgs) GetString(key string) string {
    s, ok := args.data[key]
    if ok {
        return s.(string)
    }
    return ""
}

func (args *MapArgs) GetFloat(key string) float64 {
    s, ok := args.data[key]
    if ok {
        return s.(float64)
    }
    return 0
}

func (args *MapArgs) GetMap(key string) map[string]interface{} {
    s, ok := args.data[key]
    if ok {
        return s.(map[string]interface{})
    }

    return make(map[string]interface{})
}

func (args *MapArgs) GetArray(key string) []interface{} {
    s, ok := args.data[key]
    if ok {
        return s.([]interface{})
    }

    return make([]interface{}, 0)
}

func (args *MapArgs) GetStringArray(key string) []string {
    s, ok := args.data[key]
    if ok {
        return s.([]string)
    }

    return make([]string, 0)
}

func (args *MapArgs) GetIntArray(key string) []int {
    s, ok := args.data[key]
    if ok {
        return s.([]int)
    }

    return make([]int, 0)
}

func (args *MapArgs) Set(key string, value interface{}) {
    args.data[key] = value
}

func (args *MapArgs) Clone(deep bool) Args {
    data := make(map[string]interface{})
    for k, v := range args.data {
        if deep {
            switch tv := v.(type) {
            case []int:
                l := make([]int, len(tv))
                copy(l, tv)
                data[k] = l
            case []string:
                l := make([]string, len(tv))
                copy(l, tv)
                data[k] = l
            }
        } else {
            data[k] = v
        }
    }

    return &MapArgs{
        data: data,
    }
}
