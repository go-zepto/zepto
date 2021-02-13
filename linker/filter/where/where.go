package where

import (
	"fmt"
	"reflect"
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
		query.Append("(")
		for _, n := range node.Nodes {
			f.walkGenerateSQLQuery(n, query)
		}
		query.Append(")")
		if node.Parent != nil && node.Parent.LastChild() != node {
			op, _ := node.Parent.Type.ApplySQL()
			query.Append(" " + op + " ")
		}
	case "eq", "neq":
		op, _ := node.Type.ApplySQL(node.Parent.Key)
		fmt.Println(node.Parent.Type.Key)
		fmt.Println(op)
		query.Append(op)
		query.Vars = append(query.Vars, node.Value)
		if node.Parent != nil && node.Parent.LastChild() != node {
			query.SQLAppendAND()
		}
	}
}

func (w *Where) ToSQL() Query {
	q := Query{
		Text: "",
		Vars: make([]interface{}, 0),
	}
	w.walkGenerateSQLQuery(w.Root, &q)
	return q
}
