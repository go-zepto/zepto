package where

type Query struct {
	Text  string
	Vars  []interface{}
	Error error
}

func (q *Query) Prepend(txt string) {
	q.Text = txt + q.Text
}

func (q *Query) Append(txt string) {
	q.Text = q.Text + txt
}

func (q *Query) SQLAppendAND() {
	and := TYPES["and"]
	op, _ := and.ApplySQL()
	q.Append(" " + op + " ")
}

func (q *Query) SQLAppendOR() {
	and := TYPES["or"]
	op, _ := and.ApplySQL()
	q.Append(" " + op + " ")
}

func (q *Query) AppendQuoted(txt string) {
	q.Append("\"" + txt + "\"")
}
