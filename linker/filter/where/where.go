package where

import (
	"errors"
	"reflect"
	"strings"
)

type Where struct {
	Root *Node
}

func (w *Where) GetWhereType(key string) WhereType {
	if nodeType, exists := TYPES[key]; exists {
		return nodeType
	}
	return TYPES["__field__"]
}

func NewFromMap(filter map[string]interface{}) *Where {
	v := reflect.ValueOf(filter)
	w := Where{}
	root := Node{Type: TYPES["__root__"]}
	w.walk(&root, v)
	w.Root = &root
	return &w
}

func (w *Where) walk(parent *Node, v reflect.Value) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			w.walk(parent, v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			node := &Node{
				Parent: parent,
				Nodes:  make([]*Node, 0),
			}
			parent.AddChild(node)
			node.Type = w.GetWhereType(k.String())
			node.Value = v.MapIndex(k).Elem().Interface()
			node.Key = k.String()
			w.walk(node, v.MapIndex(k))
		}
	}
}

func (f *Where) walkGenerateSQLQuery(node *Node, query *Query) {
	if query.Error != nil {
		return
	}
	nt := node.Type
	switch nt.Key {
	case "__root__":
		for _, n := range node.Nodes {
			f.walkGenerateSQLQuery(n, query)
		}
	case "__field__":
		for _, n := range node.Nodes {
			f.walkGenerateSQLQuery(n, query)
		}
		if node.Parent != nil && node.Parent.LastChild() != node {
			if node.Parent.Type.Key == "__root__" {
				query.SQLAppendAND()
			} else {
				op, _ := node.Parent.Type.ApplySQL()
				query.Append(" " + op + " ")
			}
		}
	case "and", "or":
		pk := node.Parent.Type.Key
		if pk != "and" && pk != "or" && pk != "__root__" {
			query.Error = errors.New(strings.ToUpper(node.Key) + " operator in unsupported parent")
			return
		}
		query.Append("(")
		for _, n := range node.Nodes {
			f.walkGenerateSQLQuery(n, query)
		}
		query.Append(")")
		if node.Parent != nil && node.Parent.LastChild() != node {
			op, _ := node.Parent.Type.ApplySQL()
			query.Append(" " + op + " ")
		}
	case "eq", "neq", "gt", "gte", "lt", "lte", "like", "nlike":
		op, _ := node.Type.ApplySQL(node.Parent.Key)
		query.Append(op)
		query.Vars = append(query.Vars, node.Value)
		if node.Parent != nil && node.Parent.LastChild() != node {
			query.SQLAppendAND()
		}
	case "in", "nin":
		op, _ := node.Type.ApplySQL(node.Parent.Key)
		query.Append(op)
		if reflect.TypeOf(node.Value).Kind() != reflect.Slice {
			query.Error = errors.New(strings.ToUpper(node.Key) + " operator must be an array")
			return
		}
		query.Vars = append(query.Vars, node.Value)
		if node.Parent != nil && node.Parent.LastChild() != node {
			query.SQLAppendAND()
		}
	case "between":
		op, _ := node.Type.ApplySQL(node.Parent.Key)
		query.Append(op)
		b_range, valid := node.Value.([]interface{})
		if !valid || len(b_range) != 2 {
			query.Error = errors.New(strings.ToUpper(node.Key) + " operator must be an array with lenght = 2")
			return
		}
		query.Vars = append(query.Vars, b_range[0])
		query.Vars = append(query.Vars, b_range[1])
		if node.Parent != nil && node.Parent.LastChild() != node {
			query.SQLAppendAND()
		}
	}
}

func (w *Where) ToSQL() (Query, error) {
	q := Query{
		Text: "",
		Vars: make([]interface{}, 0),
	}
	w.walkGenerateSQLQuery(w.Root, &q)
	return q, q.Error
}
